package routes

import (
	"errors"
	"net/http"

	"Fynance/internal/contracts"
	"Fynance/internal/domain/transaction"
	"Fynance/internal/utils"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
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

	c.JSON(http.StatusOK, contracts.TransactionListResponse{Transactions: transactions, Total: len(transactions)})
}

func (h *Handler) GetTransaction(c *gin.Context) {
	transactionID, err := utils.ParseULID(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, contracts.ErrorResponse{Error: "id de transação inválido"})
		return
	}

	userID, err := h.GetUserIDFromContext(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, contracts.ErrorResponse{Error: err.Error()})
		return
	}

	transactionEntity, err := h.TransactionService.GetTransactionByID(transactionID, userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) || err.Error() == "transaction does not belong to user" {
			c.JSON(http.StatusNotFound, contracts.ErrorResponse{Error: "Transação não encontrada"})
			return
		}
		c.JSON(http.StatusInternalServerError, contracts.ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, contracts.TransactionSingleResponse{Transaction: transactionEntity})
}

func (h *Handler) UpdateTransaction(c *gin.Context) {
	transactionID, err := utils.ParseULID(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, contracts.ErrorResponse{Error: "id de transação inválido"})
		return
	}

	userID, err := h.GetUserIDFromContext(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, contracts.ErrorResponse{Error: err.Error()})
		return
	}

	var body contracts.TransactionUpdateRequest
	if err = c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, contracts.ErrorResponse{Error: err.Error()})
		return
	}

	categoryID, err := utils.ParseULID(body.CategoryID)
	if err != nil {
		c.JSON(http.StatusBadRequest, contracts.ErrorResponse{Error: err.Error()})
		return
	}

	transactionEntity := transaction.Transaction{
		Id:          transactionID,
		UserId:      userID,
		CategoryId:  categoryID,
		Amount:      body.Amount,
		Description: body.Description,
		Type:        transaction.Types(body.Type),
		UpdatedAt:   utils.SetTimestamps(),
	}

	if body.Date != nil {
		transactionEntity.Date = *body.Date
	}

	if err := h.TransactionService.UpdateTransaction(&transactionEntity); err != nil {
		if err.Error() == "transaction does not belong to user" || errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, contracts.ErrorResponse{Error: "Transação não encontrada"})
			return
		}
		c.JSON(http.StatusInternalServerError, contracts.ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, contracts.MessageResponse{Message: "Transação atualizada com sucesso"})
}

func (h *Handler) DeleteTransaction(c *gin.Context) {
	transactionID, err := utils.ParseULID(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, contracts.ErrorResponse{Error: "id de transação inválido"})
		return
	}

	userID, err := h.GetUserIDFromContext(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, contracts.ErrorResponse{Error: err.Error()})
		return
	}

	if err := h.TransactionService.DeleteTransaction(transactionID, userID); err != nil {
		if err.Error() == "transaction does not exist" || err.Error() == "transaction does not belong to user" || errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, contracts.ErrorResponse{Error: "Transação não encontrada"})
			return
		}
		c.JSON(http.StatusInternalServerError, contracts.ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, contracts.MessageResponse{Message: "Transação removida com sucesso"})
}
