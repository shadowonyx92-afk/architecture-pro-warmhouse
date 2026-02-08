package models

import "time"

type TelemetryData struct {
    ID        int       `db:"id" json:"id"`
    DeviceID  int       `db:"device_id" json:"device_id"` // FK -> Device.id (из DeviceService)
    Value     string    `db:"value" json:"value"`         // любое значение: температура, on/off, locked/unlocked
    Unit      string    `db:"unit" json:"unit"`           // для числовых значений, например "C", "kWh"
    Status    string    `db:"status" json:"status"`       // состояние устройства (опционально)
    Timestamp time.Time `db:"timestamp" json:"timestamp"` // время события
}
