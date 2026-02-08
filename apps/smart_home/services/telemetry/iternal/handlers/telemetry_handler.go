package handlers

import (
    "net/http"
    "strconv"
    "time"

    "github.com/gin-gonic/gin"
    "telemetry/internal/models"
)

type TelemetryHandler struct {
    DB *DB
}

func (h *TelemetryHandler) RegisterRoutes(router *gin.RouterGroup) {
    telemetry := router.Group("/telemetry")
    {
        telemetry.POST("", h.CreateTelemetry)
        telemetry.GET("", h.GetTelemetry)
        telemetry.GET("/latest", h.GetLatestTelemetry)
    }
}

func (h *TelemetryHandler) CreateTelemetry(c *gin.Context) {
    var t models.TelemetryData
    if err := c.ShouldBindJSON(&t); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    if t.Timestamp.IsZero() {
        t.Timestamp = time.Now()
    }

    err := h.DB.InsertTelemetry(&t)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusCreated, t)
}

// Получить историю за период
func (h *TelemetryHandler) GetTelemetry(c *gin.Context) {
    deviceID, _ := strconv.Atoi(c.Query("device_id"))
    from := c.Query("from")
    to := c.Query("to")

    telemetry, err := h.DB.GetTelemetry(deviceID, from, to)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusOK, telemetry)
}

// Последнее состояние устройства
func (h *TelemetryHandler) GetLatestTelemetry(c *gin.Context) {
    deviceID, _ := strconv.Atoi(c.Query("device_id"))

    t, err := h.DB.GetLatestTelemetry(deviceID)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusOK, t)
}