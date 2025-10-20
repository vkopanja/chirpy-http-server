-- +goose Up
CREATE TABLE users
(
    id         UUID      NOT NULL PRIMARY KEY DEFAULT gen_random_uuid(),
    email      TEXT      NOT NULL UNIQUE,
    created_at TIMESTAMP NOT NULL             DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL
);

-- +goose Down
DROP TABLE users;
