-- +goose Up
-- Ensure the UUID extension is enabled
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Step 1: Temporarily drop the foreign key constraint from the `orders` table
ALTER TABLE orders DROP CONSTRAINT IF EXISTS orders_user_id_fkey;

-- Step 2: Add a new UUID column to `users`
ALTER TABLE users ADD COLUMN new_id UUID DEFAULT uuid_generate_v4();

-- Step 3: Change the `user_id` column in `orders` to UUID to prepare for the data update
ALTER TABLE orders ALTER COLUMN user_id TYPE UUID USING (uuid_generate_v4());

-- Step 4: Update the `user_id` in the `orders` table to match the new UUIDs
UPDATE orders SET user_id = users.new_id
FROM users WHERE orders.user_id::text = users.id::text;

-- Step 5: Set the new UUID column as the primary key
ALTER TABLE users DROP CONSTRAINT users_pkey;
ALTER TABLE users ADD PRIMARY KEY (new_id);

-- Step 6: Drop the old `id` column from `users`
ALTER TABLE users DROP COLUMN id;

-- Step 7: Rename `new_id` to `id`
ALTER TABLE users RENAME COLUMN new_id TO id;

-- Step 8: Re-establish the foreign key constraint between `orders` and `users`
ALTER TABLE orders
ADD CONSTRAINT orders_user_id_fkey FOREIGN KEY (user_id) REFERENCES users(id);