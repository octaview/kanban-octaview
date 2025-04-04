package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/octaview/kanban-octaview/internal/models"
	"github.com/octaview/kanban-octaview/internal/service"
)

// ErrorResponse represents a standard error response
type ErrorResponse struct {
	Error string `json:"error" example:"Error message details"`
}

// LabelRequest represents the request body for label operations
type LabelRequest struct {
	Name    string `json:"name" example:"Bug" binding:"required"`
	Color   string `json:"color" example:"#FF0000" binding:"required"`
	BoardID uint   `json:"board_id" example:"1" binding:"required"`
}

// LabelResponse represents the response for a single label
type LabelResponse struct {
	ID      uint   `json:"id" example:"1"`
	Name    string `json:"name" example:"Bug"`
	Color   string `json:"color" example:"#FF0000"`
	BoardID uint   `json:"board_id" example:"1"`
}

// LabelsResponse represents the response for multiple labels
type LabelsResponse []LabelResponse

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
// @Param request body LabelRequest true "Label information"
// @Success 201 {object} LabelResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/labels [post]
func (h *LabelHandler) CreateLabel(c *gin.Context) {
	var labelReq LabelRequest
	if err := c.ShouldBindJSON(&labelReq); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}

	label := models.Label{
		Name:    labelReq.Name,
		Color:   labelReq.Color,
		BoardID: labelReq.BoardID,
	}

	if err := h.labelService.Create(c.Request.Context(), &label); err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusCreated, LabelResponse{
		ID:      label.ID,
		Name:    label.Name,
		Color:   label.Color,
		BoardID: label.BoardID,
	})
}

// GetLabel godoc
// @Summary Get label by ID
// @Description Get label details by its ID
// @Tags labels
// @Produce json
// @Param label_id path int true "Label ID"
// @Success 200 {object} LabelResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/labels/{label_id} [get]
func (h *LabelHandler) GetLabel(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("label_id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Invalid label ID"})
		return
	}

	label, err := h.labelService.GetByID(c.Request.Context(), uint(id))
	if err != nil {
		if err.Error() == "record not found" {
			c.JSON(http.StatusNotFound, ErrorResponse{Error: "Label not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, LabelResponse{
		ID:      label.ID,
		Name:    label.Name,
		Color:   label.Color,
		BoardID: label.BoardID,
	})
}

// GetBoardLabels godoc
// @Summary Get labels by board ID
// @Description Get all labels for a specific board
// @Tags labels
// @Produce json
// @Param board_id path int true "Board ID"
// @Success 200 {array} LabelResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/boards/{board_id}/labels [get]
func (h *LabelHandler) GetBoardLabels(c *gin.Context) {
	boardID, err := strconv.ParseUint(c.Param("board_id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Invalid board ID"})
		return
	}

	labels, err := h.labelService.GetByBoardID(c.Request.Context(), uint(boardID))
	if err != nil {
		if err.Error() == "record not found" {
			c.JSON(http.StatusNotFound, ErrorResponse{Error: "Board not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	response := make(LabelsResponse, 0, len(labels))
	for _, label := range labels {
		response = append(response, LabelResponse{
			ID:      label.ID,
			Name:    label.Name,
			Color:   label.Color,
			BoardID: label.BoardID,
		})
	}

	c.JSON(http.StatusOK, response)
}

// UpdateLabel godoc
// @Summary Update label
// @Description Update label details
// @Tags labels
// @Accept json
// @Produce json
// @Param label_id path int true "Label ID"
// @Param request body LabelRequest true "Label information"
// @Success 200 {object} LabelResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/labels/{label_id} [put]
func (h *LabelHandler) UpdateLabel(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("label_id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Invalid label ID"})
		return
	}

	var labelReq LabelRequest
	if err := c.ShouldBindJSON(&labelReq); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}

	label := models.Label{
		ID:      uint(id),
		Name:    labelReq.Name,
		Color:   labelReq.Color,
		BoardID: labelReq.BoardID,
	}

	if err := h.labelService.Update(c.Request.Context(), &label); err != nil {
		if err.Error() == "record not found" {
			c.JSON(http.StatusNotFound, ErrorResponse{Error: "Label not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, LabelResponse{
		ID:      label.ID,
		Name:    label.Name,
		Color:   label.Color,
		BoardID: label.BoardID,
	})
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
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Invalid label ID"})
		return
	}

	if err := h.labelService.Delete(c.Request.Context(), uint(id)); err != nil {
		if err.Error() == "record not found" {
			c.JSON(http.StatusNotFound, ErrorResponse{Error: "Label not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}
