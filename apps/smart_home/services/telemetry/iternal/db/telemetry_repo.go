package db

import (
    "context"
    "fmt"
    "time"

    "telemetry/internal/models"

    "github.com/jackc/pgx/v5/pgxpool"
)

type DB struct {
    Pool *pgxpool.Pool
}

// New создает подключение к БД
func New(connString string) (*DB, error) {
    pool, err := pgxpool.New(context.Background(), connString)
    if err != nil {
        return nil, fmt.Errorf("unable to connect to database: %w", err)
    }

    if err := pool.Ping(context.Background()); err != nil {
        return nil, fmt.Errorf("unable to ping database: %w", err)
    }

    return &DB{Pool: pool}, nil
}

func (db *DB) Close() {
    if db.Pool != nil {
        db.Pool.Close()
    }
}

// CreateTelemetry записывает новую запись телеметрии
func (db *DB) CreateTelemetry(t *models.TelemetryData) error {
    query := `
        INSERT INTO telemetry_data (device_id, type, value, unit, timestamp)
        VALUES ($1, $2, $3, $4, $5)
        RETURNING id
    `
    if t.Timestamp.IsZero() {
        t.Timestamp = time.Now()
    }

    return db.Pool.QueryRow(context.Background(), query,
        t.DeviceID, t.Type, t.Value, t.Unit, t.Timestamp).Scan(&t.ID)
}

// GetTelemetryByDevice возвращает список телеметрий для устройства
func (db *DB) GetTelemetryByDevice(deviceID int) ([]models.TelemetryData, error) {
    query := `
        SELECT id, device_id, type, value, unit, timestamp
        FROM telemetry_data
        WHERE device_id = $1
        ORDER BY timestamp DESC
    `
    rows, err := db.Pool.Query(context.Background(), query, deviceID)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var list []models.TelemetryData
    for rows.Next() {
        var t models.TelemetryData
        if err := rows.Scan(&t.ID, &t.DeviceID, &t.Type, &t.Value, &t.Unit, &t.Timestamp); err != nil {
            return nil, err
        }
        list = append(list, t)
    }
    return list, nil
}

// GetLatestTelemetry возвращает последнюю запись телеметрии устройства
func (db *DB) GetLatestTelemetry(deviceID int, telemetryType string) (*models.TelemetryData, error) {
    query := `
        SELECT id, device_id, type, value, unit, timestamp
        FROM telemetry_data
        WHERE device_id = $1 AND type = $2
        ORDER BY timestamp DESC
        LIMIT 1
    `
    var t models.TelemetryData
    err := db.Pool.QueryRow(context.Background(), query, deviceID, telemetryType).Scan(
        &t.ID, &t.DeviceID, &t.Type, &t.Value, &t.Unit, &t.Timestamp,
    )
    if err != nil {
        return nil, err
    }
    return &t, nil
}

// DeleteTelemetry удаляет запись телеметрии по ID
func (db *DB) DeleteTelemetry(id int) error {
    query := `DELETE FROM telemetry_data WHERE id = $1`
    _, err := db.Pool.Exec(context.Background(), query, id)
    return err
}
