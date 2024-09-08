-- name: GetUserById :one
SELECT * FROM users WHERE id = $1;

-- name: SaveUser :one
INSERT INTO users (
    name,mobile_number
) VALUES ($1,$2) RETURNING *;