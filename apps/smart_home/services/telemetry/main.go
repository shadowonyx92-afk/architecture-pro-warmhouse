package main

import (
    "context"
    "log"
    "net/http"
    "os"
    "os/signal"
    "time"

    "telemetry/internal/db"
    "telemetry/internal/handlers"

    "github.com/gin-gonic/gin"
)

func main() {
    // Подключение к БД
    dbURL := getEnv("DATABASE_URL", "postgres://postgres:postgres@postgres:5432/smarthome")
    database, err := db.New(dbURL)
    if err != nil {
        log.Fatalf("Failed to connect to DB: %v", err)
    }
    defer database.Close()
    log.Println("Connected to database")

    // Создаем handler
    telemetryHandler := handlers.NewTelemetryHandler(database)

    // Настройка Gin
    router := gin.Default()
    api := router.Group("/api/v1")
    telemetryHandler.RegisterRoutes(api)

    // Запуск сервера
    port := getEnv("PORT", "8083")
    srv := &http.Server{
        Addr:    ":" + port,
        Handler: router,
    }

    go func() {
        log.Printf("TelemetryService listening on port %s", port)
        if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
            log.Fatalf("Server failed: %v", err)
        }
    }()

    // Graceful shutdown
    quit := make(chan os.Signal, 1)
    signal.Notify(quit, os.Interrupt)
    <-quit
    log.Println("Shutting down server...")
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()
    if err := srv.Shutdown(ctx); err != nil {
        log.Fatalf("Server forced to shutdown: %v", err)
    }
    log.Println("Server exited")
}

// getEnv возвращает значение переменной окружения или default
func getEnv(key, defaultValue string) string {
    if value := os.Getenv(key); value != "" {
        return value
    }
    return defaultValue
}
