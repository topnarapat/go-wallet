-- Sequence and defined type

-- Table Definition
CREATE TABLE IF NOT EXISTS wallets (
    wallet_id SERIAL PRIMARY KEY,
    balance FLOAT NOT NULL,
    wallet_status TEXT NOT NULL DEFAULT 'Active',
    created_at TIMESTAMP NOT NULL DEFAULT (now())
);

INSERT INTO wallets (balance, created_at) VALUES (1000, '2023-01-27T12:30:00Z');
INSERT INTO wallets (balance, created_at) VALUES (2000, '2023-01-27T12:30:00Z');
INSERT INTO wallets (balance, created_at) VALUES (3000, '2023-01-27T12:30:00Z');
INSERT INTO wallets (balance, created_at) VALUES (4000, '2023-01-27T12:30:00Z');
INSERT INTO wallets (balance, created_at) VALUES (5000, '2023-01-27T12:30:00Z');