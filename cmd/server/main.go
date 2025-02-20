package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"kanban-octaview/internal/db"
	"kanban-octaview/internal/handler"
	"kanban-octaview/pkg/repository"
	"kanban-octaview/pkg/service"
	"github.com/sirupsen/logrus"
)

func main() {
	// Загрузка конфигурации из .env
	if err := godotenv.Load("config/.env.example"); err != nil {
		logrus.Warn("No .env file found")
	}

	// Инициализация логгера
	logger := logrus.New()
	logger.SetFormatter(&logrus.JSONFormatter{})
	logger.SetReportCaller(true)

	// Подключение к БД
	dbConn, err := db.Connect()
	if err != nil {
		logger.Fatal("DB connection error: ", err)
	}
	defer dbConn.Close()

	// Инициализация репозиториев и сервисов
	boardRepo := repository.NewBoardRepository(dbConn)
	boardHandler := handler.NewBoardHandler(boardRepo)
	// Здесь можно инициализировать сервисы (например, TaskService) и аутентификацию

	// Настройка Gin
	router := gin.New()
	router.Use(gin.LoggerWithWriter(os.Stdout))
	router.Use(gin.Recovery())
	// Пример middleware для CORS:
	router.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type,Authorization")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	})

	// Настройка роутинга
	api := router.Group("/api")
	{
		boards := api.Group("/boards")
		{
			boards.POST("", boardHandler.CreateBoard)
			boards.GET("/:id", boardHandler.GetBoard)
		}
		// Здесь добавьте остальные группы (tasks, columns, auth и пр.)
	}

	// Запуск сервера
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	router.Run(":" + port)
}
