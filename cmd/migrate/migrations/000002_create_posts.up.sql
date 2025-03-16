CREATE TABLE IF NOT EXISTS posts (
    id bigserial PRIMARY KEY,
    content text NOT NULL,
    title text NOT NULL,
    user_id bigint NOT NULL,
    tags TEXT[] NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
