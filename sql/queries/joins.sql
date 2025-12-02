-- name: GetAccountBalance :one
SELECT SUM(transactions.amount) AS account_balance
FROM accounts
INNER JOIN transactions
ON transactions.account_id = accounts.id
WHERE accounts.id = $1;

-- name: GetUserAccountsBalances :many
SELECT accounts.*, (SUM(transactions.amount * 100))::bigint AS account_balance_cents
FROM accounts
INNER JOIN transactions
ON transactions.account_id = accounts.id
WHERE accounts.user_id = $1
GROUP BY accounts.id;

-- name: GetTransactionUserID :one
SELECT accounts.user_id AS user_id FROM transactions
INNER JOIN accounts
ON accounts.id = transactions.account_id
WHERE transactions.id = $1;
