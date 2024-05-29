// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.26.0
// source: order.sql

package store

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"
)

const createOrder = `-- name: CreateOrder :one
insert into orders (user_id, total, status, address)
values ($1, $2, $3, $4)
returning id, user_id, total, status, address, created_at
`

type CreateOrderParams struct {
	UserID  int32          `db:"user_id" json:"user_id"`
	Total   pgtype.Numeric `db:"total" json:"total"`
	Status  string         `db:"status" json:"status"`
	Address string         `db:"address" json:"address"`
}

func (q *Queries) CreateOrder(ctx context.Context, arg CreateOrderParams) (Order, error) {
	row := q.db.QueryRow(ctx, createOrder,
		arg.UserID,
		arg.Total,
		arg.Status,
		arg.Address,
	)
	var i Order
	err := row.Scan(
		&i.ID,
		&i.UserID,
		&i.Total,
		&i.Status,
		&i.Address,
		&i.CreatedAt,
	)
	return i, err
}

const createOrderItem = `-- name: CreateOrderItem :one
insert into order_items (order_id, product_id, quantity, price)
values ($1, $2, $3, $4)
returning id, order_id, product_id, quantity, price
`

type CreateOrderItemParams struct {
	OrderID   int32          `db:"order_id" json:"order_id"`
	ProductID int32          `db:"product_id" json:"product_id"`
	Quantity  int32          `db:"quantity" json:"quantity"`
	Price     pgtype.Numeric `db:"price" json:"price"`
}

func (q *Queries) CreateOrderItem(ctx context.Context, arg CreateOrderItemParams) (OrderItem, error) {
	row := q.db.QueryRow(ctx, createOrderItem,
		arg.OrderID,
		arg.ProductID,
		arg.Quantity,
		arg.Price,
	)
	var i OrderItem
	err := row.Scan(
		&i.ID,
		&i.OrderID,
		&i.ProductID,
		&i.Quantity,
		&i.Price,
	)
	return i, err
}
