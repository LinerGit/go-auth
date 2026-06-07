-- name: CreateRefreshToken :one
INSERT INTO refresh_tokens (
    user_id,
    token_hash,
    expires_at
)
VALUES (
    $1,
    $2,
    $3
)
RETURNING
    id,
    user_id,
    token_hash,
    expires_at,
    created_at;

-- name: GetRefreshToken :one
SELECT
    id,
    user_id,
    token_hash,
    expires_at,
    created_at
FROM refresh_tokens
WHERE token_hash = $1;

-- name: DeleteRefreshToken :exec
DELETE FROM refresh_tokens
WHERE token_hash = $1;

-- name: DeleteUserRefreshTokens :exec
DELETE FROM refresh_tokens
WHERE user_id = $1;

-- name: DeleteExpiredRefreshTokens :exec
DELETE FROM refresh_tokens
WHERE expires_at < NOW();

-- name: GetUserRefreshTokens :many
SELECT
    id,
    user_id,
    expires_at,
    created_at
FROM refresh_tokens
WHERE user_id = $1;