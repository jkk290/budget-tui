-- +goose Up
ALTER TABLE accounts
DROP COLUMN balance;

-- +goose Down
ALTER TABLE accounts
ADD COLUMN balance;