CREATE TYPE role_type AS ENUM ('ADMIN', 'USER', 'LIBRARIAN','AUTHOR');

CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    name text NOT NULL,
    mobile_number text,
    role role_type DEFAULT 'USER',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS librarians(
    email text PRIMARY KEY,
    user_id INT NOT NULL,
    password_hash text NOT NULL,
    library_id int NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT librarian_user_FK
        FOREIGN KEY (user_id)
        REFERENCES users (id),

    CONSTRAINT library_librarian_fk
        FOREIGN KEY (library_id)
        REFERENCES libraries (id)
);


CREATE TABLE IF NOT EXISTS libraries(
    id SERIAL PRIMARY KEY,
    name text NOT NULL,
    address TEXT NOT NULL,
    phone_number VARCHAR(20),
    email VARCHAR(100),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);


CREATE TABLE IF NOT EXISTS books(
    id SERIAL PRIMARY KEY,
    name text NOT NULL,
    authorId INT NOT NULL,
    description text NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT book_author
        FOREIGN KEY (authorId)
        REFERENCES users (id)
        ON DELETE SET NULL
);