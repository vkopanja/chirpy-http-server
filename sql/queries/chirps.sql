-- name: CreateChirp :one
INSERT INTO chirps (id, user_id, body, created_at, updated_at)
VALUES ($1,
        $2,
        $3,
        $4,
        $5)
RETURNING *;

-- name: GetChirps :many
SELECT *
FROM chirps
ORDER BY created_at;

-- name: GetChirpById :one
SELECT *
FROM chirps
WHERE id = $1;

-- name: DeleteChirp :exec
DELETE
FROM chirps
WHERE id = $1;

-- name: GetChirpsByUserId :many
SELECT *
FROM chirps
WHERE user_id = $1
ORDER BY created_at;
