package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"device_management_service/db"
	"device_management_service/devices_client"
	"device_management_service/models"

	"github.com/gin-gonic/gin"
)

// DeviceHandler обрабатывает HTTP-запросы по устройствам
type DeviceHandler struct {
	DB                *db.DB
	TemperatureClient *devices_client.TemperatureClient
	LightClient       *devices_client.LightClient
	GateClient        *devices_client.GateClient
}

// Создание нового DeviceHandler с клиентами
func NewDeviceHandler(
	database *db.DB,
	tempClient *devices_client.TemperatureClient,
	lightClient *devices_client.LightClient,
	gateClient *devices_client.GateClient,
) *DeviceHandler {
	return &DeviceHandler{
		DB:                database,
		TemperatureClient: tempClient,
		LightClient:       lightClient,
		GateClient:        gateClient,
	}
}

// Регистрация маршрутов
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

// ==========================
// Обработчики
// ==========================

// GetDevices — список всех устройств
func (h *DeviceHandler) GetDevices(c *gin.Context) {
	devices, err := h.DB.GetDevices(context.Background())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Обновление актуальных данных для внешних устройств
	for i, d := range devices {
		switch d.Type {
		case models.Temperature:
			if temp, err := h.TemperatureClient.Get(strconv.Itoa(d.ID)); err == nil {
				devices[i].Value = &temp.Value
				devices[i].Status = &temp.Status
			}
		case models.Light:
			if light, err := h.LightClient.Get(strconv.Itoa(d.ID)); err == nil {
				devices[i].Status = &light.Status
			}
		case models.Gate:
			if gate, err := h.GateClient.Get(strconv.Itoa(d.ID)); err == nil {
				devices[i].Status = &gate.Status
			}
		}
	}

	c.JSON(http.StatusOK, devices)
}

// GetDeviceByID — конкретное устройство
func (h *DeviceHandler) GetDeviceByID(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))

	device, err := h.DB.GetDeviceByID(context.Background(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Device not found"})
		return
	}

	// Подтянуть актуальные данные для внешних устройств
	switch device.Type {
	case models.Temperature:
		if temp, err := h.TemperatureClient.Get(strconv.Itoa(device.ID)); err == nil {
			device.Value = &temp.Value
			device.Status = &temp.Status
		}
	case models.Light:
		if light, err := h.LightClient.Get(strconv.Itoa(device.ID)); err == nil {
			device.Status = &light.Status
		}
	case models.Gate:
		if gate, err := h.GateClient.Get(strconv.Itoa(device.ID)); err == nil {
			device.Status = &gate.Status
		}
	}

	c.JSON(http.StatusOK, device)
}

// UpdateDeviceValue — POST /devices/:id/command
func (h *DeviceHandler) UpdateDeviceValue(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))

	var req struct {
		Value  *string `json:"value"`
		Status *string `json:"status"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	device, err := h.DB.GetDeviceByID(context.Background(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Device not found"})
		return
	}

	isTest := c.Query("isTest") == "true"

	// Отправляем команду на внешнее устройство
	if isTest {
		// Обновляем БД для консистентности
		if err := h.DB.UpdateDeviceValue(context.Background(), id, req.Value, req.Status); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"Device updated successfully1": req})
		return
	}

	switch device.Type {
	case models.Temperature:
		if req.Value != nil {
			temp, err := strconv.ParseFloat(*req.Value, 64) // 64 означает float64
			if err != nil {
				fmt.Println("Ошибка преобразования температуры:", err)
				return
			}
			if err := h.TemperatureClient.Set(strconv.Itoa(device.ID), temp); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
		}
	case models.Light:
		if req.Status != nil {
			if err := h.LightClient.Set(strconv.Itoa(device.ID), *req.Status); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
		}
	case models.Gate:
		if req.Status != nil {
			if err := h.GateClient.Set(strconv.Itoa(device.ID), *req.Status); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
		}
	}

	// Обновляем БД для консистентности
	if err := h.DB.UpdateDeviceValue(context.Background(), id, req.Value, req.Status); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"Device updated successfully2": req.Value})
}

// DeleteDevice — удалить устройство
func (h *DeviceHandler) DeleteDevice(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	err := h.DB.DeleteDevice(context.Background(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Device deleted successfully"})
}

// CreateDevice — добавить новое устройство
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

	switch device.Type {
	case models.Temperature:
		if temp, err := h.TemperatureClient.Get(strconv.Itoa(device.ID)); err == nil {
			device.Value = &temp.Value
			device.Status = &temp.Status
		}
	case models.Light:
		if light, err := h.LightClient.Get(strconv.Itoa(device.ID)); err == nil {
			device.Status = &light.Status
		}
	case models.Gate:
		if gate, err := h.GateClient.Get(strconv.Itoa(device.ID)); err == nil {
			device.Status = &gate.Status
		}
	}

	c.JSON(http.StatusCreated, device)
}

// ==========================
// main.go
// ==========================

func main() {
	dbURL := os.Getenv("DATABASE_URL")
	database, err := db.New(dbURL)
	if err != nil {
		log.Fatalf("failed to connect to DB: %v", err)
	}
	defer database.Close()

	// Инициализация клиентов устройств
	tempClient := &devices_client.TemperatureClient{BaseURL: os.Getenv("TEMPERATURE_API_URL")}
	lightClient := &devices_client.LightClient{BaseURL: os.Getenv("LIGHT_API_URL")}
	gateClient := &devices_client.GateClient{BaseURL: os.Getenv("GATE_API_URL")}

	handler := NewDeviceHandler(database, tempClient, lightClient, gateClient)

	r := gin.Default()
	api := r.Group("/api/v1")
	handler.RegisterRoutes(api)

	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "ok",
		})
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "8090"
	}
	r.Run(":" + port)
}
