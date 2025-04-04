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
	"github.com/octaview/kanban-octaview/internal/handlers"
	"github.com/octaview/kanban-octaview/internal/middleware"
	"github.com/octaview/kanban-octaview/internal/repository"
	"github.com/octaview/kanban-octaview/internal/service"
	"github.com/octaview/kanban-octaview/pkg/database"
	"github.com/octaview/kanban-octaview/pkg/logger"

	// Swagger
	_ "github.com/octaview/kanban-octaview/docs" // импорт сгенерированной документации
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title Kanban Octaview API
// @version 1.0
// @description API для приложения Kanban.
// @host localhost:8080
// @BasePath /
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
		os.Exit(1)
	}

	if err := cfg.Validate(); err != nil {
		log.Error("Invalid configuration", slog.Any("error", err))
		os.Exit(1)
	}

	if cfg.App.Env == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	db, err := database.NewPostgresDB(cfg)
	if err != nil {
		log.Error("Failed to initialize database", slog.Any("error", err))
		os.Exit(1)
	}

	sqlDB, _ := db.DB()
	defer sqlDB.Close()

	repos := repository.NewRepositories(db)

	services := service.NewServices(repos, cfg)

	authMiddleware := middleware.NewAuthMiddleware(services.Auth)

	handler := handlers.NewHandler(services, repos)

	router := gin.Default()

	// Endpoint для проверки работоспособности
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "ok",
			"time":   time.Now().Format(time.RFC3339),
		})
	})

	// Добавляем маршрут для Swagger UI
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	handler.InitRoutes(router, authMiddleware.AuthRequired())

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
