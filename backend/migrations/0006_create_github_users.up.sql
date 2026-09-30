CREATE TABLE IF NOT EXISTS github_users (
    id         BIGINT PRIMARY KEY,
    login      VARCHAR(255) NOT NULL,
    name       VARCHAR(255) NOT NULL DEFAULT '',
    email      VARCHAR(255) NOT NULL DEFAULT '',
    avatar_url VARCHAR(255) NOT NULL DEFAULT '',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
