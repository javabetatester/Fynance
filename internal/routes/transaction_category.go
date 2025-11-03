package routes

import (
	"Fynance/internal/domain/transaction"
	"net/http"

	"github.com/gin-gonic/gin"
)

// CreateCategory godoc
// @Summary      Criar categoria de transação
// @Description  Cria uma nova categoria de transação para o usuário autenticado
// @Tags         categories
// @Accept       json
// @Produce      json
// @Param        category body object true "Dados da categoria"
// @Success      201 {string} string "Categoria criada com sucesso"
// @Failure      400 {object} map[string]string "Erro de validação"
// @Failure      401 {object} map[string]string "Não autorizado"
// @Failure      500 {object} map[string]string "Erro interno do servidor"
// @Router       /api/categories [post]
// @Security     BearerAuth
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

	category.UserId = userID

	if err := h.TransactionService.CreateCategory(&category); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, "Category created with success")
}
