CREATE TABLE IF NOT EXISTS telemetry (
    id SERIAL PRIMARY KEY,
    device_id TEXT NOT NULL,
    type TEXT NOT NULL,
    value TEXT NULL,
    status TEXT NULL,
    timestamp TIMESTAMP NOT NULL
);
