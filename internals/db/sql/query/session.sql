-- name: SaveSession :one
INSERT INTO sessions (
    user_id,expires_at
) VALUES ($1, $2) RETURNING *;

-- name: GetSessionById :one
SELECT * FROM sessions WHERE id = $1;

-- name: UpdateSessionById :one
UPDATE sessions SET expires_at = $1 , fresh = $2
WHERE id = $3
RETURNING *;