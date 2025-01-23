FROM golang:1.22-bullseye AS builder

WORKDIR /app

# Copiar arquivos de dependências
COPY go.mod go.sum ./
RUN go mod download

# Copiar código fonte
COPY . .

# Compilar a aplicação
RUN CGO_ENABLED=1 GOOS=linux go build -o main ./cmd/api

# Imagem final
FROM debian:bullseye-slim

WORKDIR /app

# Instalar dependências necessárias
RUN apt-get update && \
    apt-get install -y --no-install-recommends ca-certificates && \
    rm -rf /var/lib/apt/lists/*

# Copiar binário compilado
COPY --from=builder /app/main .

# Expor porta da API
EXPOSE 8087

# Comando para executar a aplicação
CMD ["./main"] 