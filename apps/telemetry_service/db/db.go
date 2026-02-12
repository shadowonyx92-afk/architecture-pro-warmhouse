package db

import (
	"context"
	"fmt"
	"telemetry_service/models"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type DB struct {
	Pool *pgxpool.Pool
}

func New(connString string) (*DB, error) {
	pool, err := pgxpool.New(context.Background(), connString)
	if err != nil {
		return nil, fmt.Errorf("failed to connect: %w", err)
	}

	if err := pool.Ping(context.Background()); err != nil {
		return nil, fmt.Errorf("failed to ping: %w", err)
	}

	return &DB{Pool: pool}, nil
}

func (db *DB) Close() {
	if db.Pool != nil {
		db.Pool.Close()
	}
}

func (db *DB) CreateTelemetry(ctx context.Context, t models.TelemetryEvent) (models.TelemetryEvent, error) {
	query := `
		INSERT INTO telemetry (device_id, type, value, status, timestamp)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id
	`
	now := time.Now()
	var id int
	err := db.Pool.QueryRow(ctx, query, t.DeviceID, t.Type, t.Value, t.Status, now).Scan(&id)
	if err != nil {
		return models.TelemetryEvent{}, err
	}
	t.ID = id
	t.Timestamp = now
	return t, nil
}

func (db *DB) GetTelemetryByDeviceID(ctx context.Context, deviceID string) ([]models.TelemetryEvent, error) {
	query := `
		SELECT id, device_id, type, value, status, timestamp
		FROM telemetry
		WHERE device_id = $1
		ORDER BY timestamp
	`

	rows, err := db.Pool.Query(ctx, query, deviceID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var events []models.TelemetryEvent
	for rows.Next() {
		var t models.TelemetryEvent
		if err := rows.Scan(&t.ID, &t.DeviceID, &t.Type, &t.Value, &t.Status, &t.Timestamp); err != nil {
			return nil, err
		}
		events = append(events, t)
	}
	return events, nil
}
