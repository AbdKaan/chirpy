// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0
// source: users.sql

package database

import (
	"context"

	"github.com/google/uuid"
)

const createUser = `-- name: CreateUser :one
INSERT INTO users (id, created_at, updated_at, email, hashed_password)
VALUES (
    gen_random_uuid(),
    Now(),
    Now(),
    $1,
    $2
)
RETURNING id, created_at, updated_at, email, hashed_password, is_chirpy_red
`

type CreateUserParams struct {
	Email          string
	HashedPassword string
}

func (q *Queries) CreateUser(ctx context.Context, arg CreateUserParams) (User, error) {
	row := q.db.QueryRowContext(ctx, createUser, arg.Email, arg.HashedPassword)
	var i User
	err := row.Scan(
		&i.ID,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.Email,
		&i.HashedPassword,
		&i.IsChirpyRed,
	)
	return i, err
}

const deleteUsers = `-- name: DeleteUsers :exec
DELETE FROM users
`

func (q *Queries) DeleteUsers(ctx context.Context) error {
	_, err := q.db.ExecContext(ctx, deleteUsers)
	return err
}

const getUserWithEmail = `-- name: GetUserWithEmail :one
SELECT id, created_at, updated_at, email, hashed_password, is_chirpy_red FROM users
WHERE email = $1
`

func (q *Queries) GetUserWithEmail(ctx context.Context, email string) (User, error) {
	row := q.db.QueryRowContext(ctx, getUserWithEmail, email)
	var i User
	err := row.Scan(
		&i.ID,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.Email,
		&i.HashedPassword,
		&i.IsChirpyRed,
	)
	return i, err
}

const updateUserEmailAndPassword = `-- name: UpdateUserEmailAndPassword :one
UPDATE users SET email = $2, hashed_password = $3
WHERE id = $1
RETURNING id, created_at, updated_at, email, hashed_password, is_chirpy_red
`

type UpdateUserEmailAndPasswordParams struct {
	ID             uuid.UUID
	Email          string
	HashedPassword string
}

func (q *Queries) UpdateUserEmailAndPassword(ctx context.Context, arg UpdateUserEmailAndPasswordParams) (User, error) {
	row := q.db.QueryRowContext(ctx, updateUserEmailAndPassword, arg.ID, arg.Email, arg.HashedPassword)
	var i User
	err := row.Scan(
		&i.ID,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.Email,
		&i.HashedPassword,
		&i.IsChirpyRed,
	)
	return i, err
}

const upgradeIsChirpyRed = `-- name: UpgradeIsChirpyRed :one
UPDATE users SET is_chirpy_red = TRUE
WHERE id = $1
RETURNING id, created_at, updated_at, email, hashed_password, is_chirpy_red
`

func (q *Queries) UpgradeIsChirpyRed(ctx context.Context, id uuid.UUID) (User, error) {
	row := q.db.QueryRowContext(ctx, upgradeIsChirpyRed, id)
	var i User
	err := row.Scan(
		&i.ID,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.Email,
		&i.HashedPassword,
		&i.IsChirpyRed,
	)
	return i, err
}
