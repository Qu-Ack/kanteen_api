-- name: CreateCategory :one 
INSERT INTO category (name)
VALUES ($1)
RETURNING id AS "ID", name AS "Name";

-- name: GetCategory :one
SELECT 
    id AS "ID", 
    name AS "Name"
FROM category
WHERE id = $1;


-- name: GetCategories :many
SELECT * from category;

-- name: UpdateCategory :exec 
UPDATE category
SET name = $1
WHERE id = $2;

-- name: DeleteCategory :exec
DELETE FROM category
WHERE id = $1;
