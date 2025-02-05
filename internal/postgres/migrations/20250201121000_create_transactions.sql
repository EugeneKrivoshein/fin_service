-- +goose Up
CREATE TABLE transactions (
    id SERIAL PRIMARY KEY,
    user_id INT REFERENCES users(id) ON DELETE CASCADE,  
    sender_id INT REFERENCES users(id),  
    receiver_id INT REFERENCES users(id), 
    amount NUMERIC(15,2) NOT NULL CHECK (amount > 0),
    transaction_type VARCHAR(20) NOT NULL CHECK (transaction_type IN ('deposit', 'transfer')),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_user_id ON transactions(user_id);
CREATE INDEX idx_sender_id ON transactions(sender_id);
CREATE INDEX idx_receiver_id ON transactions(receiver_id);

-- +goose Down
DROP TABLE IF EXISTS transactions CASCADE;