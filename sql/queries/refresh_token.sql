-- name: CreateRefreshToken :one
INSERT INTO refresh_tokens (token, created_at, updated_at, user_id, expires_at)
VALUES (
    $1,
    Now(),
    Now(),
    $2,
    Now() + INTERVAL '60 days'
)
RETURNING *;

-- name: GetUserFromRefreshToken :one
SELECT users.* FROM users
JOIN refresh_tokens ON users.id = refresh_tokens.user_id
WHERE refresh_tokens.token = $1
AND revoked_at IS NULL
AND expires_at > NOW();

-- name: RevokeRefreshToken :one
UPDATE refresh_tokens SET revoked_at = Now(),
updated_at = Now()
WHERE token = $1
RETURNING *;