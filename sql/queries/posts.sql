-- name: CreatePost :one
INSERT INTO posts (id, created_at, updated_at, body, user_id)
VALUES (
    gen_random_uuid(),
    Now(),
    Now(),
    $1,
    $2
)
RETURNING *;

-- name: GetPosts :many
SELECT * FROM posts
ORDER BY created_at ASC;

-- name: GetPostsOfAuthor :many
SELECT * FROM posts
WHERE user_id = $1
ORDER BY created_at ASC;

-- name: GetPost :one
SELECT * FROM posts
WHERE id = $1;

-- name: DeletePost :exec
DELETE FROM posts
WHERE id = $1;