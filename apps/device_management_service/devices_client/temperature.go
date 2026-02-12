package devices_client

import (
	"fmt"
	"math/rand"
	"net/http"
	"sync"
	"time"
)

type TemperatureDevice struct {
	ID     string  `json:"id"`
	Name   string  `json:"name"`
	Value  float64 `json:"value"`
	Status string  `json:"status"`
	mu     sync.Mutex
}

type TemperatureClient struct {
	BaseURL string
}

var device = TemperatureDevice{
	ID:     "1",
	Name:   "Device",
	Value:  21,
	Status: "active",
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

	rand.Seed(time.Now().UnixNano())
	minTemp := 18
	maxTemp := 42

	value := float64(rand.Intn(maxTemp-minTemp+1) + minTemp)

	device = TemperatureDevice{
		ID:     id,
		Name:   device.Name,
		Value:  float64(value),
		Status: device.Status,
	}

	return &device, nil
}

// Set устанавливает новое значение температуры по ID
func (c *TemperatureClient) Set(id string, status string) error {

	device.mu.Lock()
	device.Status = status
	device.mu.Unlock()

	return nil
}
