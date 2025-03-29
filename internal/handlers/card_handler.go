package handlers

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/octaview/kanban-octaview/internal/models"
	"github.com/octaview/kanban-octaview/internal/service"
)

type CardHandler struct {
	cardService      service.CardServiceInterface
	cardLabelService *service.CardLabelService
}

func NewCardHandler(cardService service.CardServiceInterface, cardLabelService *service.CardLabelService) *CardHandler {
	return &CardHandler{
		cardService:      cardService,
		cardLabelService: cardLabelService,
	}
}

// CreateCard godoc
// @Summary Create a new card
// @Description Create a new card in a column
// @Tags cards
// @Accept json
// @Produce json
// @Param input body models.Card true "Card data"
// @Success 201 {object} models.Card
// @Failure 400 {object} models.ValidationError
// @Failure 404 {string} string
// @Failure 500 {string} string
// @Router /api/cards [post]
func (h *CardHandler) CreateCard(c *gin.Context) {
	var card models.Card
	if err := c.ShouldBindJSON(&card); err != nil {
		validErr := models.NewValidationError("request_body", "Invalid request body")
		c.JSON(http.StatusBadRequest, validErr)
		return
	}

	if card.Title == "" {
		validErr := models.NewValidationError("title", "Title is required")
		c.JSON(http.StatusBadRequest, validErr)
		return
	}

	if card.ColumnID == 0 {
		validErr := models.NewValidationError("column_id", "Column ID is required")
		c.JSON(http.StatusBadRequest, validErr)
		return
	}

	if err := h.cardService.Create(c.Request.Context(), &card); err != nil {
		if err == models.ErrColumnNotFound {
			c.JSON(http.StatusNotFound, err.Error())
			return
		}
		
		if validationErr, ok := err.(*models.ValidationError); ok {
			c.JSON(http.StatusBadRequest, validationErr)
			return
		}
		
		c.JSON(http.StatusInternalServerError, "Failed to create card")
		return
	}

	c.JSON(http.StatusCreated, card)
}

// GetCard godoc
// @Summary Get a card by ID
// @Description Get a card by its ID
// @Tags cards
// @Produce json
// @Param id path int true "Card ID"
// @Success 200 {object} models.Card
// @Failure 404 {string} string
// @Failure 500 {string} string
// @Router /api/cards/{id} [get]
func (h *CardHandler) GetCard(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("card_id"), 10, 32)
	if err != nil {
		validErr := models.NewValidationError("id", "Invalid card ID")
		c.JSON(http.StatusBadRequest, validErr)
		return
	}

	card, err := h.cardService.GetByID(c.Request.Context(), uint(id))
	if err != nil {
		if err == models.ErrCardNotFound {
			c.JSON(http.StatusNotFound, err.Error())
			return
		}
		c.JSON(http.StatusInternalServerError, "Failed to get card")
		return
	}

	c.JSON(http.StatusOK, card)
}

// GetCardsByColumn godoc
// @Summary Get cards by column ID
// @Description Get all cards in a column
// @Tags cards
// @Produce json
// @Param column_id path int true "Column ID"
// @Success 200 {array} models.Card
// @Failure 404 {string} string
// @Failure 500 {string} string
// @Router /api/columns/{column_id}/cards [get]
func (h *CardHandler) GetCardsByColumn(c *gin.Context) {
	columnID, err := strconv.ParseUint(c.Param("column_id"), 10, 32)
	if err != nil {
		validErr := models.NewValidationError("column_id", "Invalid column ID")
		c.JSON(http.StatusBadRequest, validErr)
		return
	}

	cards, err := h.cardService.GetByColumnID(c.Request.Context(), uint(columnID))
	if err != nil {
		if err == models.ErrColumnNotFound {
			c.JSON(http.StatusNotFound, err.Error())
			return
		}
		c.JSON(http.StatusInternalServerError, "Failed to get cards")
		return
	}

	c.JSON(http.StatusOK, cards)
}

// UpdateCard godoc
// @Summary Update a card
// @Description Update a card by its ID
// @Tags cards
// @Accept json
// @Produce json
// @Param id path int true "Card ID"
// @Param input body models.Card true "Updated card data"
// @Success 200 {object} models.Card
// @Failure 400 {object} models.ValidationError
// @Failure 404 {string} string
// @Failure 500 {string} string
// @Router /api/cards/{id} [put]
func (h *CardHandler) UpdateCard(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("card_id"), 10, 32)
	if err != nil {
		validErr := models.NewValidationError("id", "Invalid card ID")
		c.JSON(http.StatusBadRequest, validErr)
		return
	}

	var card models.Card
	if err := c.ShouldBindJSON(&card); err != nil {
		validErr := models.NewValidationError("request_body", "Invalid request body")
		c.JSON(http.StatusBadRequest, validErr)
		return
	}

	card.ID = uint(id)

	if err := h.cardService.Update(c.Request.Context(), &card); err != nil {
		if err == models.ErrCardNotFound {
			c.JSON(http.StatusNotFound, err.Error())
			return
		}
		if err == models.ErrColumnNotFound {
			c.JSON(http.StatusNotFound, err.Error())
			return
		}
		
		if validationErr, ok := err.(*models.ValidationError); ok {
			c.JSON(http.StatusBadRequest, validationErr)
			return
		}
		
		c.JSON(http.StatusInternalServerError, "Failed to update card")
		return
	}

	c.JSON(http.StatusOK, card)
}

// DeleteCard godoc
// @Summary Delete a card
// @Description Delete a card by its ID
// @Tags cards
// @Produce json
// @Param id path int true "Card ID"
// @Success 204 "No Content"
// @Failure 404 {string} string
// @Failure 500 {string} string
// @Router /api/cards/{id} [delete]
func (h *CardHandler) DeleteCard(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("card_id"), 10, 32)
	if err != nil {
		validErr := models.NewValidationError("id", "Invalid card ID")
		c.JSON(http.StatusBadRequest, validErr)
		return
	}

	if err := h.cardService.Delete(c.Request.Context(), uint(id)); err != nil {
		if err == models.ErrCardNotFound {
			c.JSON(http.StatusNotFound, err.Error())
			return
		}
		c.JSON(http.StatusInternalServerError, "Failed to delete card")
		return
	}

	c.Status(http.StatusNoContent)
}

// UpdateCardPositions godoc
// @Summary Update card positions
// @Description Update the positions of multiple cards within a column
// @Tags cards
// @Accept json
// @Produce json
// @Param input body []models.Card true "Cards with new positions"
// @Success 204 "No Content"
// @Failure 400 {object} models.ValidationError
// @Failure 404 {string} string
// @Failure 500 {string} string
// @Router /api/cards/positions [put]
func (h *CardHandler) UpdateCardPositions(c *gin.Context) {
	var cards []models.Card
	if err := c.ShouldBindJSON(&cards); err != nil {
		validErr := models.NewValidationError("request_body", "Invalid request body")
		c.JSON(http.StatusBadRequest, validErr)
		return
	}

	if len(cards) == 0 {
		validErr := models.NewValidationError("cards", "No cards provided")
		c.JSON(http.StatusBadRequest, validErr)
		return
	}

	if err := h.cardService.UpdatePositions(c.Request.Context(), cards); err != nil {
		if err == models.ErrColumnNotFound {
			c.JSON(http.StatusNotFound, err.Error())
			return
		}
		c.JSON(http.StatusInternalServerError, "Failed to update card positions")
		return
	}

	c.Status(http.StatusNoContent)
}

// MoveCardToColumn godoc
// @Summary Move a card to another column
// @Description Move a card to another column with a specific position
// @Tags cards
// @Accept json
// @Produce json
// @Param id path int true "Card ID"
// @Param input body struct{ColumnID uint `json:"column_id"`;Position int `json:"position"`} true "New column and position"
// @Success 204 "No Content"
// @Failure 400 {object} models.ValidationError
// @Failure 404 {string} string
// @Failure 500 {string} string
// @Router /api/cards/{id}/move [post]
func (h *CardHandler) MoveCardToColumn(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("card_id"), 10, 32)
	if err != nil {
		validErr := models.NewValidationError("id", "Invalid card ID")
		c.JSON(http.StatusBadRequest, validErr)
		return
	}

	var input struct {
		ColumnID uint `json:"column_id"`
		Position int  `json:"position"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		validErr := models.NewValidationError("request_body", "Invalid request body")
		c.JSON(http.StatusBadRequest, validErr)
		return
	}

	if input.ColumnID == 0 {
		validErr := models.NewValidationError("column_id", "Column ID is required")
		c.JSON(http.StatusBadRequest, validErr)
		return
	}

	if err := h.cardService.MoveCard(c.Request.Context(), uint(id), input.ColumnID, input.Position); err != nil {
		if err == models.ErrCardNotFound {
			c.JSON(http.StatusNotFound, err.Error())
			return
		}
		if err == models.ErrColumnNotFound {
			c.JSON(http.StatusNotFound, err.Error())
			return
		}
		c.JSON(http.StatusInternalServerError, "Failed to move card")
		return
	}

	c.Status(http.StatusNoContent)
}

// AssignCard godoc
// @Summary Assign a card to a user
// @Description Assign a card to a user by user ID
// @Tags cards
// @Accept json
// @Produce json
// @Param id path int true "Card ID"
// @Param input body struct{UserID uint `json:"user_id"`} true "User ID"
// @Success 204 "No Content"
// @Failure 400 {object} models.ValidationError
// @Failure 404 {string} string
// @Failure 500 {string} string
// @Router /api/cards/{id}/assign [post]
func (h *CardHandler) AssignCard(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("card_id"), 10, 32)
	if err != nil {
		validErr := models.NewValidationError("id", "Invalid card ID")
		c.JSON(http.StatusBadRequest, validErr)
		return
	}

	var input struct {
		UserID uint `json:"user_id"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		validErr := models.NewValidationError("request_body", "Invalid request body")
		c.JSON(http.StatusBadRequest, validErr)
		return
	}

	if input.UserID == 0 {
		validErr := models.NewValidationError("user_id", "User ID is required")
		c.JSON(http.StatusBadRequest, validErr)
		return
	}

	if err := h.cardService.AssignCard(c.Request.Context(), uint(id), input.UserID); err != nil {
		if err == models.ErrCardNotFound {
			c.JSON(http.StatusNotFound, err.Error())
			return
		}
		if err == models.ErrUserNotFound {
			authErr, _ := err.(*models.AuthError)
			c.JSON(http.StatusNotFound, authErr)
			return
		}
		c.JSON(http.StatusInternalServerError, "Failed to assign card")
		return
	}

	c.Status(http.StatusNoContent)
}

// UnassignCard godoc
// @Summary Unassign a card
// @Description Remove user assignment from a card
// @Tags cards
// @Produce json
// @Param id path int true "Card ID"
// @Success 204 "No Content"
// @Failure 404 {string} string
// @Failure 500 {string} string
// @Router /api/cards/{id}/unassign [post]
func (h *CardHandler) UnassignCard(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("card_id"), 10, 32)
	if err != nil {
		validErr := models.NewValidationError("id", "Invalid card ID")
		c.JSON(http.StatusBadRequest, validErr)
		return
	}

	if err := h.cardService.UnassignCard(c.Request.Context(), uint(id)); err != nil {
		if err == models.ErrCardNotFound {
			c.JSON(http.StatusNotFound, err.Error())
			return
		}
		c.JSON(http.StatusInternalServerError, "Failed to unassign card")
		return
	}

	c.Status(http.StatusNoContent)
}

// UpdateDueDate godoc
// @Summary Update card due date
// @Description Set or update the due date for a card
// @Tags cards
// @Accept json
// @Produce json
// @Param id path int true "Card ID"
// @Param input body struct{DueDate *time.Time `json:"due_date"`} true "Due date (null to remove)"
// @Success 204 "No Content"
// @Failure 400 {object} models.ValidationError
// @Failure 404 {string} string
// @Failure 500 {string} string
// @Router /api/cards/{id}/due-date [put]
func (h *CardHandler) UpdateDueDate(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("card_id"), 10, 32)
	if err != nil {
		validErr := models.NewValidationError("id", "Invalid card ID")
		c.JSON(http.StatusBadRequest, validErr)
		return
	}

	var input struct {
		DueDate *time.Time `json:"due_date"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		validErr := models.NewValidationError("request_body", "Invalid request body")
		c.JSON(http.StatusBadRequest, validErr)
		return
	}

	if err := h.cardService.UpdateDueDate(c.Request.Context(), uint(id), input.DueDate); err != nil {
		if err == models.ErrCardNotFound {
			c.JSON(http.StatusNotFound, err.Error())
			return
		}
		
		if validationErr, ok := err.(*models.ValidationError); ok {
			c.JSON(http.StatusBadRequest, validationErr)
			return
		}
		
		c.JSON(http.StatusInternalServerError, "Failed to update due date")
		return
	}

	c.Status(http.StatusNoContent)
}

// AddLabelToCard godoc
// @Summary Add a label to a card
// @Description Add a label to a card
// @Tags cards,labels
// @Accept json
// @Produce json
// @Param id path int true "Card ID"
// @Param input body struct{LabelID uint `json:"label_id"`} true "Label ID"
// @Success 204 "No Content"
// @Failure 400 {object} models.ValidationError
// @Failure 404 {string} string
// @Failure 500 {string} string
// @Router /api/cards/{id}/labels [post]
func (h *CardHandler) AddLabelToCard(c *gin.Context) {
	cardID, err := strconv.ParseUint(c.Param("card_id"), 10, 32)
	if err != nil {
		validErr := models.NewValidationError("id", "Invalid card ID")
		c.JSON(http.StatusBadRequest, validErr)
		return
	}

	var input struct {
		LabelID uint `json:"label_id"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		validErr := models.NewValidationError("request_body", "Invalid request body")
		c.JSON(http.StatusBadRequest, validErr)
		return
	}

	if input.LabelID == 0 {
		validErr := models.NewValidationError("label_id", "Label ID is required")
		c.JSON(http.StatusBadRequest, validErr)
		return
	}

	if err := h.cardLabelService.AddLabelToCard(c.Request.Context(), uint(cardID), input.LabelID); err != nil {
		if err == models.ErrCardNotFound {
			c.JSON(http.StatusNotFound, err.Error())
			return
		}
		if err == models.ErrLabelNotFound {
			c.JSON(http.StatusNotFound, err.Error())
			return
		}
		c.JSON(http.StatusInternalServerError, "Failed to add label to card")
		return
	}

	c.Status(http.StatusNoContent)
}

// RemoveLabelFromCard godoc
// @Summary Remove a label from a card
// @Description Remove a label from a card
// @Tags cards,labels
// @Accept json
// @Produce json
// @Param id path int true "Card ID"
// @Param label_id path int true "Label ID"
// @Success 204 "No Content"
// @Failure 400 {object} models.ValidationError
// @Failure 404 {string} string
// @Failure 500 {string} string
// @Router /api/cards/{id}/labels/{label_id} [delete]
func (h *CardHandler) RemoveLabelFromCard(c *gin.Context) {
	cardID, err := strconv.ParseUint(c.Param("card_id"), 10, 32)
	if err != nil {
		validErr := models.NewValidationError("id", "Invalid card ID")
		c.JSON(http.StatusBadRequest, validErr)
		return
	}

	labelID, err := strconv.ParseUint(c.Param("label_id"), 10, 32)
	if err != nil {
		validErr := models.NewValidationError("label_id", "Invalid label ID")
		c.JSON(http.StatusBadRequest, validErr)
		return
	}

	if err := h.cardLabelService.RemoveLabelFromCard(c.Request.Context(), uint(cardID), uint(labelID)); err != nil {
		if err == models.ErrCardNotFound {
			c.JSON(http.StatusNotFound, err.Error())
			return
		}
		if err == models.ErrLabelNotFound {
			c.JSON(http.StatusNotFound, err.Error())
			return
		}
		c.JSON(http.StatusInternalServerError, "Failed to remove label from card")
		return
	}

	c.Status(http.StatusNoContent)
}

// GetCardLabels godoc
// @Summary Get all labels for a card
// @Description Get all labels associated with a card
// @Tags cards,labels
// @Produce json
// @Param id path int true "Card ID"
// @Success 200 {array} models.Label
// @Failure 404 {string} string
// @Failure 500 {string} string
// @Router /api/cards/{id}/labels [get]
func (h *CardHandler) GetCardLabels(c *gin.Context) {
	cardID, err := strconv.ParseUint(c.Param("card_id"), 10, 32)
	if err != nil {
		validErr := models.NewValidationError("id", "Invalid card ID")
		c.JSON(http.StatusBadRequest, validErr)
		return
	}

	labels, err := h.cardLabelService.GetLabelsByCardID(c.Request.Context(), uint(cardID))
	if err != nil {
		if err == models.ErrCardNotFound {
			c.JSON(http.StatusNotFound, err.Error())
			return
		}
		c.JSON(http.StatusInternalServerError, "Failed to get card labels")
		return
	}

	c.JSON(http.StatusOK, labels)
}

// BatchAddLabelsToCard godoc
// @Summary Add multiple labels to a card
// @Description Add multiple labels to a card in a single request
// @Tags cards,labels
// @Accept json
// @Produce json
// @Param id path int true "Card ID"
// @Param input body struct{LabelIDs []uint `json:"label_ids"`} true "Label IDs"
// @Success 204 "No Content"
// @Failure 400 {object} models.ValidationError
// @Failure 404 {string} string
// @Failure 500 {string} string
// @Router /api/cards/{id}/labels/batch [post]
func (h *CardHandler) BatchAddLabelsToCard(c *gin.Context) {
	cardID, err := strconv.ParseUint(c.Param("card_id"), 10, 32)
	if err != nil {
		validErr := models.NewValidationError("id", "Invalid card ID")
		c.JSON(http.StatusBadRequest, validErr)
		return
	}

	var input struct {
		LabelIDs []uint `json:"label_ids"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		validErr := models.NewValidationError("request_body", "Invalid request body")
		c.JSON(http.StatusBadRequest, validErr)
		return
	}

	if len(input.LabelIDs) == 0 {
		validErr := models.NewValidationError("label_ids", "No label IDs provided")
		c.JSON(http.StatusBadRequest, validErr)
		return
	}

	if err := h.cardLabelService.BatchAddLabelsToCard(c.Request.Context(), uint(cardID), input.LabelIDs); err != nil {
		if err == models.ErrCardNotFound {
			c.JSON(http.StatusNotFound, err.Error())
			return
		}
		c.JSON(http.StatusInternalServerError, "Failed to add labels to card")
		return
	}

	c.Status(http.StatusNoContent)
}

// RemoveAllLabelsFromCard godoc
// @Summary Remove all labels from a card
// @Description Remove all labels associated with a card
// @Tags cards,labels
// @Produce json
// @Param id path int true "Card ID"
// @Success 204 "No Content"
// @Failure 404 {string} string
// @Failure 500 {string} string
// @Router /api/cards/{id}/labels [delete]
func (h *CardHandler) RemoveAllLabelsFromCard(c *gin.Context) {
	cardID, err := strconv.ParseUint(c.Param("card_id"), 10, 32)
	if err != nil {
		validErr := models.NewValidationError("id", "Invalid card ID")
		c.JSON(http.StatusBadRequest, validErr)
		return
	}

	if err := h.cardLabelService.RemoveAllLabelsFromCard(c.Request.Context(), uint(cardID)); err != nil {
		if err == models.ErrCardNotFound {
			c.JSON(http.StatusNotFound, err.Error())
			return
		}
		c.JSON(http.StatusInternalServerError, "Failed to remove labels from card")
		return
	}

	c.Status(http.StatusNoContent)
}