CREATE TABLE IF NOT EXISTS articles (
    id          BIGSERIAL PRIMARY KEY,
    created_at  TIMESTAMPTZ,
    updated_at  TIMESTAMPTZ,
    deleted_at  TIMESTAMPTZ,
    slug        VARCHAR(255) NOT NULL,
    title       VARCHAR(500) NOT NULL,
    description VARCHAR(1000),
    body        TEXT         NOT NULL,
    author_id   BIGINT       NOT NULL REFERENCES users (id)
);

CREATE UNIQUE INDEX IF NOT EXISTS idx_articles_slug       ON articles (slug)      WHERE deleted_at IS NULL;
CREATE INDEX        IF NOT EXISTS idx_articles_author_id  ON articles (author_id);
CREATE INDEX        IF NOT EXISTS idx_articles_deleted_at ON articles (deleted_at);
