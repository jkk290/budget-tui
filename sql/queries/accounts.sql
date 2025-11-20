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

-- name: GetAccountByID :one
SELECT * FROM accounts
WHERE id = $1;

-- name: GetAccountsByUserID :many
SELECT * FROM accounts
WHERE user_id = $1;

-- name: UpdateAccountBalance :one
UPDATE accounts
SET balance = $2,
updated_at = NOW()
WHERE id = $1
RETURNING *;

-- name: UpdateAccountInfo :one
UPDATE accounts
SET account_name = $2,
updated_at = NOW()
where id = $1
RETURNING *;

-- name: DeleteAccount :exec
DELETE FROM accounts 
WHERE id = $1;