-- name: GetUserByEmail :one
select * FROM librarians WHERE email = $1;

-- name: GetLibrarianById :one
SELECT * FROM librarians WHERE user_id  = $1;