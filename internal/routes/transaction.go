package routes

import (
	"Fynance/internal/contracts"
	"Fynance/internal/domain/transaction"
	"Fynance/internal/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

// CreateTransaction godoc
// @Summary      Criar transação
// @Description  Cria uma nova transação financeira para o usuário autenticado
// @Tags         transactions
// @Accept       json
// @Produce      json
// @Param        transaction body contracts.TransactionCreateRequest true "Dados da transação"
// @Success      201 {object} contracts.TransactionCreateResponse "Transação criada com sucesso"
// @Failure      400 {object} contracts.ErrorResponse "Erro de validação"
// @Failure      401 {object} contracts.ErrorResponse "Não autorizado"
// @Failure      500 {object} contracts.ErrorResponse "Erro interno do servidor"
// @Router       /api/transactions [post]
// @Security     BearerAuth
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

// GetTransactions godoc
// @Summary      Listar transações
// @Description  Lista todas as transações do usuário autenticado
// @Tags         transactions
// @Accept       json
// @Produce      json
// @Success      200 {array} transaction.Transaction "Lista de transações"
// @Failure      401 {object} contracts.ErrorResponse "Não autorizado"
// @Failure      500 {object} contracts.ErrorResponse "Erro interno do servidor"
// @Router       /api/transactions [get]
// @Security     BearerAuth
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
