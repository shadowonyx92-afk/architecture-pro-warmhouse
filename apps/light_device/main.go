package main

import (
	"log"
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
)

type LightDevice struct {
	ID     string `json:"id"`
	Name   string `json:"name"`
	Status string `json:"status"` // "on" или "off"
	mu     sync.Mutex
}

var device = LightDevice{
	ID:     "1",
	Name:   "Living Room Light",
	Status: "off",
}

func main() {
	router := gin.Default()

	// Получить состояние
	router.GET("/light/:id", func(c *gin.Context) {
		id := c.Param("id")
		if id != device.ID {
			c.JSON(http.StatusNotFound, gin.H{"error": "Device not found"})
			return
		}
		device.mu.Lock()
		defer device.mu.Unlock()
		c.JSON(http.StatusOK, device)
	})

	// Включить/выключить свет
	router.POST("/light/:id/set", func(c *gin.Context) {
		id := c.Param("id")
		if id != device.ID {
			c.JSON(http.StatusNotFound, gin.H{"error": "Device not found"})
			return
		}

		var payload struct {
			Status string `json:"status"` // "on" или "off"
		}
		if err := c.ShouldBindJSON(&payload); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		device.mu.Lock()
		device.Status = payload.Status
		device.mu.Unlock()

		c.JSON(http.StatusOK, device)
	})

	log.Println("Light Device running on :8082")
	if err := router.Run(":8082"); err != nil {
		log.Fatal(err)
	}
}
