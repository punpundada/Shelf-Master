
ALTER TABLE books 
    ADD ISBN_10 VARCHAR(25), 
    ADD ISBN_13 VARCHAR(25),
    ADD edition VARCHAR(250) NOT NULL, 
    ADD publisher VARCHAR(250) NOT NULL,
    ADD date_published TIMESTAMP NOT NULL;

ALTER TABLE user_books
    ADD returned boolean;

ALTER TABLE user_books
    ALTER COLUMN returned set default false;