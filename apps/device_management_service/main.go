package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"strconv"

	"device_management_service/db"
	"device_management_service/models"

	"github.com/gin-gonic/gin"
)

// DeviceHandler handles device-related requests
type DeviceHandler struct {
	DB *db.DB
}

// NewDeviceHandler creates a new handler
func NewDeviceHandler(database *db.DB) *DeviceHandler {
	return &DeviceHandler{DB: database}
}

// RegisterRoutes registers the device routes
func (h *DeviceHandler) RegisterRoutes(router *gin.RouterGroup) {
	devices := router.Group("/devices")
	{
		devices.GET("", h.GetDevices)
		devices.GET("/:id", h.GetDeviceByID)
		devices.POST("/:id/command", h.UpdateDeviceValue)
		devices.DELETE("/:id", h.DeleteDevice)
		devices.POST("", h.CreateDevice)
	}
}

func main() {
	dbURL := os.Getenv("DATABASE_URL")
	database, err := db.New(dbURL)
	if err != nil {
		log.Fatalf("failed to connect to DB: %v", err)
	}
	defer database.Close()

	r := gin.Default()
	handler := NewDeviceHandler(database)

	api := r.Group("/api/v1")
	handler.RegisterRoutes(api)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8090"
	}
	r.Run(":" + port)
}

// GetDevices handles GET /devices
func (h *DeviceHandler) GetDevices(c *gin.Context) {
	devices, err := h.DB.GetDevices(context.Background())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, devices)
}

// GetDeviceByID handles GET /devices/:id
func (h *DeviceHandler) GetDeviceByID(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	device, err := h.DB.GetDeviceByID(context.Background(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Device not found"})
		return
	}
	c.JSON(http.StatusOK, device)
}

// UpdateDeviceValue handles POST /devices/:id/command
func (h *DeviceHandler) UpdateDeviceValue(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))

	var req struct {
		Value  *float64 `json:"value"`
		Status *string  `json:"status"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := h.DB.UpdateDeviceValue(context.Background(), id, req.Value, req.Status)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Device updated successfully"})
}

// DeleteDevice handles DELETE /devices/:id
func (h *DeviceHandler) DeleteDevice(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	err := h.DB.DeleteDevice(context.Background(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Device deleted successfully"})
}

func (h *DeviceHandler) CreateDevice(c *gin.Context) {
	var req models.DeviceCreate

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	device, err := h.DB.CreateDevice(context.Background(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, device)
}
