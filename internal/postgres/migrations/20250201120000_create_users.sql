-- +goose Up
CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    username VARCHAR(255) UNIQUE NOT NULL,
    balance NUMERIC(15,2) DEFAULT 0 NOT NULL CHECK (balance >= 0),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

INSERT INTO users (username, balance) VALUES
    ('user1', 1000.00),
    ('user2', 1500.50),
    ('user3', 2000.75);

-- +goose Down
DROP TABLE IF EXISTS users CASCADE;