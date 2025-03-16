CREATE TABLE IF NOT EXISTS users(
    id bigserial PRIMARY KEY,
    email varchar(255) UNIQUE NOT NULL,
    username varchar(255) NOT NULL,
    password bytea NOT NULL,
    created_at timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP
)