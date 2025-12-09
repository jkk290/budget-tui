-- name: AddTransaction :one
INSERT INTO transactions (id, amount, tx_description, tx_date, created_at, updated_at, posted, account_id, category_id)
VALUES (
    $1,
    $2,
    $3,
    $4,
    $5,
    $6,
    $7,
    $8,
    $9
)
RETURNING *;

-- name: GetTransactionByID :one
SELECT * FROM transactions
WHERE id = $1;

-- name: GetTransactionsByAccount :many
SELECT * FROM transactions
WHERE account_id = $1
ORDER BY tx_date::date DESC, tx_date DESC;

-- name: GetTransactionsByCategory :many
SELECT * FROM transactions
WHERE category_id = $1
ORDER BY tx_date::date DESC, tx_date DESC;

-- name: UpdateTransaction :one
UPDATE transactions
SET amount = $2,
tx_description = $3,
tx_date = $4,
updated_at = NOW(),
posted = $5,
account_id = $6,
category_id = $7
WHERE id = $1
RETURNING *;

-- name: DeleteTransaction :exec
DELETE FROM transactions
WHERE id = $1;
