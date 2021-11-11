BEGIN;

CREATE TABLE IF NOT EXISTS users(
    id SERIAL PRIMARY KEY,
    username VARCHAR(255) NOT NULL UNIQUE,
    email citext NOT NULL UNIQUE,
    bio TEXT,
    image VARCHAR(255), 
    password_hash VARCHAR(255) NOT NULL,
    created_at timestamptz NOT NULL DEFAULT NOW(),
    updated_at timestamptz 
);

COMMIT;