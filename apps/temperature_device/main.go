package main

import (
	"log"
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
)

type TemperatureDevice struct {
	ID     string `json:"id"`
	Name   string `json:"name"`
	Value  string `json:"value"`
	Status string `json:"status"`
	mu     sync.Mutex
}

var device = TemperatureDevice{
	ID:     "1",
	Name:   "Temperature Sensor 1",
	Value:  "22.5",
	Status: "active",
}

func main() {
	router := gin.Default()

	// Получить текущее значение
	router.GET("/temperature/:id", func(c *gin.Context) {
		id := c.Param("id")
		if id != device.ID {
			c.JSON(http.StatusNotFound, gin.H{"error": "Device not found"})
			return
		}
		device.mu.Lock()
		defer device.mu.Unlock()
		c.JSON(http.StatusOK, device)
	})

	// Установить значение (для теста команд от Device Management)
	router.POST("/temperature/:id/set", func(c *gin.Context) {
		id := c.Param("id")
		if id != device.ID {
			c.JSON(http.StatusNotFound, gin.H{"error": "Device not found"})
			return
		}

		var payload struct {
			Value string `json:"value"`
		}
		if err := c.ShouldBindJSON(&payload); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		device.mu.Lock()
		device.Value = payload.Value
		device.Status = "active"
		device.mu.Unlock()

		c.JSON(http.StatusOK, device)
	})

	log.Println("Temperature Device running on :8081")
	if err := router.Run(":8081"); err != nil {
		log.Fatal(err)
	}
}
