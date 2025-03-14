package main

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/octaview/kanban-octaview/internal/config"
	"github.com/octaview/kanban-octaview/pkg/database"
	"github.com/octaview/kanban-octaview/pkg/logger"
)

func main() {
	logConfig := logger.Config{
		Level:       slog.LevelInfo,
		Development: true,
		Filename:    "logs/app.log",
	}

	log, err := logger.InitLogger(logConfig)
	if err != nil {
		fmt.Printf("Failed to initialize logger: %v\n", err)
		os.Exit(1)
	}
	defer logger.Close()

	cfg, err := config.LoadConfig()
	if err != nil {
		log.Error("Failed to load config", slog.Any("error", err))
	}

	if cfg.App.Env == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	db, err := database.NewPostgresDB(cfg)
	if err != nil {
		log.Error("Failed to initialize database", slog.Any("error", err))
	}
	
	sqlDB, _ := db.DB()
	defer sqlDB.Close()

	router := gin.Default()

	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "ok",
			"time":   time.Now().Format(time.RFC3339),
		})
	})

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", cfg.HTTP.Port),
		Handler: router,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Error("Failed to start server:", slog.Any("error", err))
		}
	}()

	log.Info("Server started", slog.String("port", cfg.HTTP.Port))

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Info("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Error("Server forced to shutdown:", slog.Any("error", err))
	}

	log.Info("Server exited properly")
}
