package handlers

import (
    "net/http"
    "strconv"
    "time"

    "github.com/gin-gonic/gin"
    "device/internal/models"
)

type DeviceHandler struct {
    DB *DB
}

func (h *DeviceHandler) RegisterRoutes(router *gin.RouterGroup) {
    devices := router.Group("/devices")
    {
        devices.GET("", h.GetDevices)
        devices.GET("/:id", h.GetDeviceByID)
        devices.POST("", h.CreateDevice)
        devices.PUT("/:id", h.UpdateDevice)
        devices.PATCH("/:id/enabled", h.UpdateDeviceEnabled)
    }
}

// Пример: включение/выключение устройства
func (h *DeviceHandler) UpdateDeviceEnabled(c *gin.Context) {
    id, _ := strconv.Atoi(c.Param("id"))
    var req struct {
        Enabled bool `json:"enabled"`
    }
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    err := h.DB.UpdateDeviceEnabled(id, req.Enabled)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusOK, gin.H{"message": "Device status updated"})
}

func (h *DeviceHandler) CreateDevice(c *gin.Context) {
    var d models.Device
    if err := c.ShouldBindJSON(&d); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }
    if err := h.DB.CreateDevice(&d); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    c.JSON(http.StatusCreated, d)
}

func (h *DeviceHandler) UpdateDevice(c *gin.Context) {
    id, _ := strconv.Atoi(c.Param("id"))
    var d models.Device
    if err := c.ShouldBindJSON(&d); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }
    d.ID = id
    if err := h.DB.UpdateDevice(&d); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    c.JSON(http.StatusOK, d)
}

func (h *DeviceHandler) DeleteDevice(c *gin.Context) {
    id, _ := strconv.Atoi(c.Param("id"))
    if err := h.DB.DeleteDevice(id); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    c.JSON(http.StatusOK, gin.H{"message": "deleted"})
}

func (h *DeviceHandler) GetDevice(c *gin.Context) {
    id, _ := strconv.Atoi(c.Param("id"))
    d, err := h.DB.GetDevice(id)
    if err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
        return
    }
    c.JSON(http.StatusOK, d)
}

func (h *DeviceHandler) GetDeviceList(c *gin.Context) {
    list, err := h.DB.GetDeviceList()
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    c.JSON(http.StatusOK, list)
}
