package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"kanban-octaview/internal/board"
	"kanban-octaview/pkg/repository"
)

type BoardHandler struct {
	boardRepo repository.BoardRepository
}

func NewBoardHandler(br repository.BoardRepository) *BoardHandler {
	return &BoardHandler{boardRepo: br}
}

func (h *BoardHandler) CreateBoard(c *gin.Context) {
	var b board.Board
	if err := c.ShouldBindJSON(&b); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	// Валидация входных данных (можно использовать github.com/go-playground/validator)
	id, err := h.boardRepo.Create(c.Request.Context(), &b)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "cannot create board"})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"id": id})
}

func (h *BoardHandler) GetBoard(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	b, err := h.boardRepo.GetByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "cannot get board"})
		return
	}
	if b == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "board not found"})
		return
	}
	c.JSON(http.StatusOK, b)
}
