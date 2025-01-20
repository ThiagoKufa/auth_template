# API de Autenticação KufaTech

API de autenticação desenvolvida em Go, fornecendo endpoints seguros para registro e autenticação de usuários.

## Características

- Registro e login de usuários
- Validação robusta de senhas
- Autenticação via JWT (JSON Web Tokens)
- Health check para monitoramento
- Rate limiting para proteção contra ataques
- Logging estruturado
- Testes automatizados
- Documentação via Swagger

## Requisitos

- Go 1.22 ou superior
- PostgreSQL 14 ou superior
- Make

## Configuração

1. Clone o repositório:
```bash
git clone https://github.com/seu-usuario/auth-template.git
cd auth-template
```

2. Configure as variáveis de ambiente:
```bash
cp .env.example .env
# Edite o arquivo .env com suas configurações
```

3. Instale as dependências:
```bash
go mod download
```

4. Execute as migrações do banco de dados:
```bash
make migrate
```

## Executando o Projeto

1. Inicie o servidor:
```bash
make run
```

2. Execute os testes:
```bash
make test
```

3. Verifique o código com o linter:
```bash
make lint
```

## Endpoints

### Health Check
```
GET /health
```
Retorna o status da API e informações do sistema

### Registro
```
POST /auth/register
Content-Type: application/json

{
    "email": "usuario@exemplo.com",
    "password": "Senha@123"
}
```

### Login
```
POST /auth/login
Content-Type: application/json

{
    "email": "usuario@exemplo.com",
    "password": "Senha@123"
}
```

## Validação de Senha

A senha deve atender aos seguintes critérios:
- Mínimo de 8 caracteres
- Pelo menos uma letra maiúscula
- Pelo menos uma letra minúscula
- Pelo menos um número
- Pelo menos um caractere especial
- Não pode conter palavras comuns como "password", "123456", etc.

## Estrutura do Projeto

```
.
├── cmd/                    # Pontos de entrada da aplicação
│   ├── api/               # Servidor API
│   └── migrate/           # Ferramenta de migração
├── internal/              # Código interno da aplicação
│   ├── auth/             # Lógica de autenticação
│   ├── config/           # Configurações
│   ├── database/         # Acesso ao banco de dados
│   ├── handlers/         # Handlers HTTP
│   ├── middleware/       # Middlewares
│   ├── models/           # Modelos de dados
│   ├── routes/           # Definição de rotas
│   ├── services/         # Lógica de negócio
│   └── validation/       # Validação de dados
├── pkg/                  # Pacotes reutilizáveis
├── migrations/           # Migrações do banco de dados
└── tests/               # Testes
```

## Contribuindo

1. Fork o projeto
2. Crie uma branch para sua feature (`git checkout -b feature/nova-feature`)
3. Commit suas mudanças (`git commit -am 'Adiciona nova feature'`)
4. Push para a branch (`git push origin feature/nova-feature`)
5. Crie um Pull Request

## Licença

Este projeto está licenciado sob a licença MIT - veja o arquivo [LICENSE](LICENSE) para mais detalhes.

## Prioridade Alta (Sprint 1)
- [ ] **Documentação da API**
  - [ ] Implementar Swagger/OpenAPI
  - [ ] Documentar todos os endpoints existentes
  - [ ] Criar exemplos de uso

- [ ] **Testes Essenciais**
  - [ ] Adicionar testes unitários para serviços críticos
  - [ ] Implementar testes de integração básicos
  - [ ] Configurar ambiente de testes automatizados

- [ ] **Segurança Básica**
  - [ ] Implementar rate limiting por IP
  - [ ] Adicionar validação de entrada em todos endpoints
  - [x] Revisar e atualizar headers de segurança (já implementado no middleware)

## Prioridade Média (Sprint 2)
- [ ] **Monitoramento**
  - [ ] Configurar Prometheus
  - [ ] Implementar métricas básicas
  - [ ] Criar dashboard de monitoramento

- [ ] **CI/CD**
  - [ ] Configurar GitHub Actions
  - [ ] Implementar pipeline de build e teste
  - [ ] Configurar deploy automático

- [ ] **Otimização de Banco**
  - [ ] Revisar e otimizar índices
  - [x] Implementar soft delete (concluído com a adição da coluna deleted_at)
  - [ ] Otimizar queries principais

## Prioridade Normal (Sprint 3)
- [ ] **Cache**
  - [x] Configurar Redis (já configurado no docker-compose)
  - [ ] Definir estratégia de invalidação
  - [ ] Cachear endpoints mais acessados

- [ ] **Logs**
  - [ ] Migrar para Zerolog
  - [ ] Implementar logs estruturados
  - [ ] Configurar rotação de logs

- [ ] **Docker**
  - [x] Configurar ambiente básico (concluído)
  - [x] Configurar health checks (implementado no docker-compose)
  - [x] Configurar volumes persistentes (implementado)
  - [ ] Otimizar Dockerfile com multi-stage

## Melhorias Concluídas ✅
1. Organização do Código:
   - [x] Mesclagem de models com entity
   - [x] Consolidação dos repositórios
   - [x] Limpeza de arquivos não utilizados

2. Banco de Dados:
   - [x] Implementação de migrations
   - [x] Configuração do PostgreSQL
   - [x] Implementação de soft delete

3. Infraestrutura:
   - [x] Configuração do Docker
   - [x] Configuração do Redis
   - [x] Health checks básicos

4. Arquitetura:
   - [x] Injeção de dependências (wire)
   - [x] Middleware básicos (CORS, Security Headers, Rate Limit)
   - [x] Estrutura base do projeto

## Próximos Passos Recomendados
1. Implementar documentação com Swagger/OpenAPI
2. Adicionar testes unitários e de integração
3. Implementar validação de entrada nos endpoints
4. Configurar monitoramento com Prometheus

Gostaria de começar com alguma dessas tarefas pendentes?
