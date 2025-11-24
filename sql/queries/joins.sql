-- name: GetAccountBalance :one
SELECT SUM(transactions.amount) AS account_balance
FROM accounts
INNER JOIN transactions
ON transactions.account_id = accounts.id
WHERE accounts.id = $1;