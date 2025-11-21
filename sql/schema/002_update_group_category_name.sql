-- +goose Up
ALTER TABLE groups
RENAME COLUMN name TO group_name;

ALTER TABLE categories
RENAME COLUMN name TO category_name;

-- +goose Down
ALTER TABLE groups
RENAME COLUMN group_name TO name;

ALTER TABLE categories
RENAME COLUMN category_name TO name;
