FROM golang:1.22-alpine AS builder

WORKDIR /app

# Instalar dependências de build
RUN apk add --no-cache gcc musl-dev

# Copiar arquivos de dependências
COPY go.mod go.sum ./
RUN go mod download

# Copiar código fonte
COPY . .

# Compilar a aplicação
RUN CGO_ENABLED=1 GOOS=linux go build -o main ./cmd/api

# Imagem final
FROM alpine:latest

WORKDIR /app

# Copiar binário compilado
COPY --from=builder /app/main .
COPY --from=builder /app/.env.docker ./.env

# Expor porta da API
EXPOSE 8081

# Comando para executar a aplicação
CMD ["./main"] 