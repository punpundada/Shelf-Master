-- name: DeleteEmailVerificationByUserId :one
DELETE FROM email_verification WHERE user_id = $1 RETURNING id;

-- name: SaveEmailVerification :one
INSERT INTO email_verification (
    code,user_id,email,expires_at
) values ($1,$2,$3,$4) RETURNING *;

-- name: GetEmailVerificationByUserId :one
SELECT * FROM email_verification where user_id = $1;

-- name: DeleteFromEmailVerificationByUserId :one
DELETE FROM email_verification where user_id = $1 RETURNING id;