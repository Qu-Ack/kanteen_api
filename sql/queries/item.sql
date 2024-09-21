-- name: CreateItem :one
INSERT INTO item (name, category_id, price, stock)
VALUES ($1, $2, $3, $4)
RETURNING id AS "ID", category_id AS "CategoryID", name AS "Name", price AS "Price", stock AS "Stock", created_at, updated_at;

-- name: GetItems :many
SELECT 
    i.id AS "ID", 
    c.name AS "CategoryName", 
    i.category_id AS "CategoryID", 
    i.name AS "Name", 
    i.price AS "Price", 
    i.stock AS "Stock", 
    i.created_at, 
    i.updated_at
FROM item i
JOIN category c ON i.category_id = c.id;


-- name: UpdateItem :exec
UPDATE item
SET name = $1, category_id = $2, price = $3, stock = $4, updated_at = CURRENT_TIMESTAMP
WHERE id = $5;

-- name: DeleteItem :exec
DELETE FROM item
WHERE id = $1;
