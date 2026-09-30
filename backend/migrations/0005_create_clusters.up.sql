CREATE TABLE IF NOT EXISTS clusters (
    id             VARCHAR(100) PRIMARY KEY,
    name           VARCHAR(100) NOT NULL,
    provider       VARCHAR(30) NOT NULL DEFAULT 'k3s',
    status         VARCHAR(30) NOT NULL DEFAULT 'provisioning',
    endpoint       VARCHAR(255) NOT NULL DEFAULT '',
    region         VARCHAR(100) NOT NULL DEFAULT '',
    zone           VARCHAR(100) NOT NULL DEFAULT '',
    credential_ref VARCHAR(255) NOT NULL,
    node_count     INT NOT NULL DEFAULT 0,
    created_at     TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at     TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
