CREATE TABLE items (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100),
    category VARCHAR(100) NOT NULL,
    price DECIMAL(10, 2) NOT NULL
);