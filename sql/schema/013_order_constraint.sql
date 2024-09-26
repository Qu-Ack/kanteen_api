-- +goose Up
ALTER TABLE orders
    ALTER COLUMN user_id SET NOT NULL,
    ALTER COLUMN total SET NOT NULL,
    ALTER COLUMN status SET NOT NULL;

ALTER TABLE orderitems
    ALTER COLUMN order_id SET NOT NULL,
    ALTER COLUMN item_id SET NOT NULL,
    ALTER COLUMN takeaway_quantity SET NOT NULL,
    ALTER COLUMN eatin_quantity SET NOT NULL,
    ALTER COLUMN price SET NOT NULL;