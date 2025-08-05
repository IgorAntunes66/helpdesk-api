CREATE TABLE IF NOT EXISTS tickets (
    id BIGSERIAL PRIMARY KEY,
    titulo VARCHAR(255) NOT NULL,
    descricao TEXT,
    status VARCHAR(50),
    diagnostico TEXT,
    solucao TEXT,
    prioridade VARCHAR(50),
    data_abertura TIMESTAMPTZ,
    data_fechamento TIMESTAMPTZ,
    data_atualizacao TIMESTAMPTZ,
    anexos TEXT[],
    tags TEXT[],
    categoria_id BIGINT,
    responsavel_id BIGINT,
    user_id BIGINT
);