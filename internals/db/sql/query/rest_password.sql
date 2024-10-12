-- name: DeleteRestPasswordByUserId :exec
DELETE FROM rest_password WHERE user_id = $1;

-- name: SavePasswordRestToken :one
INSERT INTO rest_password(token_hash,user_id,expires_at) values($1,$2,$3) RETURNING *;

-- name: GetResetPasswordFromTokenHash :one
SELECT * FROM rest_password where token_hash = $1;