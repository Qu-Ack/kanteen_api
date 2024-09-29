-- name: CreateOrder :one
INSERT INTO orders (user_id, total, status, created_at)
VALUES ($1, $2, $3, NOW())
RETURNING *;


-- name: GetOrder :one
SELECT 
  orders.id as order_id, 
  orders.status as order_status, 
  orders.total as total, 
  u.name AS user_name, 
  u.phone AS user_mobile 
FROM 
  orders 
JOIN 
  users u ON orders.user_id = u.id 
WHERE 
  orders.id = $1;  -- Ensure 'orders.id' is used to avoid ambiguity


-- name: GetPendingOrders :many
SELECT 
  orders.id as order_id, 
  orders.status as order_status, 
  orders.total as total, 
  u.name AS user_name, 
  u.phone AS user_mobile 
FROM 
  orders 
JOIN 
  users u ON orders.user_id = u.id 
WHERE 
  orders.status = 'pending' 
ORDER BY 
  orders.created_at DESC;
