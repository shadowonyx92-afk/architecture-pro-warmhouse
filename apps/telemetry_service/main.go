package main

import (
	"log"
	"net/http"
	"os"

	"telemetry_service/db"
	"telemetry_service/models"

	"github.com/gin-gonic/gin"
)

func main() {
	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		log.Fatal("DATABASE_URL not set")
	}

	dbConn, err := db.New(databaseURL)
	if err != nil {
		log.Fatalf("Failed to connect to DB: %v", err)
	}
	defer dbConn.Close()

	router := gin.Default()

	router.POST("/telemetry", func(c *gin.Context) {
		var t models.TelemetryEvent
		if err := c.ShouldBindJSON(&t); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		event, err := dbConn.CreateTelemetry(c.Request.Context(), t)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusCreated, event)
	})

	router.GET("/telemetry/:device_id", func(c *gin.Context) {
		deviceID := c.Param("device_id")
		events, err := dbConn.GetTelemetryByDeviceID(c.Request.Context(), deviceID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, events)
	})

	log.Println("Telemetry Service running on :8100")
	if err := router.Run(":8100"); err != nil {
		log.Fatal(err)
	}
}
