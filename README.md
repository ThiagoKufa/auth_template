# API de Autenticação KufaTech

API de autenticação desenvolvida em Go, fornecendo endpoints seguros para registro e autenticação de usuários.

## Características

- Registro e login de usuários
- Validação robusta de senhas e emails
- Autenticação via JWT com refresh tokens
- Proteção contra força bruta
- Rate limiting por IP
- Blacklist de tokens
- Logging estruturado
- Documentação detalhada

## Requisitos

- Go 1.22 ou superior
- PostgreSQL 16 ou superior
- Redis 7 ou superior
- Docker e Docker Compose (opcional)
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

3. Escolha um método para executar:

### Usando Docker (Recomendado)
```bash
# Inicia todos os serviços
docker compose up -d

# Verifica os logs
docker compose logs -f
```

### Manualmente
```bash
# Instala dependências
go mod download

# Executa migrações
make migrate

# Inicia o servidor
make run
```

## Documentação

A documentação completa da API está disponível em:

- [Guia de Autenticação](doc/auth_guide.md) - Detalhes sobre endpoints, erros e exemplos
- [Swagger/OpenAPI](http://localhost:8087/swagger/index.html) - Documentação interativa (quando o servidor estiver rodando)

## Validação de Senha

A senha deve atender aos seguintes critérios:
1. Mínimo de 8 caracteres
2. Pelo menos uma letra maiúscula
3. Pelo menos uma letra minúscula
4. Pelo menos um número
5. Pelo menos um caractere especial
6. Máximo de 3 caracteres repetidos consecutivos
7. Mínimo de 5 caracteres únicos
8. Não pode conter palavras comuns como "password", "123456", "qwerty"

## Estrutura do Projeto

```
.
├── cmd/                    # Pontos de entrada da aplicação
│   ├── api/               # Servidor API
│   └── migrate/           # Ferramenta de migração
├── config/                # Arquivos de configuração
├── doc/                   # Documentação
├── internal/              # Código interno da aplicação
│   ├── config/           # Configurações
│   ├── database/         # Acesso ao banco de dados
│   ├── entity/           # Entidades do domínio
│   ├── errors/           # Erros customizados
│   ├── handlers/         # Handlers HTTP
│   ├── interfaces/       # Interfaces e contratos
│   ├── middleware/       # Middlewares
│   ├── routes/           # Definição de rotas
│   └── services/         # Lógica de negócio
├── pkg/                  # Pacotes reutilizáveis
│   ├── auth/            # Autenticação e tokens
│   ├── database/        # Utilitários de banco
│   ├── logger/          # Sistema de logs
│   └── validation/      # Validação de dados
└── scripts/             # Scripts utilitários
```

## Endpoints Principais

### Autenticação
- `POST /auth/register` - Registro de usuário
- `POST /auth/login` - Login com email/senha
- `POST /auth/refresh` - Renovação de tokens
- `POST /auth/logout` - Logout (invalidação de token)
- `GET /auth/me` - Dados do usuário atual

### Sistema
- `GET /health` - Status da API e recursos

Para exemplos detalhados de uso, consulte o [Guia de Autenticação](doc/auth_guide.md).

## Segurança

1. **Proteção de Senhas**:
   - Hash bcrypt com custo configurável
   - Validação robusta de força da senha
   - Proteção contra senhas comuns

2. **Tokens**:
   - Access tokens de curta duração (15min)
   - Refresh tokens de longa duração (30 dias)
   - Rotação automática de refresh tokens
   - Blacklist de tokens invalidados

3. **Rate Limiting**:
   - 100 requisições por hora por IP
   - Proteção contra força bruta
   - Blacklist temporária de IPs suspeitos

## Contribuindo

1. Fork o projeto
2. Crie uma branch para sua feature (`git checkout -b feature/nova-feature`)
3. Commit suas mudanças (`git commit -am 'Adiciona nova feature'`)
4. Push para a branch (`git push origin feature/nova-feature`)
5. Crie um Pull Request

## Melhorias Concluídas ✅

1. **Organização do Código**:
   - [x] Consolidação dos pacotes de validação
   - [x] Reorganização da estrutura de arquivos
   - [x] Documentação atualizada

2. **Segurança**:
   - [x] Validação robusta de senhas
   - [x] Proteção contra tokens inválidos
   - [x] Rate limiting implementado

3. **Configuração**:
   - [x] Configurações unificadas
   - [x] Suporte a múltiplos ambientes
   - [x] Docker Compose otimizado

## Próximos Passos

1. **Documentação**:
   - [ ] Implementar Swagger/OpenAPI
   - [ ] Adicionar exemplos em outras linguagens
   - [ ] Criar guia de contribuição

2. **Testes**:
   - [ ] Aumentar cobertura de testes
   - [ ] Adicionar testes de integração
   - [ ] Implementar testes de carga

3. **Monitoramento**:
   - [ ] Integrar Prometheus
   - [ ] Configurar Grafana
   - [ ] Adicionar tracing distribuído

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
