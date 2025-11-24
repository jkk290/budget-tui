-- +goose Up
ALTER TABLE categories
ALTER COLUMN budget TYPE NUMERIC(12, 2);

ALTER TABLE transactions
ALTER COLUMN amount TYPE NUMERIC(12, 2);

-- +goose Down
ALTER TABLE categories
ALTER COLUMN budget TYPE REAL;

ALTER TABLE transactions
ALTER COLUMN amount TYPE REAL;