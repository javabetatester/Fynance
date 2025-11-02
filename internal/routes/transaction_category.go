package routes

import (
	"Fynance/internal/domain/transaction"
	"Fynance/internal/utils"
	"net/http"
	"github.com/gin-gonic/gin"
)

func (h *Handler) CreateCategory(c *gin.Context) {
	var category transaction.Category
	if err := c.ShouldBindJSON(&category); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID, err := h.GetUserIDFromContext(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	category.UserId = utils.ULIDToString(userID)

	if err := h.TransactionService.CreateCategory(&category); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, "Category created with success")
}
