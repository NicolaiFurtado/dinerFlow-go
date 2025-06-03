CREATE TABLE IF NOT EXISTS diner_tables (
    id SERIAL PRIMARY KEY,
    table_name VARCHAR(255) NOT NULL,
    seats INTEGER NOT NULL CHECK (seats > 0),
    status VARCHAR(50) NOT NULL DEFAULT 'available'
);