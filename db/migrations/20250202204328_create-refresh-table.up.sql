CREATE TABLE IF NOT EXISTS refresh_tokens (
    user_id VARCHAR NOT NULL,
    token VARCHAR NOT NULL,
    PRIMARY KEY (user_id, token)
);