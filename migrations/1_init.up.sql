CREATE TABLE IF NOT EXISTS users
(
    id INT PRIMARY KEY,
    username TEXT NOT NULL UNIQUE,
    pass_hash BYTEA NOT NULL
    is_active BOOLEAN
);

CREATE INDEX IF NOT EXISTS idx_username ON users (username);