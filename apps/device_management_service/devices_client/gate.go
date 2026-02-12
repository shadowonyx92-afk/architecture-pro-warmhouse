package devices_client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

type GateDevice struct {
	ID     string `json:"id"`
	Name   string `json:"name"`
	Status string `json:"status"` // "open" или "closed"
}

type GateClient struct {
	BaseURL string
}

// Get получает текущее состояние ворот
func (c *GateClient) Get(id string) (*GateDevice, error) {
	resp, err := http.Get(fmt.Sprintf("%s/gate/%s", c.BaseURL, id))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("failed to get gate device, status: %d", resp.StatusCode)
	}

	var device GateDevice
	if err := json.NewDecoder(resp.Body).Decode(&device); err != nil {
		return nil, err
	}

	return &device, nil
}

// Set устанавливает состояние ворот ("open"/"closed")
func (c *GateClient) Set(id string, status string) error {
	payload := map[string]string{"status": status}
	body, _ := json.Marshal(payload)

	resp, err := http.Post(fmt.Sprintf("%s/gate/%s/set", c.BaseURL, id), "application/json", bytes.NewReader(body))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return fmt.Errorf("failed to set gate device, status: %d", resp.StatusCode)
	}

	return nil
}
