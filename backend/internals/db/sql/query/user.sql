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

-- name: UpdateUsersEmail_verification :one
UPDATE users
    set email_verified = $1 WHERE id = $2 RETURNING id;

-- name: UpdateUserPasswordByUserId :one
UPDATE users
    set password_hash = $1 RETURNING id;

-- name: CreateLibrarian :one
UPDATE users
    SET role = 'LIBRARIAN' WHERE id = $1 RETURNING id;