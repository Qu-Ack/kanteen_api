-- name: CreateOrderItem :one
INSERT INTO orderitems (order_id, item_id, takeaway_quantity, eatin_quantity, price)
VALUES ($1, $2, $3, $4, $5) RETURNING *;


-- name: GetOrderItemsForOrder :many
SELECT * FROM orderitems where order_id=$1;


