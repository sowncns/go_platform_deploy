CREATE TABLE IF NOT EXISTS deployments (
    id          BIGSERIAL PRIMARY KEY,
    project_id  BIGINT NOT NULL REFERENCES projects(id) ON DELETE CASCADE,
    commit_sha  VARCHAR(40),
    commit_msg  TEXT,
    image_tag   VARCHAR(255) NOT NULL,
    status      VARCHAR(30) NOT NULL DEFAULT 'queued',
    build_logs  TEXT,
    deployed_at TIMESTAMPTZ,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_deployments_project_id ON deployments(project_id);
