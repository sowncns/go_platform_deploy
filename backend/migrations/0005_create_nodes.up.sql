CREATE TABLE IF NOT EXISTS nodes (
    id              TEXT PRIMARY KEY,
    cluster_id      TEXT NOT NULL REFERENCES clusters(id) ON DELETE CASCADE,
    name            VARCHAR(255) NOT NULL,
    role            VARCHAR(50) NOT NULL,
    status          VARCHAR(50) NOT NULL,
    cpu_capacity    BIGINT NOT NULL DEFAULT 0,
    memory_capacity BIGINT NOT NULL DEFAULT 0,
    pod_capacity    INT NOT NULL DEFAULT 0,
    cpu_usage       BIGINT NOT NULL DEFAULT 0,
    memory_usage    BIGINT NOT NULL DEFAULT 0,
    pod_usage       INT NOT NULL DEFAULT 0,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_nodes_cluster_id ON nodes(cluster_id);
