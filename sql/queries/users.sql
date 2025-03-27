-- name: CreateUser :one
INSERT INTO users (id, created_at, updated_at, email, hashed_password)
VALUES (
    gen_random_uuid(),
    NOW(),
    NOW(),
    $1,
    $2
)
RETURNING id, email, is_chirpy_red, created_at, updated_at;

-- name: GetUserByEmail :one
SELECT *
FROM users
WHERE email = $1;

-- name: UpdateUserEmailAndPasswordByID :one
UPDATE users
SET email = $2,
    hashed_password = $3,
    updated_at = NOW()
WHERE id = $1
RETURNING id, email, is_chirpy_red, created_at, updated_at;

-- name: UpgradeUserMembership :exec
UPDATE users
SET is_chirpy_red = true,
    updated_at = NOW()
WHERE id = $1;
