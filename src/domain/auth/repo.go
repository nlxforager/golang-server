package auth

// Code generated by sqlc.
// versions:
//   sqlc v1.25.0
// source: query.sql

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
)

const getUser = `-- name: GetUser :one
select username, auth_mode from accounts where username = $1 and auth_mode = $2
`

type Repository struct {
	db sql.DB
}

type GetUserParams struct {
	Username       string
	HashedPassword string
}

type GetUserRow struct {
	Username string
	AuthMode string
}

func (q *Repository) GetUser(ctx context.Context, arg GetUserParams) (GetUserRow, error) {
	row := q.db.QueryRowContext(ctx, getUser, arg.Username, arg.HashedPassword)
	var i GetUserRow
	err := row.Scan(&i.Username, &i.AuthMode)
	return i, err
}

const createUser = `-- name: CreateUser :one
INSERT into accounts(username, hashed_password,auth_mode) VALUES($1,$2,$3) RETURNING username
`

type CreateUserParams struct {
	Username       string
	HashedPassword string
	AuthMode       AUTH_MODE
}

func (q *Repository) CreateUser(ctx context.Context, arg CreateUserParams) (string, error) {
	row := q.db.QueryRowContext(ctx, createUser, arg.Username, arg.HashedPassword, arg.AuthMode)
	var username string
	err := row.Scan(&username)
	return username, err
}

type UpdateUserAuthModeParams struct {
	AuthMode *string
	Email    *string
	Username string
}

func (q *Repository) UpdateUser(ctx context.Context, arg UpdateUserAuthModeParams) error {
	idx := 1

	var args []any
	var sb strings.Builder
	sb.WriteString("UPDATE accounts SET")
	if arg.AuthMode != nil {
		sb.WriteString(fmt.Sprintf(" auth_mode=$%d", idx))
		args = append(args, arg.AuthMode)
		idx += 1
	}
	if arg.Email != nil {
		sb.WriteString(fmt.Sprintf(" email=$%d", idx))
		args = append(args, arg.Email)
		idx += 1
	}
	if idx == 1 {
		return errors.New("no change set")
	}
	sb.WriteString(fmt.Sprintf(" WHERE username = $%d", idx))
	args = append(args, arg.Username)

	_, err := q.db.ExecContext(ctx, sb.String(), args...)
	return err
}

const getEmail = `-- name: GetEmail :many
SELECT email from accounts left join emails on emails.account_id = accounts.id where username = $1
`

func (q *Repository) GetEmail(ctx context.Context, username string) ([]sql.NullString, error) {
	rows, err := q.db.QueryContext(ctx, getEmail, username)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []sql.NullString
	for rows.Next() {
		var email sql.NullString
		if err := rows.Scan(&email); err != nil {
			return nil, err
		}
		items = append(items, email)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}
