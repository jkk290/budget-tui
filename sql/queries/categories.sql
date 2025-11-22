-- name: CreateCategory :one
INSERT INTO categories (id, category_name, created_at, updated_at, budget, user_id, group_id)
VALUES (
    $1,
    $2,
    $3,
    $4,
    $5,
    $6,
    $7
)
RETURNING *;

-- name: GetCategoriesByUser :many
SELECT * FROM categories
WHERE user_id = $1;

-- name: GetCategoryByID :one
SELECT * FROM categories
WHERE id = $1;

-- name: UpdateCategory :one
UPDATE categories
SET category_name = $2,
budget = $3,
group_id = $4,
updated_at = NOW()
WHERE id = $1
RETURNING *;

-- name: DeleteCategory :exec
DELETE FROM categories
WHERE id = $1;