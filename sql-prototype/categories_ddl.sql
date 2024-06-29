DROP TABLE IF EXISTS categories;

CREATE TABLE categories (
id int primary key,
name string not null
);

INSERT INTO categories (id, name) values (1, 'Home');
