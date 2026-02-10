package main

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"os"

	"device_management_service/models"

	"github.com/gin-gonic/gin"
)

var devices = []models.Device{
	{ID: "1", Name: "Temperature Sensor 1", Type: models.Temperature, Status: "active", Value: 0},
	// Можно добавить Light и Gate позже
}

func main() {
	telemetryURL := os.Getenv("TELEMETRY_SERVICE_URL")
	if telemetryURL == "" {
		log.Fatal("TELEMETRY_SERVICE_URL not set")
	}

	router := gin.Default()

	router.GET("/devices", func(c *gin.Context) {
		c.JSON(http.StatusOK, devices)
	})

	router.GET("/devices/:id", func(c *gin.Context) {
		id := c.Param("id")
		for _, d := range devices {
			if d.ID == id {
				c.JSON(http.StatusOK, d)
				return
			}
		}
		c.JSON(http.StatusNotFound, gin.H{"error": "Device not found"})
	})

	router.POST("/devices/:id/command", func(c *gin.Context) {
		id := c.Param("id")
		var cmd models.DeviceCommand
		if err := c.ShouldBindJSON(&cmd); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		var device *models.Device
		for i := range devices {
			if devices[i].ID == id {
				device = &devices[i]
				break
			}
		}
		if device == nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Device not found"})
			return
		}

		// Обрабатываем команду только для Temperature для MVP
		if device.Type == models.Temperature && cmd.Value != nil {
			device.Value = *cmd.Value
			device.Status = "active"
		}

		// Отправка события в Telemetry Service
		event := map[string]interface{}{
			"device_id": device.ID,
			"type":      device.Type,
			"value":     device.Value,
			"status":    device.Status,
		}
		payload, _ := json.Marshal(event)
		_, err := http.Post(telemetryURL+"/telemetry", "application/json", bytes.NewBuffer(payload))
		if err != nil {
			log.Printf("Failed to send telemetry: %v", err)
		}

		c.JSON(http.StatusOK, device)
	})

	log.Println("Device Management Service running on :8090")
	if err := router.Run(":8090"); err != nil {
		log.Fatal(err)
	}
}
