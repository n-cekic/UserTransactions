-- Create the database
CREATE DATABASE transactions;

-- Connect to the newly created database
\c transactions;

-- Crea the 'user_balance' table
CREATE TABLE "user_balance" (
    user_id SERIAL PRIMARY KEY,
    balance VARCHAR(255) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);