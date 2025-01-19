CREATE TABLE IF NOT EXISTS users (
    id BIGSERIAL PRIMARY KEY,
    email VARCHAR(255) NOT NULL UNIQUE,
    password VARCHAR(255) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_users_email ON users(email);

-- Inserir usuário de teste para desenvolvimento
INSERT INTO users (email, password, created_at, updated_at)
VALUES (
    'test@example.com',
    '$2a$10$ZkKPmWrC5DNzLQP1oEK3y.X0RHPqQltq0LMZVGKKVVvKp5pZHvQmG', -- senha: Test@123
    CURRENT_TIMESTAMP,
    CURRENT_TIMESTAMP
) ON CONFLICT (email) DO NOTHING; 