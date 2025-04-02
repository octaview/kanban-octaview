package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/octaview/kanban-octaview/internal/models"
	"github.com/octaview/kanban-octaview/internal/service"
)

type CommentHandler struct {
	commentService service.CommentServiceInterface
}

func NewCommentHandler(commentService service.CommentServiceInterface) *CommentHandler {
	return &CommentHandler{
		commentService: commentService,
	}
}

// CreateComment godoc
// @Summary Create a new comment
// @Description Create a new comment for a card
// @Tags comments
// @Accept json
// @Produce json
// @Param input body models.Comment true "Comment data"
// @Success 201 {object} models.Comment
// @Failure 400 {object} models.ErrorResponse
// @Failure 404 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /api/comments [post]
func (h *CommentHandler) CreateComment(c *gin.Context) {
	var input models.Comment
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Message: "Invalid input: " + err.Error()})
		return
	}

	// Get user ID from context (set by auth middleware)
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Message: "User ID not found in context"})
		return
	}
	
	// Set the user ID from the authenticated user
	input.UserID = userID.(uint)

	if err := h.commentService.Create(c.Request.Context(), &input); err != nil {
		statusCode := http.StatusInternalServerError
		if err == models.ErrCardNotFound || err == models.ErrUserNotFound {
			statusCode = http.StatusNotFound
		}
		c.JSON(statusCode, models.ErrorResponse{Message: err.Error()})
		return
	}

	c.JSON(http.StatusCreated, input)
}

// GetCommentByID godoc
// @Summary Get comment by ID
// @Description Get a comment by its ID
// @Tags comments
// @Produce json
// @Param comment_id path int true "Comment ID"
// @Success 200 {object} models.Comment
// @Failure 404 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /api/comments/{comment_id} [get]
func (h *CommentHandler) GetCommentByID(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("comment_id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Message: "Invalid comment ID"})
		return
	}

	comment, err := h.commentService.GetByID(c.Request.Context(), uint(id))
	if err != nil {
		statusCode := http.StatusInternalServerError
		if err == models.ErrCommentNotFound {
			statusCode = http.StatusNotFound
		}
		c.JSON(statusCode, models.ErrorResponse{Message: err.Error()})
		return
	}

	c.JSON(http.StatusOK, comment)
}

// GetCommentsByCard godoc
// @Summary Get comments by card ID
// @Description Get all comments for a specific card
// @Tags comments
// @Produce json
// @Param card_id path int true "Card ID"
// @Success 200 {array} models.Comment
// @Failure 404 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /api/cards/{card_id}/comments [get]
func (h *CommentHandler) GetCommentsByCard(c *gin.Context) {
	cardID, err := strconv.ParseUint(c.Param("card_id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Message: "Invalid card ID"})
		return
	}

	comments, err := h.commentService.GetByCardID(c.Request.Context(), uint(cardID))
	if err != nil {
		statusCode := http.StatusInternalServerError
		if err == models.ErrCardNotFound {
			statusCode = http.StatusNotFound
		}
		c.JSON(statusCode, models.ErrorResponse{Message: err.Error()})
		return
	}

	c.JSON(http.StatusOK, comments)
}

// UpdateComment godoc
// @Summary Update a comment
// @Description Update a comment's content
// @Tags comments
// @Accept json
// @Produce json
// @Param comment_id path int true "Comment ID"
// @Param input body models.Comment true "Updated comment data"
// @Success 200 {object} models.Comment
// @Failure 400 {object} models.ErrorResponse
// @Failure 404 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /api/comments/{comment_id} [put]
func (h *CommentHandler) UpdateComment(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("comment_id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Message: "Invalid comment ID"})
		return
	}

	var input models.Comment
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Message: "Invalid input: " + err.Error()})
		return
	}

	// Get the existing comment to check ownership
	existingComment, err := h.commentService.GetByID(c.Request.Context(), uint(id))
	if err != nil {
		statusCode := http.StatusInternalServerError
		if err == models.ErrCommentNotFound {
			statusCode = http.StatusNotFound
		}
		c.JSON(statusCode, models.ErrorResponse{Message: err.Error()})
		return
	}

	// Get user ID from context (set by auth middleware)
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Message: "User ID not found in context"})
		return
	}

	// Check if the user is the owner of the comment
	if existingComment.UserID != userID.(uint) {
		c.JSON(http.StatusForbidden, models.ErrorResponse{Message: "You can only update your own comments"})
		return
	}

	// Set the ID from the URL parameter
	input.ID = uint(id)
	// Keep the original user and card IDs
	input.UserID = existingComment.UserID
	input.CardID = existingComment.CardID

	if err := h.commentService.Update(c.Request.Context(), &input); err != nil {
		statusCode := http.StatusInternalServerError
		if err == models.ErrCommentNotFound {
			statusCode = http.StatusNotFound
		}
		c.JSON(statusCode, models.ErrorResponse{Message: err.Error()})
		return
	}

	c.JSON(http.StatusOK, input)
}

// DeleteComment godoc
// @Summary Delete a comment
// @Description Delete a comment by ID
// @Tags comments
// @Produce json
// @Param comment_id path int true "Comment ID"
// @Success 204 "No Content"
// @Failure 400 {object} models.ErrorResponse
// @Failure 403 {object} models.ErrorResponse
// @Failure 404 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /api/comments/{comment_id} [delete]
func (h *CommentHandler) DeleteComment(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("comment_id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Message: "Invalid comment ID"})
		return
	}

	// Get the existing comment to check ownership
	existingComment, err := h.commentService.GetByID(c.Request.Context(), uint(id))
	if err != nil {
		statusCode := http.StatusInternalServerError
		if err == models.ErrCommentNotFound {
			statusCode = http.StatusNotFound
		}
		c.JSON(statusCode, models.ErrorResponse{Message: err.Error()})
		return
	}

	// Get user ID from context (set by auth middleware)
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Message: "User ID not found in context"})
		return
	}

	// Check if the user is the owner of the comment
	if existingComment.UserID != userID.(uint) {
		c.JSON(http.StatusForbidden, models.ErrorResponse{Message: "You can only delete your own comments"})
		return
	}

	if err := h.commentService.Delete(c.Request.Context(), uint(id)); err != nil {
		statusCode := http.StatusInternalServerError
		if err == models.ErrCommentNotFound {
			statusCode = http.StatusNotFound
		}
		c.JSON(statusCode, models.ErrorResponse{Message: err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}