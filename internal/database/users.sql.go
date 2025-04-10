// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.28.0
// source: users.sql

package database

import (
	"context"
	"time"

	"github.com/google/uuid"
)

const createUser = `-- name: CreateUser :one
INSERT INTO users (id, created_at, updated_at, email, hashed_password)
VALUES (
    gen_random_uuid(),
    NOW(),
    NOW(),
    $1,
    $2
)
RETURNING id, email, is_chirpy_red, created_at, updated_at
`

type CreateUserParams struct {
	Email          string
	HashedPassword string
}

type CreateUserRow struct {
	ID          uuid.UUID
	Email       string
	IsChirpyRed bool
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

func (q *Queries) CreateUser(ctx context.Context, arg CreateUserParams) (CreateUserRow, error) {
	row := q.db.QueryRowContext(ctx, createUser, arg.Email, arg.HashedPassword)
	var i CreateUserRow
	err := row.Scan(
		&i.ID,
		&i.Email,
		&i.IsChirpyRed,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const getUserByEmail = `-- name: GetUserByEmail :one
SELECT id, created_at, updated_at, email, hashed_password, is_chirpy_red
FROM users
WHERE email = $1
`

func (q *Queries) GetUserByEmail(ctx context.Context, email string) (User, error) {
	row := q.db.QueryRowContext(ctx, getUserByEmail, email)
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

const updateUserEmailAndPasswordByID = `-- name: UpdateUserEmailAndPasswordByID :one
UPDATE users
SET email = $2,
    hashed_password = $3,
    updated_at = NOW()
WHERE id = $1
RETURNING id, email, is_chirpy_red, created_at, updated_at
`

type UpdateUserEmailAndPasswordByIDParams struct {
	ID             uuid.UUID
	Email          string
	HashedPassword string
}

type UpdateUserEmailAndPasswordByIDRow struct {
	ID          uuid.UUID
	Email       string
	IsChirpyRed bool
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

func (q *Queries) UpdateUserEmailAndPasswordByID(ctx context.Context, arg UpdateUserEmailAndPasswordByIDParams) (UpdateUserEmailAndPasswordByIDRow, error) {
	row := q.db.QueryRowContext(ctx, updateUserEmailAndPasswordByID, arg.ID, arg.Email, arg.HashedPassword)
	var i UpdateUserEmailAndPasswordByIDRow
	err := row.Scan(
		&i.ID,
		&i.Email,
		&i.IsChirpyRed,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const upgradeUserMembership = `-- name: UpgradeUserMembership :exec
UPDATE users
SET is_chirpy_red = true,
    updated_at = NOW()
WHERE id = $1
`

func (q *Queries) UpgradeUserMembership(ctx context.Context, id uuid.UUID) error {
	_, err := q.db.ExecContext(ctx, upgradeUserMembership, id)
	return err
}
