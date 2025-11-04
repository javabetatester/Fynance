package routes

import (
	"Fynance/internal/contracts"
	"Fynance/internal/domain/transaction"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *Handler) CreateCategory(c *gin.Context) {
	var body contracts.CategoryCreateRequest
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, contracts.ErrorResponse{Error: err.Error()})
		return
	}

	userID, err := h.GetUserIDFromContext(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, contracts.ErrorResponse{Error: err.Error()})
		return
	}

	category := transaction.Category{
		UserId: userID,
		Name:   body.Name,
		Icon:   body.Icon,
	}

	if err := h.TransactionService.CreateCategory(&category); err != nil {
		c.JSON(http.StatusInternalServerError, contracts.ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusCreated, contracts.MessageResponse{Message: "Categoria criada com sucesso"})
}
