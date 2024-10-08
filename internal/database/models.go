// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.26.0

package database

import (
	"database/sql"

	"github.com/google/uuid"
)

type Category struct {
	ID   int32
	Name string
}

type Item struct {
	ID         int32
	Name       string
	CategoryID int32
	Price      int32
	Stock      int32
	CreatedAt  sql.NullTime
	UpdatedAt  sql.NullTime
}

type Order struct {
	ID        int32
	UserID    uuid.UUID
	Total     string
	Status    string
	CreatedAt sql.NullTime
}

type Orderitem struct {
	ID               int32
	OrderID          int32
	ItemID           int32
	TakeawayQuantity int32
	EatinQuantity    int32
	Price            string
}

type Otp struct {
	ID     uuid.UUID
	Mobile string
	Otp    string
}

type User struct {
	Name  string
	Phone string
	ID    uuid.NullUUID
}
