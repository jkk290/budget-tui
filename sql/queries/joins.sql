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
LEFT JOIN groups
ON groups.id = categories.group_id
WHERE categories.user_id = $1;

-- name: GetUserBudgetOverviewForMonth :many
SELECT categories.id AS category_id,
categories.category_name,
categories.budget,
categories.group_id,
groups.group_name,
COALESCE(SUM(-transactions.amount), 0)::numeric AS total_spent
FROM categories
LEFT JOIN groups ON groups.id = categories.group_id
LEFT JOIN transactions
ON transactions.category_id = categories.id
AND transactions.tx_date >= $2
AND transactions.tx_date < $3
WHERE categories.user_id = $1
GROUP BY categories.id, categories.category_name, categories.budget,
categories.group_id, groups.group_name
ORDER BY groups.group_name NULLS LAST, categories.category_name;
