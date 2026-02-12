package handlers

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"smarthome/db"
	"smarthome/models"
	"smarthome/services"

	"github.com/gin-gonic/gin"
)

// SensorHandler обрабатывает HTTP-запросы сенсоров
type SensorHandler struct {
	DB                     *db.DB
	TemperatureService     *services.TemperatureService
	DeviceManagementClient *services.DeviceManagementClient
}

// NewSensorHandler создаёт новый обработчик
func NewSensorHandler(
	db *db.DB,
	tempService *services.TemperatureService,
	dmClient *services.DeviceManagementClient,
) *SensorHandler {
	return &SensorHandler{
		DB:                     db,
		TemperatureService:     tempService,
		DeviceManagementClient: dmClient,
	}
}

// RegisterRoutes registers the sensor routes
func (h *SensorHandler) RegisterRoutes(router *gin.RouterGroup) {
	sensors := router.Group("/sensors")
	{
		sensors.GET("", h.GetSensors)
		sensors.GET("/:id", h.GetSensorByID)
		sensors.POST("", h.CreateSensor)
		sensors.PUT("/:id", h.UpdateSensor)
		sensors.DELETE("/:id", h.DeleteSensor)
		sensors.PATCH("/:id/value", h.UpdateSensorValue)
		sensors.GET("/temperature/:location", h.GetTemperatureByLocation)
	}
}

// GetSensors handles GET /api/v1/sensors
func (h *SensorHandler) GetSensors(c *gin.Context) {
	sensors, err := h.DB.GetSensors(context.Background())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Update temperature sensors with real-time data from the external API
	for i, sensor := range sensors {
		if sensor.Type == models.Temperature {
			tempData, err := h.TemperatureService.GetTemperatureByID(fmt.Sprintf("%d", sensor.ID))
			if err == nil {
				// Update sensor with real-time data
				sensors[i].Value = &tempData.Value
				sensors[i].Status = &tempData.Status
				sensors[i].LastUpdated = tempData.Timestamp
				log.Printf("Updated temperature data for sensor %d from external API", sensor.ID)
			} else {
				log.Printf("Failed to fetch temperature data for sensor %d: %v", sensor.ID, err)
			}
		}
	}

	devices, err := h.DeviceManagementClient.GetDevices()
	if err == nil {
		for _, d := range devices {
			sensors = append(sensors, models.Sensor{
				ID:          d.ID,
				Name:        d.Name,
				Type:        models.SensorType(d.Type),
				Location:    d.Location,
				Value:       d.Value,
				Unit:        d.Unit,
				Status:      d.Status,
				LastUpdated: d.LastUpdated,
				CreatedAt:   d.CreatedAt,
			})
		}
	}

	c.JSON(http.StatusOK, sensors)
}

// GetSensorByID handles GET /api/v1/sensors/:id
func (h *SensorHandler) GetSensorByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid sensor ID"})
		return
	}

	sensor, err := h.DB.GetSensorByID(context.Background(), id)
	if err != nil {
		device, err := h.DeviceManagementClient.GetDeviceByID(id)

		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Sensor not found"})
			return
		}
		c.JSON(http.StatusOK, device)
		return
	}

	// If this is a temperature sensor, fetch real-time data from the temperature API
	if sensor.Type == models.Temperature {
		tempData, err := h.TemperatureService.GetTemperatureByID(fmt.Sprintf("%d", sensor.ID))
		if err == nil {
			// Update sensor with real-time data
			sensor.Value = &tempData.Value
			sensor.Status = &tempData.Status
			sensor.LastUpdated = tempData.Timestamp
			log.Printf("Updated temperature data for sensor %d from external API", sensor.ID)
		} else {
			log.Printf("Failed to fetch temperature data for sensor %d: %v", sensor.ID, err)
		}
	}

	c.JSON(http.StatusOK, sensor)
}

// GetTemperatureByLocation handles GET /api/v1/sensors/temperature/:location
func (h *SensorHandler) GetTemperatureByLocation(c *gin.Context) {
	location := c.Param("location")
	if location == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Location is required"})
		return
	}

	// Fetch temperature data from the external API
	tempData, err := h.TemperatureService.GetTemperature(location)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": fmt.Sprintf("Failed to fetch temperature data: %v", err),
		})
		return
	}

	// Return the temperature data
	c.JSON(http.StatusOK, gin.H{
		"location":    tempData.Location,
		"value":       tempData.Value,
		"unit":        tempData.Unit,
		"status":      tempData.Status,
		"timestamp":   tempData.Timestamp,
		"description": tempData.Description,
	})
}

// ==========================
// Создание и обновление идут в Device Management
// ==========================

func (h *SensorHandler) CreateSensor(c *gin.Context) {
	var req models.SensorCreate
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 1. Прокидываем создание в Device Management Service
	device, err := h.DeviceManagementClient.CreateDevice(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	sensor := models.Sensor{
		ID:          device.ID,
		Name:        device.Name,
		Type:        models.SensorType(device.Type),
		Location:    device.Location,
		Value:       device.Value,
		Unit:        device.Unit,
		Status:      device.Status,
		CreatedAt:   device.CreatedAt,
		LastUpdated: device.LastUpdated,
	}

	c.JSON(http.StatusCreated, sensor)
}

// UpdateSensor handles PUT /api/v1/sensors/:id
func (h *SensorHandler) UpdateSensor(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid sensor ID"})
		return
	}

	var sensorUpdate models.SensorUpdate
	_, err = h.DB.GetSensorByID(context.Background(), id)

	if err == nil {
		if err := c.ShouldBindJSON(&sensorUpdate); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		sensor, err := h.DB.UpdateSensor(context.Background(), id, sensorUpdate)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, sensor)
	} else {
		_, err := h.DeviceManagementClient.GetDeviceByID(id)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		device, err := h.DeviceManagementClient.UpdateDevice(id, sensorUpdate)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, device)
	}
}

// DeleteSensor handles DELETE /api/v1/sensors/:id
func (h *SensorHandler) DeleteSensor(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid sensor ID"})
		return
	}

	err = h.DB.DeleteSensor(context.Background(), id)
	device_err := h.DeviceManagementClient.DeleteDevice(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if device_err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": device_err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Sensor deleted successfully"})
}

// UpdateSensorValue handles PATCH /api/v1/sensors/:id/value
func (h *SensorHandler) UpdateSensorValue(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid sensor ID"})
		return
	}

	var request struct {
		Value  float64 `json:"value" binding:"required"`
		Status string  `json:"status" binding:"required"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	_, err = h.DB.GetSensorByID(context.Background(), id)
	if err != nil {
		_, err = h.DeviceManagementClient.UpdateDeviceValueAndStatus(id, &request.Value, &request.Status)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Sensor value updated successfully"})
		return
	}

	err = h.DB.UpdateSensorValue(context.Background(), id, request.Value, request.Status)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Sensor value updated successfully"})
}

func (h *SensorHandler) fallbackToDeviceService() []models.Sensor {
	if h.DeviceManagementClient == nil {
		return []models.Sensor{}
	}

	devices, err := h.DeviceManagementClient.GetDevices()
	if err != nil {
		return []models.Sensor{}
	}

	var sensors []models.Sensor
	for _, d := range devices {
		sensors = append(sensors, models.Sensor{
			ID:          d.ID,
			Name:        d.Name,
			Type:        models.SensorType(d.Type),
			Location:    d.Location,
			Value:       d.Value,
			Unit:        d.Unit,
			Status:      d.Status,
			LastUpdated: d.LastUpdated,
			CreatedAt:   d.CreatedAt,
		})
	}

	return sensors
}
