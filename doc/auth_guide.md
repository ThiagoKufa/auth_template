# Guia de Autenticação

Este documento descreve o fluxo de autenticação da API, incluindo todos os endpoints disponíveis e exemplos de uso.

## Fluxo de Autenticação

1. O usuário se registra usando o endpoint `/auth/register`
2. O usuário faz login usando `/auth/login` e recebe um par de tokens:
   - `access_token`: usado para acessar endpoints protegidos
   - `refresh_token`: usado para obter novos tokens quando o access_token expirar
3. O cliente usa o `access_token` para fazer requisições autenticadas
4. Quando o `access_token` expirar (15 minutos), o cliente usa o `refresh_token` para obter um novo par de tokens
5. O processo se repete até o usuário fazer logout ou o `refresh_token` expirar (30 dias)

## Endpoints

### 1. Registro de Usuário
- **Endpoint**: `POST /auth/register`
- **Descrição**: Registra um novo usuário no sistema
- **Corpo da Requisição**:
```json
{
    "email": "usuario@exemplo.com",
    "password": "senha123"
}
```
- **Resposta de Sucesso**: `201 Created`
- **Possíveis Erros**:
  - `400 Bad Request`:
    - "email inválido": Formato de email incorreto
    - "email muito longo": Email excede 255 caracteres
    - "domínio do email inválido": Domínio não atende aos requisitos
    - "senha deve ter pelo menos X caracteres": Comprimento mínimo não atingido
    - "senha não pode ter mais que X caracteres repetidos": Muitos caracteres repetidos
    - "senha deve ter pelo menos X caracteres únicos": Poucos caracteres únicos
    - "senha deve conter pelo menos uma letra maiúscula": Falta letra maiúscula
    - "senha deve conter pelo menos uma letra minúscula": Falta letra minúscula
    - "senha deve conter pelo menos um número": Falta número
    - "senha deve conter pelo menos um caractere especial": Falta caractere especial
    - "senha contém uma sequência de caracteres proibida": Senha muito comum ou insegura
  - `409 Conflict`: "email já cadastrado"

### 2. Login
- **Endpoint**: `POST /auth/login`
- **Descrição**: Autentica o usuário e retorna tokens de acesso
- **Corpo da Requisição**:
```json
{
    "email": "usuario@exemplo.com",
    "password": "senha123"
}
```
- **Resposta de Sucesso** (200 OK):
```json
{
    "access_token": "eyJhbGciOiJIUzI1...",
    "refresh_token": "eyJhbGciOiJIUzI1..."
}
```
- **Possíveis Erros**:
  - `401 Unauthorized`: "credenciais inválidas"

### 3. Refresh Token
- **Endpoint**: `POST /auth/refresh`
- **Descrição**: Renova os tokens usando um refresh token válido
- **Importante**: 
  - O refresh token usado será automaticamente invalidado
  - Um novo par de tokens (access + refresh) será retornado
  - O novo refresh token deve ser armazenado para futuras renovações
- **Corpo da Requisição**:
```json
{
    "refresh_token": "eyJhbGciOiJIUzI1..."
}
```
- **Resposta de Sucesso** (200 OK):
```json
{
    "access_token": "eyJhbGciOiJIUzI1...",
    "refresh_token": "eyJhbGciOiJIUzI1..."
}
```
- **Possíveis Erros**:
  - `401 Unauthorized`: 
    - "refresh token inválido": Token expirado, malformado ou já utilizado
    - "erro ao verificar token": Erro ao verificar blacklist
    - "erro ao invalidar token": Erro ao adicionar token à blacklist

### 4. Logout
- **Endpoint**: `POST /auth/logout`
- **Descrição**: Invalida o refresh token atual
- **Corpo da Requisição**:
```json
{
    "refresh_token": "eyJhbGciOiJIUzI1..."
}
```
- **Resposta de Sucesso**: `204 No Content`
- **Possíveis Erros**:
  - `401 Unauthorized`: "refresh token inválido"

### 5. Dados do Usuário
- **Endpoint**: `GET /auth/me`
- **Descrição**: Retorna os dados do usuário autenticado
- **Headers**: 
  - `Authorization: Bearer <access_token>`
- **Resposta de Sucesso** (200 OK):
```json
{
    "id": 1,
    "email": "usuario@exemplo.com"
}
```
- **Possíveis Erros**:
  - `401 Unauthorized`: "token inválido"

## Requisitos de Senha

A senha deve atender aos seguintes critérios:
1. Mínimo de 8 caracteres
2. Pelo menos uma letra maiúscula
3. Pelo menos uma letra minúscula
4. Pelo menos um número
5. Pelo menos um caractere especial
6. Máximo de 3 caracteres repetidos consecutivos
7. Mínimo de 5 caracteres únicos
8. Não pode conter palavras comuns como "password", "123456", "qwerty"

## Exemplos de Uso com cURL

### Registro
```bash
curl -X POST http://localhost:8087/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "email": "usuario@exemplo.com",
    "password": "Senha@123"
  }'
```

### Login
```bash
curl -X POST http://localhost:8087/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "usuario@exemplo.com",
    "password": "Senha@123"
  }'
```

### Refresh Token
```bash
curl -X POST http://localhost:8087/auth/refresh \
  -H "Content-Type: application/json" \
  -d '{
    "refresh_token": "seu_refresh_token_aqui"
  }'
```

### Logout
```bash
curl -X POST http://localhost:8087/auth/logout \
  -H "Content-Type: application/json" \
  -d '{
    "refresh_token": "seu_refresh_token_aqui"
  }'
```

### Dados do Usuário
```bash
curl -X GET http://localhost:8087/auth/me \
  -H "Authorization: Bearer seu_access_token_aqui"
```

## Notas Importantes

1. **Tokens**:
   - O `access_token` tem validade de 15 minutos
   - O `refresh_token` tem validade de 30 dias
   - Cada refresh token só pode ser usado uma vez
   - Após usar um refresh token, ele é invalidado e um novo é gerado
   - Sempre armazene os tokens de forma segura

2. **Segurança**:
   - Todas as senhas são armazenadas com hash bcrypt
   - Os tokens são invalidados após o logout
   - Rate limiting de 100 requisições por hora por IP
   - Refresh tokens usados são automaticamente invalidados (rotação de tokens)

3. **Boas Práticas**:
   - Sempre use HTTPS em produção
   - Renove o access_token antes de expirar
   - Nunca envie tokens em URLs
   - Armazene tokens apenas em locais seguros (ex: HttpOnly cookies)
   - Mantenha apenas o refresh token mais recente
   - Implemente renovação automática do access_token quando faltar 1 minuto para expirar 