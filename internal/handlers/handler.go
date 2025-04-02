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
	Label     *LabelHandler
	Comment   *CommentHandler
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
		Label:     NewLabelHandler(services.Label),
		Comment:   NewCommentHandler(services.Comment), // Initialize CommentHandler
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
            boardID := boards.Group("/:board_id")  // Changed from ":id" to ":board_id"
            {
                boardID.GET("", h.Board.GetBoard)
                boardID.PUT("", h.Board.UpdateBoard)
                boardID.DELETE("", h.Board.DeleteBoard)
                
                // Now using ":board_id" consistently
                boardID.GET("/columns", h.Column.GetBoardColumns)
            }
        }
        
        columns := api.Group("/columns")
        {
            columns.POST("", h.Column.CreateColumn)
            columns.GET("/:column_id", h.Column.GetColumn)  // Changed from ":id" to ":column_id"
            columns.PUT("/:column_id", h.Column.UpdateColumn)  // Changed from ":id" to ":column_id"
            columns.DELETE("/:column_id", h.Column.DeleteColumn)  // Changed from ":id" to ":column_id"
            columns.PUT("/positions", h.Column.UpdateColumnPositions)
            
            // Column cards routes - already using specific parameter name
            columns.GET("/:column_id/cards", h.Card.GetCardsByColumn)
        }

        // Rest of the routes remain unchanged
        cards := api.Group("/cards")
        {
            cards.POST("", h.Card.CreateCard)
            cards.GET("/:card_id", h.Card.GetCard)  // Changed from ":id" to ":card_id"
            cards.PUT("/:card_id", h.Card.UpdateCard)  // Changed from ":id" to ":card_id"
            cards.DELETE("/:card_id", h.Card.DeleteCard)  // Changed from ":id" to ":card_id"
            cards.PUT("/positions", h.Card.UpdateCardPositions)
            
            // Card actions - using consistent parameter names
            cards.POST("/:card_id/move", h.Card.MoveCardToColumn)  // Changed from ":id" to ":card_id"
            cards.POST("/:card_id/assign", h.Card.AssignCard)  // Changed from ":id" to ":card_id"
            cards.POST("/:card_id/unassign", h.Card.UnassignCard)  // Changed from ":id" to ":card_id"
            cards.PUT("/:card_id/due-date", h.Card.UpdateDueDate)  // Changed from ":id" to ":card_id"
            
            // Card labels - using consistent parameter names
            cards.GET("/:card_id/labels", h.Card.GetCardLabels)  // Changed from ":id" to ":card_id"
            cards.POST("/:card_id/labels", h.Card.AddLabelToCard)  // Changed from ":id" to ":card_id"
            cards.DELETE("/:card_id/labels", h.Card.RemoveAllLabelsFromCard)  // Changed from ":id" to ":card_id"
            cards.DELETE("/:card_id/labels/:label_id", h.Card.RemoveLabelFromCard)  // Changed from ":id" to ":card_id"
            cards.POST("/:card_id/labels/batch", h.Card.BatchAddLabelsToCard)  // Changed from ":id" to ":card_id"
            
            // Card comments - add routes for comments
            cards.GET("/:card_id/comments", h.Comment.GetCommentsByCard)
        }
        
        labels := api.Group("/labels")
        {
            labels.POST("", h.Label.CreateLabel)
            labels.GET("/:label_id", h.Label.GetLabel)
            labels.PUT("/:label_id", h.Label.UpdateLabel)
            labels.DELETE("/:label_id", h.Label.DeleteLabel)
        }
        
        // Add comment routes
        comments := api.Group("/comments")
        {
            comments.POST("", h.Comment.CreateComment)
            comments.GET("/:comment_id", h.Comment.GetCommentByID)
            comments.PUT("/:comment_id", h.Comment.UpdateComment)
            comments.DELETE("/:comment_id", h.Comment.DeleteComment)
        }
    }
}