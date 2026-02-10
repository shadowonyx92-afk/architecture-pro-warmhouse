CREATE TABLE IF NOT EXISTS devices (
    id SERIAL PRIMARY KEY,
    name TEXT NOT NULL,
    type TEXT NOT NULL,
    location TEXT,
    value TEXT DEFAULT NULL,
    unit TEXT,
    status TEXT,
    last_updated TIMESTAMP DEFAULT now(),
    created_at TIMESTAMP DEFAULT now()
);
