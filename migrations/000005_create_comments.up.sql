CREATE TABLE IF NOT EXISTS comments (
    id         BIGSERIAL PRIMARY KEY,
    created_at TIMESTAMPTZ,
    updated_at TIMESTAMPTZ,
    deleted_at TIMESTAMPTZ,
    body       TEXT   NOT NULL,
    author_id  BIGINT NOT NULL REFERENCES users    (id),
    article_id BIGINT NOT NULL REFERENCES articles (id) ON DELETE CASCADE
);

CREATE INDEX IF NOT EXISTS idx_comments_author_id  ON comments (author_id);
CREATE INDEX IF NOT EXISTS idx_comments_article_id ON comments (article_id);
CREATE INDEX IF NOT EXISTS idx_comments_deleted_at ON comments (deleted_at);
