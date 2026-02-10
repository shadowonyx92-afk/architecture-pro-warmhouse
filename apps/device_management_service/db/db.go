package db

import (
	"context"
	"errors"
	"fmt"
	"time"

	"device_management_service/models"

	"github.com/jackc/pgx/v5/pgxpool"
)

// DB хранит подключение к базе
type DB struct {
	Pool *pgxpool.Pool
}

// New создаёт новое подключение к базе
func New(connString string) (*DB, error) {
	pool, err := pgxpool.New(context.Background(), connString)
	if err != nil {
		return nil, fmt.Errorf("unable to connect to database: %w", err)
	}

	// Проверка соединения
	if err := pool.Ping(context.Background()); err != nil {
		return nil, fmt.Errorf("unable to ping database: %w", err)
	}

	return &DB{Pool: pool}, nil
}

// Close закрывает подключение к базе
func (db *DB) Close() {
	if db.Pool != nil {
		db.Pool.Close()
	}
}

// GetDevices возвращает список всех устройств
func (db *DB) GetDevices(ctx context.Context) ([]models.Device, error) {
	query := `
		SELECT id, name, type, location, value, unit, status, last_updated, created_at
		FROM devices
		ORDER BY id
	`

	rows, err := db.Pool.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("error querying devices: %w", err)
	}
	defer rows.Close()

	var devices []models.Device
	for rows.Next() {
		var d models.Device
		err := rows.Scan(
			&d.ID,
			&d.Name,
			&d.Type,
			&d.Location,
			&d.Value,
			&d.Unit,
			&d.Status,
			&d.LastUpdated,
			&d.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("error scanning device row: %w", err)
		}
		devices = append(devices, d)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating device rows: %w", err)
	}

	return devices, nil
}

// GetDeviceByID возвращает устройство по ID
func (db *DB) GetDeviceByID(ctx context.Context, id int) (models.Device, error) {
	query := `
		SELECT id, name, type, location, value, unit, status, last_updated, created_at
		FROM devices
		WHERE id = $1
	`

	var d models.Device
	err := db.Pool.QueryRow(ctx, query, id).Scan(
		&d.ID,
		&d.Name,
		&d.Type,
		&d.Location,
		&d.Value,
		&d.Unit,
		&d.Status,
		&d.LastUpdated,
		&d.CreatedAt,
	)
	if err != nil {
		return models.Device{}, fmt.Errorf("error getting device by ID: %w", err)
	}

	return d, nil
}

// UpdateDeviceValue обновляет значение и статус устройства
func (db *DB) UpdateDeviceValue(ctx context.Context, id int, value *float64, status *string) error {
	query := `
		UPDATE devices
		SET value = $1, status = $2, last_updated = $3
		WHERE id = $4
	`

	result, err := db.Pool.Exec(ctx, query, value, status, time.Now(), id)
	if err != nil {
		return fmt.Errorf("error updating device value: %w", err)
	}

	if result.RowsAffected() == 0 {
		return errors.New("device not found")
	}

	return nil
}

func (db *DB) DeleteDevice(ctx context.Context, id int) error {
	query := `
		DELETE FROM devices WHERE id = $1
	`

	result, err := db.Pool.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("error delete device value: %w", err)
	}

	if result.RowsAffected() == 0 {
		return errors.New("device not found")
	}

	return nil
}

// Creaete
func (db *DB) CreateDevice(ctx context.Context, d models.DeviceCreate) (models.Device, error) {
	query := `
		INSERT INTO devices (name, type, location, unit, created_at, last_updated)
		VALUES ($1, $2, $3, $4, NOW(), NOW())
		RETURNING id, name, type, location, value, unit, status, last_updated, created_at
	`

	var device models.Device
	err := db.Pool.QueryRow(
		ctx,
		query,
		d.Name,
		d.Type,
		d.Location,
		d.Unit,
	).Scan(
		&device.ID,
		&device.Name,
		&device.Type,
		&device.Location,
		&device.Value,
		&device.Unit,
		&device.Status,
		&device.LastUpdated,
		&device.CreatedAt,
	)

	if err != nil {
		return models.Device{}, err
	}

	return device, nil
}
