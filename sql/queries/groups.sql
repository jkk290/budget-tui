-- name: CreateGroup :one
INSERT INTO groups (id, group_name, created_at, updated_at, user_id)
VALUES (
    $1,
    $2,
    $3,
    $4,
    $5
)
RETURNING *;

-- name: GetGroupsByUser :many
SELECT * FROM groups
WHERE user_id = $1;

-- name: GetGroupByID :one
SELECT * FROM groups
WHERE id = $1;

-- name: UpdateGroup :one
UPDATE groups
SET group_name = $2,
updated_at = NOW()
WHERE id = $1
RETURNING *;

-- name: DeleteGroup :exec
DELETE FROM groups
WHERE id = $1;