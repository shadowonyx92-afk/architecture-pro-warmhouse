package models

import "time"

type DeviceType string

const (
	Temperature DeviceType = "temperature"
	Light       DeviceType = "light"
	Gate        DeviceType = "gate"
)

// Device представляет устройство в умном доме
type Device struct {
	ID          int        `json:"id"`
	Name        string     `json:"name"`
	Type        DeviceType `json:"type"`
	Location    string     `json:"location"`
	Value       *float64   `json:"value,omitempty"`
	Unit        *string    `json:"unit,omitempty"`
	Status      *string    `json:"status,omitempty"`
	LastUpdated *time.Time `json:"last_updated,omitempty"`
	CreatedAt   *time.Time `json:"created_at,omitempty"`
}

type DeviceCommand struct {
	Command string   `json:"command"`
	Value   *float64 `json:"value,omitempty"` // для установки значения, если нужно
}

type DeviceCreate struct {
	Name     string   `json:"name" binding:"required"`
	Type     string   `json:"type" binding:"required"`
	Location *string  `json:"location"`
	Unit     *string  `json:"unit"`
	Value    *float64 `json:"value"`
	Status   *string  `json:"status"`
}

type DeviceUpdate struct {
	Name     *string  `json:"name"`
	Location *string  `json:"location"`
	Unit     *string  `json:"unit"`
	Value    *float64 `json:"value"`
	Status   *string  `json:"status"`
}
