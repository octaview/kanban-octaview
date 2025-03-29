package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/octaview/kanban-octaview/internal/service"
	"github.com/octaview/kanban-octaview/internal/repository"
)

type Handler struct {
	Auth      *AuthHandler
	User      *UserHandler
	Board     *BoardHandler
	Column    *ColumnHandler
	Card      *CardHandler
	// Additional handlers will be added here as needed
}

func NewHandler(services *service.Services, repos *repository.Repositories) *Handler {
	cardLabelService := service.NewCardLabelService(
		repos.CardLabel,
		repos.Card,
		repos.Label,
		repos.Board,
		repos.Column,
	)

	return &Handler{
		Auth:      NewAuthHandler(services.Auth, services.User),
		User:      NewUserHandler(services.User),
		Board:     NewBoardHandler(services.Board),
		Column:    NewColumnHandler(services.Column),
		Card:      NewCardHandler(services.Card, cardLabelService),
		// Initialize other handlers
	}
}

func (h *Handler) InitRoutes(router *gin.Engine, authMiddleware gin.HandlerFunc) {
    // Public routes remain unchanged
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
            
            // Individual board operations
            boardID := boards.Group("/:id")
            {
                boardID.GET("", h.Board.GetBoard)
                boardID.PUT("", h.Board.UpdateBoard)
                boardID.DELETE("", h.Board.DeleteBoard)
                
                // Change this line - use the same parameter name
                boardID.GET("/columns", h.Column.GetBoardColumns)
            }
        }
        
        columns := api.Group("/columns")
        {
            columns.POST("", h.Column.CreateColumn)
            columns.GET("/:id", h.Column.GetColumn)
            columns.PUT("/:id", h.Column.UpdateColumn)
            columns.DELETE("/:id", h.Column.DeleteColumn)
            columns.PUT("/positions", h.Column.UpdateColumnPositions)
            
            // Column cards routes
            columns.GET("/:column_id/cards", h.Card.GetCardsByColumn)
        }

        // Rest of the routes remain unchanged
        cards := api.Group("/cards")
		{
			cards.POST("", h.Card.CreateCard)
			cards.GET("/:id", h.Card.GetCard)
			cards.PUT("/:id", h.Card.UpdateCard)
			cards.DELETE("/:id", h.Card.DeleteCard)
			cards.PUT("/positions", h.Card.UpdateCardPositions)
			
			// Card actions
			cards.POST("/:id/move", h.Card.MoveCardToColumn)
			cards.POST("/:id/assign", h.Card.AssignCard)
			cards.POST("/:id/unassign", h.Card.UnassignCard)
			cards.PUT("/:id/due-date", h.Card.UpdateDueDate)
			
			// Card labels
			cards.GET("/:id/labels", h.Card.GetCardLabels)
			cards.POST("/:id/labels", h.Card.AddLabelToCard)
			cards.DELETE("/:id/labels", h.Card.RemoveAllLabelsFromCard)
			cards.DELETE("/:id/labels/:label_id", h.Card.RemoveLabelFromCard)
			cards.POST("/:id/labels/batch", h.Card.BatchAddLabelsToCard)
		}
    }
}