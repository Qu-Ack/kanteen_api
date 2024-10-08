// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.26.0
// source: user.sql

package database

import (
	"context"

	"github.com/google/uuid"
)

const createUser = `-- name: CreateUser :one
INSERT INTO users (name, phone)
VALUES ($1, $2)
RETURNING id, name, phone
`

type CreateUserParams struct {
	Name  string
	Phone string
}

type CreateUserRow struct {
	ID    uuid.NullUUID
	Name  string
	Phone string
}

func (q *Queries) CreateUser(ctx context.Context, arg CreateUserParams) (CreateUserRow, error) {
	row := q.db.QueryRowContext(ctx, createUser, arg.Name, arg.Phone)
	var i CreateUserRow
	err := row.Scan(&i.ID, &i.Name, &i.Phone)
	return i, err
}

const getUser = `-- name: GetUser :one
SELECT name, phone, id from users where phone = $1
`

func (q *Queries) GetUser(ctx context.Context, phone string) (User, error) {
	row := q.db.QueryRowContext(ctx, getUser, phone)
	var i User
	err := row.Scan(&i.Name, &i.Phone, &i.ID)
	return i, err
}

const getUserByID = `-- name: GetUserByID :one
SELECT name, phone, id from users where id = $1
`

func (q *Queries) GetUserByID(ctx context.Context, id uuid.NullUUID) (User, error) {
	row := q.db.QueryRowContext(ctx, getUserByID, id)
	var i User
	err := row.Scan(&i.Name, &i.Phone, &i.ID)
	return i, err
}
