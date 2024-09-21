-- +goose Up
-- Create category table
CREATE TABLE IF NOT EXISTS category (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL UNIQUE
);

-- Create item table
CREATE TABLE IF NOT EXISTS item (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    category_id INTEGER NOT NULL REFERENCES category(id) ON DELETE CASCADE,
    price INTEGER NOT NULL,
    stock INTEGER NOT NULL,
    created_at TIMESTAMP WITHOUT TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITHOUT TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Create users table
CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255),
    email VARCHAR(255) UNIQUE,
    password VARCHAR(255)
);

-- Create orders table
CREATE TABLE IF NOT EXISTS orders (
    id SERIAL PRIMARY KEY,
    user_id INTEGER REFERENCES users(id),
    total NUMERIC(10, 2),
    status VARCHAR(50),
    created_at TIMESTAMP WITHOUT TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Create orderitems table
CREATE TABLE IF NOT EXISTS orderitems (
    id SERIAL PRIMARY KEY,
    order_id INTEGER REFERENCES orders(id),
    item_id INTEGER REFERENCES item(id),
    takeaway_quantity INTEGER,
    eatin_quantity INTEGER,
    price NUMERIC(10, 2)
);

-- +goose Down
-- Drop orderitems table
DROP TABLE IF EXISTS orderitems;

-- Drop orders table
DROP TABLE IF EXISTS orders;

-- Drop users table
DROP TABLE IF EXISTS users;

-- Drop item table
DROP TABLE IF EXISTS item;

-- Drop category table
DROP TABLE IF EXISTS category;
