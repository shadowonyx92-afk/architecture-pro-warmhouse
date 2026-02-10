CREATE TABLE IF NOT EXISTS telemetry (
    id SERIAL PRIMARY KEY,
    device_id TEXT NOT NULL,
    type TEXT NOT NULL,
    value DOUBLE PRECISION NOT NULL,
    status TEXT NOT NULL,
    timestamp TIMESTAMP NOT NULL
);
