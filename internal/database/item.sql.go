// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.26.0
// source: item.sql

package database

import (
	"context"
	"database/sql"
)

const createItem = `-- name: CreateItem :one
INSERT INTO item (name, category_id, price, stock)
VALUES ($1, $2, $3, $4)
RETURNING id AS "ID", category_id AS "CategoryID", name AS "Name", price AS "Price", stock AS "Stock", created_at, updated_at
`

type CreateItemParams struct {
	Name       string
	CategoryID int32
	Price      int32
	Stock      int32
}

type CreateItemRow struct {
	ID         int32
	CategoryID int32
	Name       string
	Price      int32
	Stock      int32
	CreatedAt  sql.NullTime
	UpdatedAt  sql.NullTime
}

func (q *Queries) CreateItem(ctx context.Context, arg CreateItemParams) (CreateItemRow, error) {
	row := q.db.QueryRowContext(ctx, createItem,
		arg.Name,
		arg.CategoryID,
		arg.Price,
		arg.Stock,
	)
	var i CreateItemRow
	err := row.Scan(
		&i.ID,
		&i.CategoryID,
		&i.Name,
		&i.Price,
		&i.Stock,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const deleteItem = `-- name: DeleteItem :exec
DELETE FROM item
WHERE id = $1
`

func (q *Queries) DeleteItem(ctx context.Context, id int32) error {
	_, err := q.db.ExecContext(ctx, deleteItem, id)
	return err
}

const getItems = `-- name: GetItems :many
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
JOIN category c ON i.category_id = c.id
`

type GetItemsRow struct {
	ID           int32
	CategoryName string
	CategoryID   int32
	Name         string
	Price        int32
	Stock        int32
	CreatedAt    sql.NullTime
	UpdatedAt    sql.NullTime
}

func (q *Queries) GetItems(ctx context.Context) ([]GetItemsRow, error) {
	rows, err := q.db.QueryContext(ctx, getItems)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetItemsRow
	for rows.Next() {
		var i GetItemsRow
		if err := rows.Scan(
			&i.ID,
			&i.CategoryName,
			&i.CategoryID,
			&i.Name,
			&i.Price,
			&i.Stock,
			&i.CreatedAt,
			&i.UpdatedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const updateItem = `-- name: UpdateItem :exec
UPDATE item
SET name = $1, category_id = $2, price = $3, stock = $4, updated_at = CURRENT_TIMESTAMP
WHERE id = $5
`

type UpdateItemParams struct {
	Name       string
	CategoryID int32
	Price      int32
	Stock      int32
	ID         int32
}

func (q *Queries) UpdateItem(ctx context.Context, arg UpdateItemParams) error {
	_, err := q.db.ExecContext(ctx, updateItem,
		arg.Name,
		arg.CategoryID,
		arg.Price,
		arg.Stock,
		arg.ID,
	)
	return err
}
