package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/octaview/kanban-octaview/internal/models"
	"github.com/octaview/kanban-octaview/internal/service"
)

type ColumnHandler struct {
	columnService service.ColumnServiceInterface
}

func NewColumnHandler(columnService service.ColumnServiceInterface) *ColumnHandler {
	return &ColumnHandler{
		columnService: columnService,
	}
}

func (h *ColumnHandler) CreateColumn(c *gin.Context) {
	var input models.Column
	if err := c.BindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input data"})
		return
	}

	if input.BoardID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Board ID is required"})
		return
	}

	if input.Title == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Title is required"})
		return
	}

	if err := h.columnService.Create(c.Request.Context(), &input); err != nil {
		if err == models.ErrBoardNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Board not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, input)
}

func (h *ColumnHandler) GetColumn(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("column_id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid column ID"})
		return
	}

	column, err := h.columnService.GetByID(c.Request.Context(), uint(id))
	if err != nil {
		if err == models.ErrColumnNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Column not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, column)
}

// In internal/handlers/column_handler.go
func (h *ColumnHandler) GetBoardColumns(c *gin.Context) {
    // Change this line from board_id to id
    boardID, err := strconv.ParseUint(c.Param("column_id"), 10, 32)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid board ID"})
        return
    }

    // Rest of the function remains the same
    columns, err := h.columnService.GetByBoardID(c.Request.Context(), uint(boardID))
    if err != nil {
        if err == models.ErrBoardNotFound {
            c.JSON(http.StatusNotFound, gin.H{"error": "Board not found"})
            return
        }
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusOK, columns)
}

func (h *ColumnHandler) UpdateColumn(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("column_id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid column ID"})
		return
	}

	var input models.Column
	if err := c.BindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input data"})
		return
	}

	input.ID = uint(id)

	if err := h.columnService.Update(c.Request.Context(), &input); err != nil {
		if err == models.ErrColumnNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Column not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	updatedColumn, err := h.columnService.GetByID(c.Request.Context(), uint(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, updatedColumn)
}

func (h *ColumnHandler) DeleteColumn(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("column_id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid column ID"})
		return
	}

	if err := h.columnService.Delete(c.Request.Context(), uint(id)); err != nil {
		if err == models.ErrColumnNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Column not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}

// UpdateColumnPositions godoc
// @Summary Update column positions
// @Description Update positions of multiple columns within a board
// @Tags columns
// @Accept json
// @Produce json
// @Param input body []models.Column true "Columns with new positions"
// @Success 200 {string} string "Positions updated successfully"
// @Failure 400 {object} errorResponse
// @Failure 404 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Router /api/columns/positions [put]
func (h *ColumnHandler) UpdateColumnPositions(c *gin.Context) {
	var columns []models.Column
	if err := c.BindJSON(&columns); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input data"})
		return
	}

	if len(columns) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No columns provided"})
		return
	}

	if err := h.columnService.UpdatePositions(c.Request.Context(), columns); err != nil {
		if err == models.ErrBoardNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Board not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Positions updated successfully"})
}