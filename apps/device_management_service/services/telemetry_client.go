package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

type TelemetryEvent struct {
	DeviceID string  `json:"device_id"`
	Type     string  `json:"type"`
	Value    *string `json:"value,omitempty"`
	Status   *string `json:"status,omitempty"`
}

type TelemetryClient struct {
	BaseURL string
}

func (c *TelemetryClient) Send(event TelemetryEvent) error {
	data, _ := json.Marshal(event)
	fmt.Println("→ TELEMETRY POST:", string(data))
	resp, err := http.Post(c.BaseURL+"/telemetry", "application/json", bytes.NewBuffer(data))
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 300 {
		return fmt.Errorf("failed to send telemetry, status: %d", resp.StatusCode)
	}
	return nil
}
