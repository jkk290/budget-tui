-- +goose Up
ALTER TABLE categories
DROP CONSTRAINT fk_group_id;

ALTER TABLE categories
ADD CONSTRAINT fk_group_id
FOREIGN KEY (group_id)
REFERENCES groups(id)
ON DELETE SET NULL;

-- +goose Down
ALTER TABLE categories
DROP CONSTRAINT fk_group_id;

ALTER TABLE categories
ADD CONSTRAINT fk_group_id
FOREIGN KEY (group_id)
REFERENCES groups(id)
ON DELETE CASCADE;
