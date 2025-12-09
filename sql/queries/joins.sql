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

-- name: GetUserTransactions :many
SELECT transactions.*,
accounts.account_name,
categories.category_name
FROM transactions
INNER JOIN accounts
ON accounts.id = transactions.account_id
INNER JOIN categories
ON categories.id = transactions.category_id
WHERE accounts.user_id = $1
ORDER BY transactions.tx_date DESC;

-- name: GetUserCategoriesDetailed :many
SELECT categories.*,
groups.group_name
FROM categories
INNER JOIN groups
ON groups.id = categories.group_id
WHERE categories.user_id = $1;
