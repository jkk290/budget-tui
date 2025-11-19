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

-- name: UpdateUser :one
UPDATE users 
SET created_at = $2,
updated_at = $3,
username = $4,
hashed_pw = $5
WHERE id = $1
RETURNING *;