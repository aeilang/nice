-- name: CreateOrder :one
insert into orders (user_id, total, status, address)
values ($1, $2, $3, $4)
returning *;

-- name: CreateOrderItem :one
insert into order_items (order_id, product_id, quantity, price)
values ($1, $2, $3, $4)
returning *;