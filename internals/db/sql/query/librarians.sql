-- name: GetUserByEmail :one
select * FROM librarians WHERE email = $1;