CREATE TABLE IF NOT EXISTS projects (
    id             BIGSERIAL PRIMARY KEY,
    user_id        BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    name           VARCHAR(100) NOT NULL,
    repo_full_name VARCHAR(255) NOT NULL DEFAULT '',
    git_url        VARCHAR(255) NOT NULL,
    branch         VARCHAR(100) NOT NULL DEFAULT 'main',
    port           INT NOT NULL DEFAULT 8080,
    domain         VARCHAR(255) UNIQUE,
    namespace      VARCHAR(100) NOT NULL DEFAULT '',
    current_image  VARCHAR(255),
    status         VARCHAR(30) NOT NULL DEFAULT 'pending',
    webhook_id     BIGINT,
    env            JSONB NOT NULL DEFAULT '{}'::jsonb,
    created_at     TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at     TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_projects_user_id ON projects(user_id);
