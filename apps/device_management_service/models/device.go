package models

type DeviceType string

const (
	Temperature DeviceType = "temperature"
	Light       DeviceType = "light"
	Gate        DeviceType = "gate"
)

type Device struct {
	ID     string     `json:"id"`
	Name   string     `json:"name"`
	Type   DeviceType `json:"type"`
	Status string     `json:"status"`
	Value  float64    `json:"value,omitempty"` // для температуры
}

type DeviceCommand struct {
	Command string   `json:"command"`         // "toggle", "open", "close"
	Value   *float64 `json:"value,omitempty"` // для установки значения, если нужно
}
