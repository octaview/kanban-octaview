package main
import (
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/octaview/kanban-backend/internal/db"
	"github.com/octaview/kanban-backend/internal/handler"
	"github.com/octaview/kanban-backend/pkg/repository"
	"github.com/sirupsen/logrus"
)

func main() {
	if err := godotenv.Load("config/.env.example"); err != nil {
		logrus.Warn("No .env file found")
	}

	logger := logrus.New()
	logger.SetFormatter(&logrus.JSONFormatter{})
	logger.SetReportCaller(true)

	dbConn, err := db.Connect()
	if err != nil {
		logger.Fatal("DB connection error: ", err)
	}
	defer dbConn.Close()

	boardRepo := repository.NewBoardRepository(dbConn)
	boardHandler := handler.NewBoardHandler(boardRepo)

	router := gin.New()
	setupMiddleware(router)
	router.Use(gin.LoggerWithWriter(os.Stdout))
	router.Use(gin.Recovery())
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

func setupMiddleware(r *gin.Engine) {
	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	r.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}
		c.Next()
	})
}
