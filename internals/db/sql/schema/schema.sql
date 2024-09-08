CREATE TYPE role_type AS ENUM ('ADMIN', 'USER', 'LIBRARIAN','AUTHOR');

CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    name text NOT NULL,
    mobile_number text,
    role role_type DEFAULT 'USER',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT mobile_number_length_check
        CHECK (mobile_number IS NULL OR LENGTH(mobile_number) = 10)
);

CREATE TABLE IF NOT EXISTS sessions(
    id text PRIMARY KEY,
    user_id int NOT NULL,
    expires_at TIMESTAMP NOT NULL,
    fresh BOOLEAN DEFAULT true,

    CONSTRAINT session_user_fk
        FOREIGN KEY (user_id)
        REFERENCES users (id)
        ON DELETE CASCADE
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
        REFERENCES libraries (id),
    CONSTRAINT email_format_check
        CHECK (email IS NOT NULL AND email ~ '^[A-Za-z0-9._%+-]+@[A-Za-z0-9.-]+\.[A-Za-z]{2,}$')
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

CREATE TABLE IF NOT EXISTS book_inventory (
    id SERIAL PRIMARY KEY,
    book_id INT NOT NULL,
    total_quantity INT NOT NULL DEFAULT 0,
    available_quantity INT NOT NULL DEFAULT 0,
    
    CONSTRAINT book_inventory_fk
        FOREIGN KEY (book_id)
        REFERENCES books (id)
        ON DELETE CASCADE,

    CONSTRAINT available_quantity_check
        CHECK (available_quantity >= 0)
);

CREATE TABLE IF NOT EXISTS user_books (
    user_id INT NOT NULL,
    book_id INT NOT NULL,
    borrowed_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    due_date TIMESTAMP NOT NULL,

    PRIMARY KEY (user_id, book_id),

    CONSTRAINT user_books_user_fk
        FOREIGN KEY (user_id)
        REFERENCES users (id)
        ON DELETE CASCADE,

    CONSTRAINT user_books_book_fk
        FOREIGN KEY (book_id)
        REFERENCES books (id)
        ON DELETE CASCADE
);