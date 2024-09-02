CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    name text NOT NULL,
    email text NOT NULL,
    password_hash text NOT NULL
);

CREATE TABLE IF NOT EXISTS author(
    id serial PRIMARY KEY,
    name text NOT NULL
);


CREATE TABLE IF NOT EXISTS books (
    id SERIAL PRIMARY KEY,
    name text NOT NULL,
    authorId INT NOT NULL,
    description text NOT NULL,
    CONSTRAINT book_author
    FOREIGN KEY (authorId)
    REFERENCES author(id)
    ON DELETE SET NULL
);