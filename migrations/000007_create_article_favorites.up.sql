CREATE TABLE IF NOT EXISTS article_favorites (
    user_id    BIGINT NOT NULL REFERENCES users    (id) ON DELETE CASCADE,
    article_id BIGINT NOT NULL REFERENCES articles (id) ON DELETE CASCADE,
    PRIMARY KEY (user_id, article_id)
);

CREATE INDEX IF NOT EXISTS idx_article_favorites_user    ON article_favorites (user_id);
CREATE INDEX IF NOT EXISTS idx_article_favorites_article ON article_favorites (article_id);
