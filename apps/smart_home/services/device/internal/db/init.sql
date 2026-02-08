CREATE DATABASE device_db;

CREATE TABLE device_modules (
    id SERIAL PRIMARY KEY,
    name TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT now()
);

CREATE TABLE devices (
    id SERIAL PRIMARY KEY,
    type_id INT NOT NULL,
    module_id INT REFERENCES device_modules(id),
    name TEXT NOT NULL,
    serial_number TEXT UNIQUE,
    enabled BOOLEAN DEFAULT FALSE,
    location TEXT,
    remote_access BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP DEFAULT now(),
    updated_at TIMESTAMP DEFAULT now()
);
