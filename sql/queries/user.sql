-- name: CreateUser :one
INSERT INTO users (name, phone)
VALUES ($1, $2)
RETURNING id, name, phone;


-- name: GetUser :one
SELECT * from users where phone = $1;


-- name: GetUserByID :one
SELECT * from users where id = $1;