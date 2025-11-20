-- name: AddAccount :one
INSERT INTO accounts (id, account_name, account_type, balance, created_at, updated_at, user_id)
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

-- name: GetAccountsByUserID :many
SELECT * FROM accounts
WHERE user_id = $1;