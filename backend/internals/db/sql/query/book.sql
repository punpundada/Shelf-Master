-- name: SaveBookQuery :one
insert into books(name,authorid,description,created_at,updated_at,isbn_10,isbn_13,edition,publisher,date_published)
    values($1,$2,$3,$4,$5,$6,$7,$8,$9,$10) RETURNING *;

-- name: SaveBookInventory :one
insert into  book_inventory (book_id,total_quantity,available_quantity,library_id)
    values($1,$2,$3,$4) RETURNING id;

-- name: GetBooksByUserId :many
select b.id,b.name,b.authorid,b.description,b.created_at,b.updated_at 
    from books as b join user_books as ub 
    on b.id = ub.book_id 
        where ub.user_id = $1
        order by ub.borrowed_at;

-- name: ReturnBook :one
UPDATE user_books as ub 
    set returned = true 
    where ub.user_id = $1 
    and ub.book_id = $2 RETURNING ub.returned;    