CREATE TABLE IF NOT EXISTS tags (
    id         BIGSERIAL PRIMARY KEY,
    created_at TIMESTAMPTZ,
    updated_at TIMESTAMPTZ,
    deleted_at TIMESTAMPTZ,
    name       VARCHAR(100) NOT NULL
);

CREATE UNIQUE INDEX IF NOT EXISTS idx_tags_name       ON tags (name) WHERE deleted_at IS NULL;
CREATE INDEX        IF NOT EXISTS idx_tags_deleted_at ON tags (deleted_at);
