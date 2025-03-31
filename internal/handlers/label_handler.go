package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/octaview/kanban-octaview/internal/models"
	"github.com/octaview/kanban-octaview/internal/service"
)

type LabelHandler struct {
	labelService service.LabelServiceInterface
}

func NewLabelHandler(labelService service.LabelServiceInterface) *LabelHandler {
	return &LabelHandler{
		labelService: labelService,
	}
}

// CreateLabel godoc
// @Summary Create a new label
// @Description Create a new label for a board
// @Tags labels
// @Accept json
// @Produce json
// @Param label body models.Label true "Label object"
// @Success 201 {object} models.Label
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/labels [post]
func (h *LabelHandler) CreateLabel(c *gin.Context) {
	var label models.Label
	if err := c.ShouldBindJSON(&label); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.labelService.Create(c.Request.Context(), &label); err != nil {
		handleServiceError(c, err)
		return
	}

	c.JSON(http.StatusCreated, label)
}

// GetLabel godoc
// @Summary Get label by ID
// @Description Get label details by its ID
// @Tags labels
// @Produce json
// @Param label_id path int true "Label ID"
// @Success 200 {object} models.Label
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/labels/{label_id} [get]
func (h *LabelHandler) GetLabel(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("label_id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid label ID"})
		return
	}

	label, err := h.labelService.GetByID(c.Request.Context(), uint(id))
	if err != nil {
		handleServiceError(c, err)
		return
	}

	c.JSON(http.StatusOK, label)
}

// GetBoardLabels godoc
// @Summary Get labels by board ID
// @Description Get all labels for a specific board
// @Tags labels
// @Produce json
// @Param board_id path int true "Board ID"
// @Success 200 {array} models.Label
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/boards/{board_id}/labels [get]
func (h *LabelHandler) GetBoardLabels(c *gin.Context) {
	boardID, err := strconv.ParseUint(c.Param("board_id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid board ID"})
		return
	}

	labels, err := h.labelService.GetByBoardID(c.Request.Context(), uint(boardID))
	if err != nil {
		handleServiceError(c, err)
		return
	}

	c.JSON(http.StatusOK, labels)
}

// UpdateLabel godoc
// @Summary Update label
// @Description Update label details
// @Tags labels
// @Accept json
// @Produce json
// @Param label_id path int true "Label ID"
// @Param label body models.Label true "Label object"
// @Success 200 {object} models.Label
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/labels/{label_id} [put]
func (h *LabelHandler) UpdateLabel(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("label_id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid label ID"})
		return
	}

	var label models.Label
	if err := c.ShouldBindJSON(&label); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	label.ID = uint(id)
	if err := h.labelService.Update(c.Request.Context(), &label); err != nil {
		handleServiceError(c, err)
		return
	}

	c.JSON(http.StatusOK, label)
}

// DeleteLabel godoc
// @Summary Delete label
// @Description Delete a label by its ID
// @Tags labels
// @Produce json
// @Param label_id path int true "Label ID"
// @Success 204 "No Content"
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/labels/{label_id} [delete]
func (h *LabelHandler) DeleteLabel(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("label_id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid label ID"})
		return
	}

	if err := h.labelService.Delete(c.Request.Context(), uint(id)); err != nil {
		handleServiceError(c, err)
		return
	}

	c.Status(http.StatusNoContent)
}