.PHONY: all build test clean run docker-up docker-down lint help

# Variáveis
APP_NAME=kufatech
MAIN_PATH=./cmd/api
BUILD_DIR=./build
DOCKER_COMPOSE=docker-compose

# Cores para output
GREEN=\033[0;32m
NC=\033[0m # No Color

help: ## Mostra esta mensagem de ajuda
	@echo "Comandos disponíveis:"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "$(GREEN)%-20s$(NC) %s\n", $$1, $$2}'

all: clean build test ## Limpa, compila e testa o projeto

build: ## Compila o projeto
	@echo "Compilando o projeto..."
	@go build -o $(BUILD_DIR)/$(APP_NAME) $(MAIN_PATH)

run: ## Executa o projeto localmente
	@echo "Executando o projeto..."
	@go run $(MAIN_PATH)

test: ## Executa os testes
	@echo "Executando testes..."
	@go test -v ./... -cover

test-coverage: ## Executa os testes com relatório de cobertura
	@echo "Executando testes com cobertura..."
	@go test -coverprofile=coverage.out ./...
	@go tool cover -html=coverage.out -o coverage.html
	@rm coverage.out

clean: ## Limpa arquivos gerados
	@echo "Limpando arquivos gerados..."
	@rm -rf $(BUILD_DIR)
	@go clean
	@find . -type f -name '*.test' -delete
	@find . -type f -name '*.out' -delete
	@find . -type f -name '*.html' -delete

docker-build: ## Constrói as imagens Docker
	@echo "Construindo imagens Docker..."
	@$(DOCKER_COMPOSE) build

docker-up: ## Inicia os containers Docker
	@echo "Iniciando containers..."
	@$(DOCKER_COMPOSE) up -d

docker-down: ## Para os containers Docker
	@echo "Parando containers..."
	@$(DOCKER_COMPOSE) down -v

docker-logs: ## Mostra logs dos containers
	@$(DOCKER_COMPOSE) logs -f

lint: ## Executa o linter
	@echo "Executando linter..."
	@golangci-lint run

migrate-up: ## Executa as migrações do banco de dados
	@echo "Executando migrações..."
	@go run cmd/migrate/main.go up

migrate-down: ## Reverte as migrações do banco de dados
	@echo "Revertendo migrações..."
	@go run cmd/migrate/main.go down

dev: docker-up ## Inicia o ambiente de desenvolvimento
	@echo "Ambiente de desenvolvimento iniciado"
	@make run

mock: ## Gera os mocks para testes
	@echo "Gerando mocks..."
	@mockgen -source=internal/interfaces/repository/user_repository.go -destination=internal/mocks/repository/user_repository_mock.go
	@mockgen -source=internal/interfaces/service/auth.go -destination=internal/mocks/service/auth_service_mock.go

deps: ## Instala dependências de desenvolvimento
	@echo "Instalando dependências..."
	@go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	@go install github.com/golang/mock/mockgen@latest

# Define o alvo padrão
.DEFAULT_GOAL := help 