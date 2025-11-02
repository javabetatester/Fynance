package routes

import (
	"Fynance/internal/domain/transaction"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/oklog/ulid/v2"
)

type CreateTransactionRequest struct {
	Type        transaction.Types `json:"type" binding:"required"`
	CategoryId  ulid.ULID         `json:"category_id" binding:"required"`
	Amount      float64           `json:"amount" binding:"required,gt=0"`
	Description string            `json:"description"`
	Date        time.Time         `json:"date" binding:"required"`
}

func (h *Handler) CreateTransaction(c *gin.Context) {
	var req CreateTransactionRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID, err := h.GetUserIDFromContext(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	transaction := transaction.Transaction{
		Type:        req.Type,
		UserId:      userID,
		CategoryId:  req.CategoryId,
		Amount:      req.Amount,
		Description: req.Description,
		Date:        req.Date,
	}

	if err := h.TransactionService.CreateTransaction(&transaction); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message":     "Transaction created successfully",
		"transaction": transaction,
	})
}

func (h *Handler) GetTransactions(c *gin.Context) {

	userID, err := h.GetUserIDFromContext(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	transactions, err := h.TransactionService.GetAllTransactions(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, transactions)
}
