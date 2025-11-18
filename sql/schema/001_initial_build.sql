-- +goose Up
CREATE TABLE users (
    id UUID PRIMARY KEY,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    username TEXT NOT NULL UNIQUE,
    hashed_pw TEXT NOT NULL
);

CREATE TABLE groups (
    id UUID PRIMARY KEY,
    name TEXT NOT NULL,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    user_id UUID NOT NULL,
    CONSTRAINT fk_user_id
    FOREIGN KEY (user_id)
    REFERENCES users(id)
    ON DELETE CASCADE
);

CREATE TABLE categories (
    id UUID PRIMARY KEY,
    name TEXT NOT NULL,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    budget REAL NOT NULL DEFAULT (0.00),
    user_id UUID NOT NULL,
    group_id UUID,
    CONSTRAINT fk_user_id
    FOREIGN KEY (user_id)
    REFERENCES users(id)
    ON DELETE CASCADE,
    CONSTRAINT fk_group_id
    FOREIGN KEY (group_id)
    REFERENCES groups(id)
    ON DELETE CASCADE
);

CREATE TABLE accounts (
    id UUID PRIMARY KEY,
    account_name TEXT NOT NULL,
    account_type TEXT NOT NULL,
    balance REAL NOT NULL DEFAULT (0.00),
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    user_id UUID NOT NULL,
    CONSTRAINT fk_user_id,
    FOREIGN KEY (user_id)
    REFERENCES users(id)
    ON DELETE CASCADE
);

CREATE TABLE transactions (
    id UUID PRIMARY KEY,
    amount REAL NOT NULL,
    tx_description TEXT NOT NULL,
    tx_date TIMESTAMP NOT NULL,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    posted BOOLEAN NOT NULL DEFAULT (false),
    account_id UUID NOT NULL,
    category_id UUID NOT NULL,
    CONSTRAINT fk_account_id,
    FOREIGN KEY (account_id)
    REFERENCES accounts(id)
    ON DELETE CASCADE,
    CONSTRAINT fk_category_id,
    FOREIGN KEY (category_id)
    REFERENCES categories(id)
);

-- +goose Down
DROP TABLE users;
DROP TABLE groups;
DROP TABLE categories;
DROP TABLE accounts;
DROP TABLE transactions;
