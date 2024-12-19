CREATE TABLE IF NOT EXISTS reset_password (
    token_hash VARCHAR(100) UNIQUE,
    user_id INT NOT NULL,
    expires_at Date NOT NULL
);