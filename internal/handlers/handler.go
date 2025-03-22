package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/octaview/kanban-octaview/internal/service"
)

type Handler struct {
	Auth   *AuthHandler
	User   *UserHandler
	Board  *BoardHandler
	// Additional handlers will be added here as needed
}

func NewHandler(services *service.Services) *Handler {
	return &Handler{
		Auth:   NewAuthHandler(services.Auth, services.User),
		User:   NewUserHandler(services.User),
		Board:  NewBoardHandler(services.Board),
		// Initialize other handlers
	}
}

func (h *Handler) InitRoutes(router *gin.Engine, authMiddleware gin.HandlerFunc) {
	// Public routes
	auth := router.Group("/auth")
	{
		auth.POST("/register", h.Auth.Register)
		auth.POST("/login", h.Auth.Login)
		auth.POST("/refresh", h.Auth.RefreshToken)
		auth.GET("/me", authMiddleware, h.Auth.GetMe)
	}

	// Protected routes
	api := router.Group("/api", authMiddleware)
	{
		users := api.Group("/users")
		{
			users.GET("/:id", h.User.GetUser)
			users.PUT("/:id", h.User.UpdateUser)
			users.POST("/:id/change-password", h.User.ChangePassword)
			users.DELETE("/:id", h.User.DeleteUser)
		}

		boards := api.Group("/boards")
		{
			boards.POST("", h.Board.CreateBoard)
			boards.GET("", h.Board.GetUserBoards)
			boards.GET("/:id", h.Board.GetBoard)
			boards.PUT("/:id", h.Board.UpdateBoard)
			boards.DELETE("/:id", h.Board.DeleteBoard)
		}

		// Additional routes will be added here
	}
}