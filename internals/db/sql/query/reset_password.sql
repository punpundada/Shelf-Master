-- name: DeleteRestPasswordByUserId :exec
DELETE FROM reset_password WHERE user_id = $1;

-- name: SavePasswordRestToken :one
INSERT INTO reset_password(token_hash,user_id,expires_at) values($1,$2,$3) RETURNING *;

-- name: GetResetPasswordFromTokenHash :one
SELECT * FROM reset_password where token_hash = $1;