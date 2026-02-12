package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"smarthome/models"
)

// DeviceManagementClient для взаимодействия с device_management_service
type DeviceManagementClient struct {
	BaseURL string
	Client  *http.Client
}

// NewDeviceManagementClient создаёт новый клиент
func NewDeviceManagementClient(baseURL string) *DeviceManagementClient {
	return &DeviceManagementClient{
		BaseURL: baseURL,
		Client:  &http.Client{},
	}
}

// ==========================
// Чтение устройств
// ==========================

func (c *DeviceManagementClient) GetDevices() ([]*models.Sensor, error) {
	url := fmt.Sprintf("%s/api/v1/devices", c.BaseURL)
	resp, err := c.Client.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var devices []*models.Sensor
	if err := json.NewDecoder(resp.Body).Decode(&devices); err != nil {
		return nil, err
	}
	return devices, nil
}

func (c *DeviceManagementClient) GetDeviceByID(id int) (*models.Sensor, error) {
	url := fmt.Sprintf("%s/api/v1/devices/%d", c.BaseURL, id)
	resp, err := c.Client.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("device not found")
	}

	var device models.Sensor
	if err := json.NewDecoder(resp.Body).Decode(&device); err != nil {
		return nil, err
	}
	return &device, nil
}

// ==========================
// Создание/обновление устройств
// ==========================

func (c *DeviceManagementClient) CreateDevice(req models.SensorCreate) (*models.Sensor, error) {
	url := fmt.Sprintf("%s/api/v1/devices", c.BaseURL)
	data, _ := json.Marshal(req)

	resp, err := c.Client.Post(url, "application/json", bytes.NewBuffer(data))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var device models.Device
	fmt.Println("resp.Body", resp.Body)
	if err := json.NewDecoder(resp.Body).Decode(&device); err != nil {
		return nil, err
	}

	fmt.Println("&device", &device)

	sensor := models.Sensor{
		ID:          device.ID,
		Name:        device.Name,
		Type:        models.SensorType("temperature"),
		Location:    device.Location,
		Value:       &device.Value,
		Unit:        device.Unit,
		Status:      &device.Status,
		LastUpdated: device.LastUpdated,
		CreatedAt:   device.CreatedAt,
	}
	return &sensor, nil
}

func (c *DeviceManagementClient) UpdateDevice(id int, req models.SensorUpdate) (*models.Sensor, error) {
	url := fmt.Sprintf("%s/api/v1/devices/%d", c.BaseURL, id)
	data, _ := json.Marshal(req)

	request, err := http.NewRequest(http.MethodPut, url, bytes.NewBuffer(data))
	if err != nil {
		return nil, err
	}
	request.Header.Set("Content-Type", "application/json")

	resp, err := c.Client.Do(request)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var device models.Sensor
	fmt.Println("resp.Body", resp.Body)

	if err := json.NewDecoder(resp.Body).Decode(&device); err != nil {
		return nil, err
	}
	fmt.Println("&device", &device)
	return &device, nil
}

func (c *DeviceManagementClient) UpdateDeviceValueAndStatus(id int, value *float64, status *string) (*models.Sensor, error) {
	url := fmt.Sprintf("%s/api/v1/devices/%d/command", c.BaseURL, id)

	payload := map[string]interface{}{}
	if value != nil {
		payload["value"] = *value
	}
	if status != nil {
		payload["status"] = *status
	}

	data, _ := json.Marshal(payload)

	request, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(data))
	if err != nil {
		return nil, err
	}
	request.Header.Set("Content-Type", "application/json")

	resp, err := c.Client.Do(request)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var device models.Device
	body, _ := io.ReadAll(resp.Body)
	fmt.Println("resp.Body:----", string(body))
	if err := json.Unmarshal(body, &device); err != nil {
		return nil, err
	}

	return &models.Sensor{
		ID:          device.ID,
		Name:        device.Name,
		Type:        models.SensorType(device.Type),
		Location:    device.Location,
		Value:       &device.Value,
		Unit:        device.Unit,
		Status:      &device.Status,
		LastUpdated: device.LastUpdated,
		CreatedAt:   device.CreatedAt,
	}, nil
}

// ==========================
// Удаление устройства
// ==========================

func (c *DeviceManagementClient) DeleteDevice(id int) error {
	url := fmt.Sprintf("%s/api/v1/devices/%d", c.BaseURL, id)
	request, err := http.NewRequest(http.MethodDelete, url, nil)
	if err != nil {
		return err
	}

	resp, err := c.Client.Do(request)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to delete device, status: %d", resp.StatusCode)
	}

	return nil
}
