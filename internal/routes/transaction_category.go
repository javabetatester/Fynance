package routes

import (
	"Fynance/internal/contracts"
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
// @Param        category body contracts.CategoryCreateRequest true "Dados da categoria"
// @Success      201 {object} contracts.MessageResponse "Categoria criada com sucesso"
// @Failure      400 {object} contracts.ErrorResponse "Erro de validação"
// @Failure      401 {object} contracts.ErrorResponse "Não autorizado"
// @Failure      500 {object} contracts.ErrorResponse "Erro interno do servidor"
// @Router       /api/categories [post]
// @Security     BearerAuth
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
