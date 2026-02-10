package main

import (
	"log"
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
)

type GateDevice struct {
	ID     string `json:"id"`
	Name   string `json:"name"`
	Status string `json:"status"` // "open" или "closed"
	mu     sync.Mutex
}

var device = GateDevice{
	ID:     "1",
	Name:   "Main Gate",
	Status: "closed",
}

func main() {
	router := gin.Default()

	// Получить состояние ворот
	router.GET("/gate/:id", func(c *gin.Context) {
		id := c.Param("id")
		if id != device.ID {
			c.JSON(http.StatusNotFound, gin.H{"error": "Device not found"})
			return
		}
		device.mu.Lock()
		defer device.mu.Unlock()
		c.JSON(http.StatusOK, device)
	})

	// Открыть/закрыть ворота
	router.POST("/gate/:id/set", func(c *gin.Context) {
		id := c.Param("id")
		if id != device.ID {
			c.JSON(http.StatusNotFound, gin.H{"error": "Device not found"})
			return
		}

		var payload struct {
			Status string `json:"status"` // "open" или "closed"
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

	log.Println("Gate Device running on :8083")
	if err := router.Run(":8083"); err != nil {
		log.Fatal(err)
	}
}
