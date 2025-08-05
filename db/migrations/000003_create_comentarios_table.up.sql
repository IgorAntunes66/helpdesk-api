CREATE TABLE IF NOT EXISTS comentarios (
    id BIGSERIAL PRIMARY KEY,
    descricao TEXT NOT NULL,
    data TIMESTAMPTZ DEFAULT NOW(),
    user_id BIGINT REFERENCES users(id), -- Chave estrangeira para a tabela de usuários
    ticket_id BIGINT NOT NULL REFERENCES tickets(id) ON DELETE CASCADE -- Chave estrangeira para a tabela de tickets
);