package routes

import (
	"Fynance/internal/contracts"
	"Fynance/internal/domain/transaction"
	appErrors "Fynance/internal/errors"
	"Fynance/internal/pkg"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *Handler) CreateTransaction(c *gin.Context) {
	var body contracts.TransactionCreateRequest

	if err := c.ShouldBindJSON(&body); err != nil {
		h.respondError(c, appErrors.ErrBadRequest.WithError(err))
		return
	}

	userID, err := h.GetUserIDFromContext(c)
	if err != nil {
		h.respondError(c, err)
		return
	}

	categoryID, err := pkg.ParseULID(body.CategoryID)
	if err != nil {
		h.respondError(c, appErrors.NewValidationError("category_id", "formato inválido"))
		return
	}

	transactionEntity := transaction.Transaction{
		Type:        transaction.Types(body.Type),
		UserId:      userID,
		CategoryId:  categoryID,
		Amount:      body.Amount,
		Description: body.Description,
		Date:        pkg.SetTimestamps(),
	}

	ctx := c.Request.Context()
	if err := h.TransactionService.CreateTransaction(ctx, &transactionEntity); err != nil {
		h.respondError(c, err)
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
		h.respondError(c, err)
		return
	}

	ctx := c.Request.Context()
	transactions, err := h.TransactionService.GetAllTransactions(ctx, userID)
	if err != nil {
		h.respondError(c, err)
		return
	}

	c.JSON(http.StatusOK, contracts.TransactionListResponse{Transactions: transactions, Total: len(transactions)})
}

func (h *Handler) GetTransaction(c *gin.Context) {
	transactionID, err := pkg.ParseULID(c.Param("id"))
	if err != nil {
		h.respondError(c, appErrors.NewValidationError("id", "formato inválido"))
		return
	}

	userID, err := h.GetUserIDFromContext(c)
	if err != nil {
		h.respondError(c, err)
		return
	}

	ctx := c.Request.Context()
	transactionEntity, err := h.TransactionService.GetTransactionByID(ctx, transactionID, userID)
	if err != nil {
		h.respondError(c, err)
		return
	}

	c.JSON(http.StatusOK, contracts.TransactionSingleResponse{Transaction: transactionEntity})
}

func (h *Handler) UpdateTransaction(c *gin.Context) {
	transactionID, err := pkg.ParseULID(c.Param("id"))
	if err != nil {
		h.respondError(c, appErrors.NewValidationError("id", "formato inválido"))
		return
	}

	userID, err := h.GetUserIDFromContext(c)
	if err != nil {
		h.respondError(c, err)
		return
	}

	var body contracts.TransactionUpdateRequest
	if err = c.ShouldBindJSON(&body); err != nil {
		h.respondError(c, appErrors.ErrBadRequest.WithError(err))
		return
	}

	categoryID, err := pkg.ParseULID(body.CategoryID)
	if err != nil {
		h.respondError(c, appErrors.NewValidationError("category_id", "formato inválido"))
		return
	}

	transactionEntity := transaction.Transaction{
		Id:          transactionID,
		UserId:      userID,
		CategoryId:  categoryID,
		Amount:      body.Amount,
		Description: body.Description,
		Type:        transaction.Types(body.Type),
		UpdatedAt:   pkg.SetTimestamps(),
	}

	if body.Date != nil {
		transactionEntity.Date = *body.Date
	}

	ctx := c.Request.Context()
	if err := h.TransactionService.UpdateTransaction(ctx, &transactionEntity); err != nil {
		h.respondError(c, err)
		return
	}

	c.JSON(http.StatusOK, contracts.MessageResponse{Message: "Transação atualizada com sucesso"})
}

func (h *Handler) DeleteTransaction(c *gin.Context) {
	transactionID, err := pkg.ParseULID(c.Param("id"))
	if err != nil {
		h.respondError(c, appErrors.NewValidationError("id", "formato inválido"))
		return
	}

	userID, err := h.GetUserIDFromContext(c)
	if err != nil {
		h.respondError(c, err)
		return
	}

	ctx := c.Request.Context()
	if err := h.TransactionService.DeleteTransaction(ctx, transactionID, userID); err != nil {
		h.respondError(c, err)
		return
	}

	c.JSON(http.StatusOK, contracts.MessageResponse{Message: "Transação removida com sucesso"})
}
