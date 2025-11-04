package routes

import (
	"Fynance/internal/contracts"
	"Fynance/internal/domain/transaction"
	"Fynance/internal/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *Handler) CreateTransaction(c *gin.Context) {
	var body contracts.TransactionCreateRequest

	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, contracts.ErrorResponse{Error: err.Error()})
		return
	}

	userID, err := h.GetUserIDFromContext(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, contracts.ErrorResponse{Error: err.Error()})
		return
	}

	categoryID, err := utils.ParseULID(body.CategoryID)
	if err != nil {
		c.JSON(http.StatusBadRequest, contracts.ErrorResponse{Error: err.Error()})
		return
	}

	transactionEntity := transaction.Transaction{
		Type:        transaction.Types(body.Type),
		UserId:      userID,
		CategoryId:  categoryID,
		Amount:      body.Amount,
		Description: body.Description,
		Date:        utils.SetTimestamps(),
	}

	if err := h.TransactionService.CreateTransaction(&transactionEntity); err != nil {
		c.JSON(http.StatusInternalServerError, contracts.ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusCreated, contracts.TransactionCreateResponse{
		Message:     "Transação criada com sucesso",
		Transaction: transactionEntity,
	})
}

func (h *Handler) GetTransactions(c *gin.Context) {
	userID, err := h.GetUserIDFromContext(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, contracts.ErrorResponse{Error: err.Error()})
		return
	}

	transactions, err := h.TransactionService.GetAllTransactions(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, contracts.ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, transactions)
}
