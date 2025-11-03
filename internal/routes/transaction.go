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

// CreateTransaction godoc
// @Summary      Criar transação
// @Description  Cria uma nova transação financeira para o usuário autenticado
// @Tags         transactions
// @Accept       json
// @Produce      json
// @Param        transaction body object true "Dados da transação"
// @Success      201 {object} map[string]interface{} "Transação criada com sucesso"
// @Failure      400 {object} map[string]string "Erro de validação"
// @Failure      401 {object} map[string]string "Não autorizado"
// @Failure      500 {object} map[string]string "Erro interno do servidor"
// @Router       /api/transactions [post]
// @Security     BearerAuth
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

// GetTransactions godoc
// @Summary      Listar transações
// @Description  Lista todas as transações do usuário autenticado
// @Tags         transactions
// @Accept       json
// @Produce      json
// @Success      200 {array} object "Lista de transações"
// @Failure      401 {object} map[string]string "Não autorizado"
// @Failure      500 {object} map[string]string "Erro interno do servidor"
// @Router       /api/transactions [get]
// @Security     BearerAuth
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
