package handlers

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/octaview/kanban-octaview/internal/models"
	"github.com/octaview/kanban-octaview/internal/service"
)

type BoardHandler struct {
	boardService service.BoardServiceInterface
}

func NewBoardHandler(boardService service.BoardServiceInterface) *BoardHandler {
	return &BoardHandler{
		boardService: boardService,
	}
}

func (h *BoardHandler) CreateBoard(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	var input models.Board
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	input.OwnerID = userID.(uint)
	
	err := h.boardService.Create(c.Request.Context(), &input)
	if err != nil {
		if errors.Is(err, models.ErrUserNotFound) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "user not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create board"})
		return
	}

	c.JSON(http.StatusCreated, input)
}

func (h *BoardHandler) GetBoard(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("board_id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid board ID"})
		return
	}

	board, err := h.boardService.GetByID(c.Request.Context(), uint(id))
	if err != nil {
		if errors.Is(err, models.ErrBoardNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "board not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get board"})
		return
	}

	c.JSON(http.StatusOK, board)
}

func (h *BoardHandler) GetUserBoards(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	boards, err := h.boardService.GetByOwnerID(c.Request.Context(), userID.(uint))
	if err != nil {
		if errors.Is(err, models.ErrUserNotFound) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "user not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get boards"})
		return
	}

	c.JSON(http.StatusOK, boards)
}

func (h *BoardHandler) UpdateBoard(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	id, err := strconv.ParseUint(c.Param("board_id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid board ID"})
		return
	}

	var input models.Board
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	input.ID = uint(id)

	existingBoard, err := h.boardService.GetByID(c.Request.Context(), uint(id))
	if err != nil {
		if errors.Is(err, models.ErrBoardNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "board not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get board"})
		return
	}

	if existingBoard.OwnerID != userID.(uint) {
		c.JSON(http.StatusForbidden, gin.H{"error": "you don't have permission to update this board"})
		return
	}

	input.OwnerID = existingBoard.OwnerID

	if err := h.boardService.Update(c.Request.Context(), &input); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update board"})
		return
	}

	c.JSON(http.StatusOK, input)
}

func (h *BoardHandler) DeleteBoard(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	id, err := strconv.ParseUint(c.Param("board_id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid board ID"})
		return
	}

	existingBoard, err := h.boardService.GetByID(c.Request.Context(), uint(id))
	if err != nil {
		if errors.Is(err, models.ErrBoardNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "board not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get board"})
		return
	}

	if existingBoard.OwnerID != userID.(uint) {
		c.JSON(http.StatusForbidden, gin.H{"error": "you don't have permission to delete this board"})
		return
	}

	if err := h.boardService.Delete(c.Request.Context(), uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete board"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "board deleted successfully"})
}