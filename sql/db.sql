-- Connect to the go_crud database.
\c go_crud;

-- Create tables.
CREATE TABLE todoList (
     id SERIAL PRIMARY KEY,
     name VARCHAR(100) NOT NULL UNIQUE,
     bio TEXT
  
);

CREATE TABLE grocerys (
    NumberItems SERIAL PRIMARY KEY,
    GroceryItem VARCHAR(200) NOT NULL,
    Price NUMERIC(10,2)

);