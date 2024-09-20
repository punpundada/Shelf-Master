ALTER table users ADD COLUMN email_verified BOOLEAN DEFAULT false;


CREATE TABLE email_verification (
    id SERIAL PRIMARY KEY,
    code VARCHAR(15) NOT NULL,
    user_id INT NOT NULL,
    email text NOT NULL,
    expires_at Date NOT NULL,

    CONSTRAINT email_verification_users_fk
        FOREIGN KEY (user_id)
        REFERENCES users (id)
);

ALTER TABLE users
ADD CONSTRAINT email_unique UNIQUE (email);

ALTER TABLE users DROP COLUMN library_id;

ALTER TABLE book_inventory ADD COLUMN library_id INT NOT NULL;