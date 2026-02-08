package models

import "time"

type Device struct {
    ID           int       `db:"id" json:"id"`
    TypeID       int       `db:"type_id" json:"type_id"`
    ModuleID     *int      `db:"module_id" json:"module_id,omitempty"`
    Name         string    `db:"name" json:"name"`
    SerialNumber string    `db:"serial_number" json:"serial_number"`
    Enabled      bool      `db:"enabled" json:"enabled"`
    Location     string    `db:"location" json:"location"`
    RemoteAccess bool      `db:"remote_access" json:"remote_access"`
    CreatedAt    time.Time `db:"created_at" json:"created_at"`
    UpdatedAt    time.Time `db:"updated_at" json:"updated_at"`
}
