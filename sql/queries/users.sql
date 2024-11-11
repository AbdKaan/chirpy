-- name: CreateUser :one
INSERT INTO users (id, created_at, updated_at, email, hashed_password)
VALUES (
    gen_random_uuid(),
    Now(),
    Now(),
    $1,
    $2
)
RETURNING *;

-- name: DeleteUsers :exec
DELETE FROM users;

-- name: GetUserWithEmail :one
SELECT * FROM users
WHERE email = $1;

-- name: UpdateUserEmailAndPassword :one
UPDATE users SET email = $2, hashed_password = $3
WHERE id = $1
RETURNING *;

-- name: UpgradeIsChirpyRed :one
UPDATE users SET is_chirpy_red = TRUE
WHERE id = $1
RETURNING *;