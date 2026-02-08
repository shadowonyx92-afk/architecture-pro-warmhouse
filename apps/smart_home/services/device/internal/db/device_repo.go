package db

import (
    "context"
    "fmt"
    "time"
    "device/internal/models"
)

func (db *DB) CreateDevice(d *models.Device) error {
    query := `
        INSERT INTO devices (type_id, module_id, name, serial_number, enabled, location, remote_access, created_at, updated_at)
        VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$8)
        RETURNING id
    `
    d.CreatedAt = time.Now()
    d.UpdatedAt = d.CreatedAt
    return db.Pool.QueryRow(context.Background(), query,
        d.TypeID, d.ModuleID, d.Name, d.SerialNumber, d.Enabled, d.Location, d.RemoteAccess, d.CreatedAt,
    ).Scan(&d.ID)
}

func (db *DB) UpdateDevice(d *models.Device) error {
    d.UpdatedAt = time.Now()
    query := `
        UPDATE devices
        SET type_id=$1,module_id=$2,name=$3,serial_number=$4,enabled=$5,location=$6,remote_access=$7,updated_at=$8
        WHERE id=$9
    `
    result, err := db.Pool.Exec(context.Background(), query,
        d.TypeID, d.ModuleID, d.Name, d.SerialNumber, d.Enabled, d.Location, d.RemoteAccess, d.UpdatedAt, d.ID,
    )
    if err != nil {
        return err
    }
    if result.RowsAffected() == 0 {
        return fmt.Errorf("device with id %d not found", d.ID)
    }
    return nil
}

func (db *DB) DeleteDevice(id int) error {
    result, err := db.Pool.Exec(context.Background(), "DELETE FROM devices WHERE id=$1", id)
    if err != nil {
        return err
    }
    if result.RowsAffected() == 0 {
        return fmt.Errorf("device with id %d not found", id)
    }
    return nil
}

func (db *DB) GetDevice(id int) (*models.Device, error) {
    var d models.Device
    query := `SELECT id,type_id,module_id,name,serial_number,enabled,location,remote_access,created_at,updated_at FROM devices WHERE id=$1`
    err := db.Pool.QueryRow(context.Background(), query, id).Scan(
        &d.ID, &d.TypeID, &d.ModuleID, &d.Name, &d.SerialNumber, &d.Enabled, &d.Location, &d.RemoteAccess, &d.CreatedAt, &d.UpdatedAt,
    )
    if err != nil {
        return nil, err
    }
    return &d, nil
}

func (db *DB) GetDeviceList() ([]models.Device, error) {
    query := `SELECT id,type_id,module_id,name,serial_number,enabled,location,remote_access,created_at,updated_at FROM devices ORDER BY id`
    rows, err := db.Pool.Query(context.Background(), query)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    devices := []models.Device{}
    for rows.Next() {
        var d models.Device
        if err := rows.Scan(&d.ID, &d.TypeID, &d.ModuleID, &d.Name, &d.SerialNumber, &d.Enabled, &d.Location, &d.RemoteAccess, &d.CreatedAt, &d.UpdatedAt); err != nil {
            return nil, err
        }
        devices = append(devices, d)
    }
    return devices, nil
}
