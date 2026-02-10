package devices_client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

type TemperatureDevice struct {
	ID     string `json:"id"`
	Name   string `json:"name"`
	Value  string `json:"value"`
	Status string `json:"status"`
}

type TemperatureClient struct {
	BaseURL string
}

// Get получает текущее значение температуры по ID
func (c *TemperatureClient) Get(id string) (*TemperatureDevice, error) {
	resp, err := http.Get(fmt.Sprintf("%s/temperature/%s", c.BaseURL, id))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("failed to get temperature device, status: %d", resp.StatusCode)
	}

	var device TemperatureDevice
	if err := json.NewDecoder(resp.Body).Decode(&device); err != nil {
		return nil, err
	}

	return &device, nil
}

// Set устанавливает новое значение температуры по ID
func (c *TemperatureClient) Set(id string, value float64) error {
	payload := map[string]float64{"value": value}
	body, _ := json.Marshal(payload)

	resp, err := http.Post(fmt.Sprintf("%s/temperature/%s/set", c.BaseURL, id), "application/json", bytes.NewReader(body))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return fmt.Errorf("failed to set temperature device, status: %d", resp.StatusCode)
	}

	return nil
}
