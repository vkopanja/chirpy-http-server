-- name: CreateRefreshToken :one
INSERT INTO refresh_tokens (token, user_id, expires_at, revoked_at, created_at, updated_at)
VALUES ($1,
        $2,
        $3,
        $4,
        $5,
        $6)
RETURNING *;

-- name: GetTokenForRefreshToken :one
SELECT *
FROM refresh_tokens
WHERE token = $1
  AND revoked_at IS NULL
  AND expires_at > NOW();

-- name: RevokeRefreshToken :exec
UPDATE refresh_tokens
SET revoked_at = NOW(),
    updated_at = NOW()
WHERE token = $1;
