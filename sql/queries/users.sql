-- name: CreateUser :one
INSERT INTO users (id, email, created_at, updated_at, hashed_password)
VALUES ($1,
        $2,
        $3,
        $4,
        $5)
RETURNING *;

-- name: ClearUsers :exec
TRUNCATE users CASCADE;

-- name: GetUserByEmail :one
SELECT *
FROM users
WHERE email = $1;

-- name: GetUserByRefreshToken :one
SELECT u.*
FROM users u
         JOIN refresh_tokens rt ON rt.user_id = u.id
WHERE rt.user_id = $1
  AND rt.revoked_at IS NULL
  AND rt.expires_at > NOW();

-- name: GetUserByID :one
SELECT *
FROM users
WHERE id = $1;

-- name: UpdateUserByID :one
UPDATE users
SET email           = $1,
    hashed_password = $2,
    updated_at      = NOW()
WHERE id = $3
RETURNING *;

-- name: UpdateUserChirpyRedByID :one
UPDATE users
SET is_chirpy_red = $1,
    updated_at    = NOW()
WHERE id = $2
RETURNING *;
