-- name: CreateUser :one
INSERT INTO users (id, created_at, updated_at, username, hashed_pw) 
VALUES (
    $1,
    $2,
    $3,
    $4,
    $5
)
RETURNING *;

-- name: GetUserByUsername :one
SELECT * FROM users
WHERE username = $1;