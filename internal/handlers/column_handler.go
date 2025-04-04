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

// CreateColumn godoc
// @Summary Создать колонку
// @Description Создает новую колонку для указанной доски.
// @Tags columns
// @Accept json
// @Produce json
// @Param column body models.Column true "Данные колонки (должны содержать BoardID и Title)"
// @Success 201 {object} models.Column "Колонка успешно создана"
// @Failure 400 {object} map[string]string "Неверные входные данные или отсутствует BoardID/Title"
// @Failure 404 {object} map[string]string "Доска не найдена"
// @Failure 500 {object} map[string]string "Внутренняя ошибка сервера"
// @Router /columns [post]
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

// GetColumn godoc
// @Summary Получить колонку
// @Description Возвращает колонку по её ID.
// @Tags columns
// @Produce json
// @Param column_id path int true "ID колонки"
// @Success 200 {object} models.Column "Колонка успешно найдена"
// @Failure 400 {object} map[string]string "Неверный формат ID колонки"
// @Failure 404 {object} map[string]string "Колонка не найдена"
// @Failure 500 {object} map[string]string "Внутренняя ошибка сервера"
// @Router /columns/{column_id} [get]
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

// GetBoardColumns godoc
// @Summary Получить колонки доски
// @Description Возвращает все колонки, принадлежащие указанной доске.
// @Tags columns
// @Produce json
// @Param board_id path int true "ID доски"
// @Success 200 {array} models.Column "Список колонок доски"
// @Failure 400 {object} map[string]string "Неверный формат ID доски"
// @Failure 404 {object} map[string]string "Доска не найдена"
// @Failure 500 {object} map[string]string "Внутренняя ошибка сервера"
// @Router /boards/{board_id}/columns [get]
func (h *ColumnHandler) GetBoardColumns(c *gin.Context) {
	// Обратите внимание: параметр в URL рекомендуется назвать board_id, хотя в коде используется "column_id"
	boardID, err := strconv.ParseUint(c.Param("column_id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid board ID"})
		return
	}

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

// UpdateColumn godoc
// @Summary Обновить колонку
// @Description Обновляет данные колонки по ID.
// @Tags columns
// @Accept json
// @Produce json
// @Param column_id path int true "ID колонки"
// @Param column body models.Column true "Новые данные колонки"
// @Success 200 {object} models.Column "Колонка успешно обновлена"
// @Failure 400 {object} map[string]string "Неверный формат ID или входные данные"
// @Failure 404 {object} map[string]string "Колонка не найдена"
// @Failure 500 {object} map[string]string "Внутренняя ошибка сервера"
// @Router /columns/{column_id} [put]
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

// DeleteColumn godoc
// @Summary Удалить колонку
// @Description Удаляет колонку по её ID.
// @Tags columns
// @Produce json
// @Param column_id path int true "ID колонки"
// @Success 204 {string} string "Колонка успешно удалена"
// @Failure 400 {object} map[string]string "Неверный формат ID колонки"
// @Failure 404 {object} map[string]string "Колонка не найдена"
// @Failure 500 {object} map[string]string "Внутренняя ошибка сервера"
// @Router /columns/{column_id} [delete]
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
// @Summary Обновить позиции колонок
// @Description Обновляет позиции нескольких колонок в рамках одной доски.
// @Tags columns
// @Accept json
// @Produce json
// @Param input body []models.Column true "Список колонок с новыми позициями"
// @Success 200 {object} map[string]string "Positions updated successfully"
// @Failure 400 {object} map[string]string "Неверные входные данные или отсутствуют колонки"
// @Failure 404 {object} map[string]string "Доска не найдена"
// @Failure 500 {object} map[string]string "Внутренняя ошибка сервера"
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
