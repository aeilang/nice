// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.26.0
// source: users.sql

package store

import (
	"context"
)

const createUser = `-- name: CreateUser :one
insert into users (
  name, email, password, role 
) values ($1, $2, $3, $4)
returning id, name, email, password, role, created_at, updated_at
`

type CreateUserParams struct {
	Name     string `db:"name" json:"name"`
	Email    string `db:"email" json:"email"`
	Password string `db:"password" json:"password"`
	Role     string `db:"role" json:"role"`
}

func (q *Queries) CreateUser(ctx context.Context, arg CreateUserParams) (User, error) {
	row := q.queryRow(ctx, q.createUserStmt, createUser,
		arg.Name,
		arg.Email,
		arg.Password,
		arg.Role,
	)
	var i User
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.Email,
		&i.Password,
		&i.Role,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const getUserByEmail = `-- name: GetUserByEmail :one
select id, name, email, password, role, created_at, updated_at from users
where email = $1 limit 1
`

func (q *Queries) GetUserByEmail(ctx context.Context, email string) (User, error) {
	row := q.queryRow(ctx, q.getUserByEmailStmt, getUserByEmail, email)
	var i User
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.Email,
		&i.Password,
		&i.Role,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const getUserByID = `-- name: GetUserByID :one
select id, name, email, password, role, created_at, updated_at from users
where id = $1 limit 1
`

func (q *Queries) GetUserByID(ctx context.Context, id int32) (User, error) {
	row := q.queryRow(ctx, q.getUserByIDStmt, getUserByID, id)
	var i User
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.Email,
		&i.Password,
		&i.Role,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const updatePasswordByEmail = `-- name: UpdatePasswordByEmail :exec
update users
set password = $1
where email = $2
`

type UpdatePasswordByEmailParams struct {
	Password string `db:"password" json:"password"`
	Email    string `db:"email" json:"email"`
}

func (q *Queries) UpdatePasswordByEmail(ctx context.Context, arg UpdatePasswordByEmailParams) error {
	_, err := q.exec(ctx, q.updatePasswordByEmailStmt, updatePasswordByEmail, arg.Password, arg.Email)
	return err
}