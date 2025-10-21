-- +goose Up
CREATE TABLE chirps
(
    id         UUID      NOT NULL PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id    UUID      NOT NULL REFERENCES users (id) ON DELETE CASCADE,
    body       TEXT      NOT NULL,
    created_at TIMESTAMP NOT NULL             DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL
);

-- +goose Down
DROP TABLE chirps;
