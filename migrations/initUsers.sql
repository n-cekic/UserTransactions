-- Create the database
CREATE DATABASE users;

-- Connect to the newly created database
\c users;

-- Create the 'user' table
CREATE TABLE "user" (
    user_id SERIAL PRIMARY KEY,
    email VARCHAR(255) UNIQUE NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);