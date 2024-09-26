-- name: CreateUser :one
INSERT INTO users (name, phone)
VALUES ($1, $2)
RETURNING id, name, phone;