package devices_client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

type LightDevice struct {
	ID     string `json:"id"`
	Name   string `json:"name"`
	Status string `json:"status"` // "on" или "off"
}

type LightClient struct {
	BaseURL string
}

// Get получает текущее состояние света
func (c *LightClient) Get(id string) (*LightDevice, error) {
	resp, err := http.Get(fmt.Sprintf("%s/light/%s", c.BaseURL, id))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("failed to get light device, status: %d", resp.StatusCode)
	}

	var device LightDevice
	if err := json.NewDecoder(resp.Body).Decode(&device); err != nil {
		return nil, err
	}

	return &device, nil
}

// Set устанавливает состояние света ("on"/"off")
func (c *LightClient) Set(id string, status string) error {
	payload := map[string]string{"status": status}
	body, _ := json.Marshal(payload)

	resp, err := http.Post(fmt.Sprintf("%s/light/%s/set", c.BaseURL, id), "application/json", bytes.NewReader(body))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return fmt.Errorf("failed to set light device, status: %d", resp.StatusCode)
	}

	return nil
}
