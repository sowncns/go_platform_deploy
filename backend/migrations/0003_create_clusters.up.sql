CREATE TABLE IF NOT EXISTS clusters (
    id             TEXT PRIMARY KEY,
    name           VARCHAR(255) NOT NULL,
    provider       VARCHAR(50) NOT NULL,
    status         VARCHAR(50) NOT NULL,
    endpoint       TEXT NOT NULL DEFAULT '',
    region         VARCHAR(255) NOT NULL DEFAULT '',
    zone           VARCHAR(255) NOT NULL DEFAULT '',
    credential_ref TEXT NOT NULL DEFAULT '',
    node_count     INT NOT NULL DEFAULT 0,
    created_at     TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at     TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
