-- name: GetProductByID :one
select * from products 
where id = $1 limit 1;

-- name: GetProductsByIDs :many
select * from products
where id = any($1::int[]);

-- name: GetProducts :many
select * from products;

-- name: CreateProduct :one
insert into products (
  name, price, image, description, quantity
) values ($1, $2, $3, $4, $5)
returning *;

-- name: UpdateProduct :exec
update products set name = $1, price = $2, 
image = $3, description = $4, quantity = $5
where id = $6;

