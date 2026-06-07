-- name: CreateUser :one
INSERT INTO users (
    username,
    password_hash,
    role
)
VALUES (
    $1,
    $2,
    $3
)
RETURNING
    id,
    username,
    password_hash,
    role,
    created_at;

-- name: GetUserByUsername :one
SELECT
    id,
    username,
    password_hash,
    role,
    created_at
FROM users
WHERE username = $1;

-- name: GetUserByID :one
SELECT
    id,
    username,
    password_hash,
    role,
    created_at
FROM users
WHERE id = $1;

-- name: ListUsers :many
SELECT
    id,
    username,
    role,
    created_at
FROM users
ORDER BY created_at DESC
LIMIT $1
OFFSET $2;

-- name: UpdateUserRole :one
UPDATE users
SET role = $2
WHERE id = $1
RETURNING
    id,
    username,
    password_hash,
    role,
    created_at;

-- name: DeleteUser :exec
DELETE FROM users
WHERE id = $1;