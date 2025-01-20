package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/go-chi/chi/v5"
	chimiddleware "github.com/go-chi/chi/v5/middleware"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"

	"auth-template/internal/config"
	"auth-template/internal/di"
	"auth-template/internal/middleware"
	"auth-template/internal/routes"
)

type Application struct {
	container *di.Container
	router    http.Handler
}

func initializeTestApplication() *Application {
	// Configurar ambiente de teste
	os.Setenv("APP_ENV", "test")
	os.Setenv("DATABASE_URL", "postgres://postgres:postgres@localhost:5000/kufatech_test?sslmode=disable")

	// Carregar configuração de teste
	cfg, err := config.Load()
	if err != nil {
		panic(err)
	}

	// Inicializar container com dependências
	container, err := di.InitializeContainer(cfg)
	if err != nil {
		panic(err)
	}

	// Configurar rotas
	router := setupRouter(container)

	return &Application{
		container: container,
		router:    router,
	}
}

var app *Application
var db *gorm.DB

func TestMain(m *testing.M) {
	// Inicializa a aplicação com configuração de teste
	app = initializeTestApplication()
	db = app.container.DB

	// Limpa o banco antes de executar os testes
	cleanDatabase()

	// Executa os testes
	code := m.Run()

	// Limpa o ambiente
	os.Unsetenv("APP_ENV")
	os.Unsetenv("DATABASE_URL")

	os.Exit(code)
}

func cleanDatabase() {
	// Limpa todas as tabelas relevantes
	db.Exec("DELETE FROM users")
}

func setupRouter(container *di.Container) http.Handler {
	r := chi.NewRouter()

	rateLimiter := middleware.NewRateLimiter(100, time.Minute) // Limite maior para testes
	r.Use(
		chimiddleware.Compress(5),
		chimiddleware.Timeout(30*time.Second),
		rateLimiter.RateLimit,
	)

	routes.SetupRoutes(r, container.Logger, container.AuthHandler, container.HealthHandler)
	return r
}

func TestAuthEndpoints(t *testing.T) {
	t.Run("Registro_com_sucesso", func(t *testing.T) {
		cleanDatabase()

		body := map[string]string{
			"email":    "test@example.com",
			"password": "Teste@7890Ab",
		}
		jsonBody, _ := json.Marshal(body)
		req := httptest.NewRequest(http.MethodPost, "/auth/register", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		app.router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusCreated, w.Code)
	})

	t.Run("Registro_com_email_duplicado", func(t *testing.T) {
		cleanDatabase()

		// Primeiro registro
		body := map[string]string{
			"email":    "test@example.com",
			"password": "Teste@7890Ab",
		}
		jsonBody, _ := json.Marshal(body)
		req := httptest.NewRequest(http.MethodPost, "/auth/register", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		app.router.ServeHTTP(w, req)

		// Tenta registrar o mesmo email
		req = httptest.NewRequest(http.MethodPost, "/auth/register", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		w = httptest.NewRecorder()
		app.router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusConflict, w.Code)
		var response map[string]string
		json.Unmarshal(w.Body.Bytes(), &response)
		assert.Contains(t, response["error"], "email já cadastrado")
	})

	t.Run("Login_com_sucesso", func(t *testing.T) {
		cleanDatabase()

		// Registra um usuário
		body := map[string]string{
			"email":    "test@example.com",
			"password": "Teste@7890Ab",
		}
		jsonBody, _ := json.Marshal(body)
		req := httptest.NewRequest(http.MethodPost, "/auth/register", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		app.router.ServeHTTP(w, req)

		// Tenta fazer login
		req = httptest.NewRequest(http.MethodPost, "/auth/login", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		w = httptest.NewRecorder()
		app.router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		var response map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &response)
		assert.Contains(t, response, "access_token")
		assert.Contains(t, response, "refresh_token")
	})

	t.Run("Login_com_senha_incorreta", func(t *testing.T) {
		cleanDatabase()

		// Registra um usuário
		body := map[string]string{
			"email":    "test@example.com",
			"password": "Teste@7890Ab",
		}
		jsonBody, _ := json.Marshal(body)
		req := httptest.NewRequest(http.MethodPost, "/auth/register", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		app.router.ServeHTTP(w, req)

		// Tenta fazer login com senha errada
		body["password"] = "wrongpassword"
		jsonBody, _ = json.Marshal(body)
		req = httptest.NewRequest(http.MethodPost, "/auth/login", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		w = httptest.NewRecorder()
		app.router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
		var response map[string]string
		json.Unmarshal(w.Body.Bytes(), &response)
		assert.Contains(t, response["error"], "credenciais inválidas")
	})

	t.Run("Login_com_email_não_cadastrado", func(t *testing.T) {
		cleanDatabase()

		body := map[string]string{
			"email":    "naoexiste@example.com",
			"password": "Teste@7890Ab",
		}
		jsonBody, _ := json.Marshal(body)
		req := httptest.NewRequest(http.MethodPost, "/auth/login", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		app.router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
		var response map[string]string
		json.Unmarshal(w.Body.Bytes(), &response)
		assert.Contains(t, response["error"], "credenciais inválidas")
	})

	t.Run("Senha_muito_curta", func(t *testing.T) {
		cleanDatabase()

		body := map[string]string{
			"email":    "test@example.com",
			"password": "Abc@1",
		}
		jsonBody, _ := json.Marshal(body)
		req := httptest.NewRequest(http.MethodPost, "/auth/register", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		app.router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		var response map[string]string
		json.Unmarshal(w.Body.Bytes(), &response)
		assert.Contains(t, response["error"], "senha deve ter pelo menos 8 caracteres")
	})

	t.Run("Email_inválido", func(t *testing.T) {
		cleanDatabase()

		body := map[string]string{
			"email":    "invalid-email",
			"password": "Teste@7890Ab",
		}
		jsonBody, _ := json.Marshal(body)
		req := httptest.NewRequest(http.MethodPost, "/auth/register", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		app.router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		var response map[string]string
		json.Unmarshal(w.Body.Bytes(), &response)
		assert.Contains(t, response["error"], "email inválido")
	})

	t.Run("Requisição_malformada", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/auth/register", bytes.NewBufferString("invalid json"))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		app.router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		var response map[string]string
		json.Unmarshal(w.Body.Bytes(), &response)
		assert.Contains(t, response["error"], "requisição inválida")
	})

	t.Run("Método_não_permitido", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/auth/register", nil)
		w := httptest.NewRecorder()
		app.router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusMethodNotAllowed, w.Code)
	})

	t.Run("Senha_sem_caractere_especial", func(t *testing.T) {
		cleanDatabase()

		body := map[string]string{
			"email":    "test@example.com",
			"password": "Teste1234Ab",
		}
		jsonBody, _ := json.Marshal(body)
		req := httptest.NewRequest(http.MethodPost, "/auth/register", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		app.router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		var response map[string]string
		json.Unmarshal(w.Body.Bytes(), &response)
		assert.Contains(t, response["error"], "senha deve conter pelo menos um caractere especial")
	})

	t.Run("Senha_sem_número", func(t *testing.T) {
		cleanDatabase()

		body := map[string]string{
			"email":    "test@example.com",
			"password": "Teste@abcdef",
		}
		jsonBody, _ := json.Marshal(body)
		req := httptest.NewRequest(http.MethodPost, "/auth/register", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		app.router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		var response map[string]string
		json.Unmarshal(w.Body.Bytes(), &response)
		assert.Contains(t, response["error"], "senha deve conter pelo menos um número")
	})

	t.Run("Senha_com_palavra_proibida", func(t *testing.T) {
		cleanDatabase()

		body := map[string]string{
			"email":    "test@example.com",
			"password": "Password@123",
		}
		jsonBody, _ := json.Marshal(body)
		req := httptest.NewRequest(http.MethodPost, "/auth/register", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		app.router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		var response map[string]string
		json.Unmarshal(w.Body.Bytes(), &response)
		assert.Contains(t, response["error"], "senha contém uma sequência de caracteres proibida")
	})

	t.Run("Health_Check", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/health", nil)
		w := httptest.NewRecorder()
		app.router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		var response map[string]string
		json.Unmarshal(w.Body.Bytes(), &response)
		assert.Equal(t, "healthy", response["status"])
	})
}

// ... rest of the tests ...
