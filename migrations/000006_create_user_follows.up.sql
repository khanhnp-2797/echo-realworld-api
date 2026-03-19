CREATE TABLE IF NOT EXISTS user_follows (
    follower_id BIGINT NOT NULL REFERENCES users (id) ON DELETE CASCADE,
    followed_id BIGINT NOT NULL REFERENCES users (id) ON DELETE CASCADE,
    PRIMARY KEY (follower_id, followed_id)
);

CREATE INDEX IF NOT EXISTS idx_user_follows_follower ON user_follows (follower_id);
CREATE INDEX IF NOT EXISTS idx_user_follows_followed ON user_follows (followed_id);
