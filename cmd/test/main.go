package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"
)

const (
	baseURL      = "http://localhost:8081"
	testEmail    = "test@example.com"
	testPassword = "Teste@7890Ab" // Senha forte que atende aos requisitos
)

type RegisterRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type RefreshRequest struct {
	RefreshToken string `json:"refresh_token"`
}

type ErrorResponse struct {
	Error string `json:"error"`
	Code  int    `json:"code"`
}

func main() {
	fmt.Println("Iniciando testes...")

	// Teste de refresh token
	fmt.Println("\nTestando fluxo de refresh token...")

	// Registrar usuário
	fmt.Println("Criando usuário de teste...")
	registerReq := RegisterRequest{
		Email:    testEmail,
		Password: testPassword,
	}
	registerBody, _ := json.Marshal(registerReq)

	resp, err := http.Post(baseURL+"/auth/register", "application/json", bytes.NewBuffer(registerBody))
	if err != nil {
		log.Fatalf("Erro ao fazer requisição de registro: %v", err)
	}

	body, _ := io.ReadAll(resp.Body)
	resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		var errResp ErrorResponse
		if err := json.Unmarshal(body, &errResp); err == nil {
			fmt.Printf("Aviso: erro ao registrar usuário. Status: %d, Erro: %s\n", resp.StatusCode, errResp.Error)
		} else {
			fmt.Printf("Aviso: erro ao registrar usuário. Status: %d, Body: %s\n", resp.StatusCode, string(body))
		}
		if resp.StatusCode != http.StatusBadRequest {
			log.Fatalf("Falha no registro do usuário")
		}
		fmt.Println("Usuário já existe, continuando com login...")
	} else {
		fmt.Println("Usuário criado com sucesso!")
	}

	// Login
	fmt.Println("\nFazendo login...")
	loginReq := LoginRequest{
		Email:    testEmail,
		Password: testPassword,
	}
	loginBody, _ := json.Marshal(loginReq)

	resp, err = http.Post(baseURL+"/auth/login", "application/json", bytes.NewBuffer(loginBody))
	if err != nil {
		log.Fatalf("Erro ao fazer requisição de login: %v", err)
	}

	body, _ = io.ReadAll(resp.Body)
	resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		var errResp ErrorResponse
		if err := json.Unmarshal(body, &errResp); err == nil {
			log.Fatalf("Login falhou. Status: %d, Erro: %s", resp.StatusCode, errResp.Error)
		} else {
			log.Fatalf("Login falhou. Status: %d, Body: %s", resp.StatusCode, string(body))
		}
	}

	var loginResp LoginResponse
	if err := json.Unmarshal(body, &loginResp); err != nil {
		log.Fatalf("Erro ao decodificar resposta do login: %v", err)
	}

	if loginResp.AccessToken == "" || loginResp.RefreshToken == "" {
		log.Fatal("Tokens vazios na resposta")
	}

	fmt.Println("Login realizado com sucesso!")

	// Aguardar um pouco para simular o uso do token
	fmt.Println("\nAguardando 2 segundos para simular uso do token...")
	time.Sleep(2 * time.Second)

	// Testar refresh token
	fmt.Println("Testando refresh token...")
	refreshReq := RefreshRequest{
		RefreshToken: loginResp.RefreshToken,
	}
	refreshBody, _ := json.Marshal(refreshReq)

	resp, err = http.Post(baseURL+"/auth/refresh", "application/json", bytes.NewBuffer(refreshBody))
	if err != nil {
		log.Fatalf("Erro ao fazer requisição de refresh: %v", err)
	}

	body, _ = io.ReadAll(resp.Body)
	resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		var errResp ErrorResponse
		if err := json.Unmarshal(body, &errResp); err == nil {
			log.Fatalf("Refresh falhou. Status: %d, Erro: %s", resp.StatusCode, errResp.Error)
		} else {
			log.Fatalf("Refresh falhou. Status: %d, Body: %s", resp.StatusCode, string(body))
		}
	}

	var refreshResp LoginResponse
	if err := json.Unmarshal(body, &refreshResp); err != nil {
		log.Fatalf("Erro ao decodificar resposta do refresh: %v", err)
	}

	if refreshResp.AccessToken == "" || refreshResp.RefreshToken == "" {
		log.Fatal("Tokens vazios na resposta do refresh")
	}

	if refreshResp.AccessToken == loginResp.AccessToken || refreshResp.RefreshToken == loginResp.RefreshToken {
		log.Fatal("Tokens não foram rotacionados corretamente")
	}

	fmt.Println("Refresh token rotacionado com sucesso!")

	// Testar refresh token antigo (deve falhar)
	fmt.Println("\nTestando refresh token antigo (deve falhar)...")
	oldRefreshReq := RefreshRequest{
		RefreshToken: loginResp.RefreshToken,
	}
	oldRefreshBody, _ := json.Marshal(oldRefreshReq)

	resp, err = http.Post(baseURL+"/auth/refresh", "application/json", bytes.NewBuffer(oldRefreshBody))
	if err != nil {
		log.Fatalf("Erro ao fazer requisição com refresh token antigo: %v", err)
	}

	body, _ = io.ReadAll(resp.Body)
	resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		log.Fatal("Refresh token antigo não foi invalidado corretamente")
	}

	fmt.Println("Refresh token antigo rejeitado corretamente!")
	fmt.Println("\nTeste concluído com sucesso! ✨")
}
