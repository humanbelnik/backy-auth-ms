CREATE TABLE IF NOT EXISTS users
(
    id        BIGINT PRIMARY KEY,
    email     TEXT NOT NULL UNIQUE,
    nickname TEXT NOT NULL UNIQUE,
    is_admin BOOLEAN NOT NULL DEFAULT FALSE,
    pass_hash bytea NOT NULL
);
