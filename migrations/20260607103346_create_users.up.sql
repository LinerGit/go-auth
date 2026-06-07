CREATE TABLE users (
    id BIGSERIAL PRIMARY KEY,

    username VARCHAR(100)
    UNIQUE NOT NULL,

    password_hash TEXT NOT NULL,

    role VARCHAR(20)
    NOT NULL DEFAULT 'USER',

    created_at TIMESTAMP
    NOT NULL DEFAULT NOW()
);