package models

import "time"

type DeviceModule struct {
    ID           int       `db:"id" json:"id"`
    Name         string    `db:"name" json:"name"`
    Description  string    `db:"description" json:"description"`
    Vendor       string    `db:"vendor" json:"vendor"`
    Connectivity string    `db:"connectivity" json:"connectivity"` // Wi-Fi, Zigbee, Z-Wave
    CreatedAt    time.Time `db:"created_at" json:"created_at"`
    UpdatedAt    time.Time `db:"updated_at" json:"updated_at"`
}
