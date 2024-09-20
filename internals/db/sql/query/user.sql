-- name: GetUserById :one
SELECT * FROM users WHERE id = $1;

-- name: SaveUser :one
INSERT INTO users (
    name,mobile_number,email,password_hash
) VALUES ($1,$2,$3,$4) RETURNING *;

-- name: GetUserByEmail :one
SELECT * FROM users WHERE email = $1;

-- name: CreateAdmin :one
UPDATE users
    SET role = 'ADMIN' WHERE id = $1 RETURNING id;