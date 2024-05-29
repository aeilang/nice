-- name: CreateUser :one
insert into users (
  firstname, lastname, email, password
) values ($1, $2, $3, $4)
returning *;

-- name: GetUserByEmail :one
select * from users
where email = $1 limit 1;


-- name: GetUserByID :one
select * from users
where id = $1 limit 1;
