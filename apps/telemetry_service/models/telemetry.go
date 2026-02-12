package models

import "time"

type TelemetryEvent struct {
	ID        int       `json:"id"`
	DeviceID  string    `json:"device_id"`
	Type      string    `json:"type"`
	Value     float64   `json:"value"`
	Status    string    `json:"status"`
	Timestamp time.Time `json:"timestamp"`
}
