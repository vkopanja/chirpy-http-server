-- +goose Up
ALTER TABLE users
    ADD COLUMN IF NOT EXISTS is_chirpy_red BOOLEAN DEFAULT FALSE;
UPDATE users
SET is_chirpy_red = FALSE
WHERE is_chirpy_red IS NULL;
ALTER TABLE users
    ALTER COLUMN is_chirpy_red SET NOT NULL;

-- +goose Down
ALTER TABLE users
    DROP COLUMN IF EXISTS is_chirpy_red;
