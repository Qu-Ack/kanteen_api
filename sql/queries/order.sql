-- name: CreateOrder :one
INSERT INTO orders (user_id, total, status, created_at)
VALUES ($1, $2, $3, NOW())
RETURNING *;


-- name: GetOrder :one
SELECT * from orders where id = $1; 

