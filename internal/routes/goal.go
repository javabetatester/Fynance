package routes

import (
	"Fynance/internal/contracts"
	"Fynance/internal/domain/goal"
	"Fynance/internal/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

// CreateGoal godoc
// @Summary      Criar meta financeira
// @Description  Cria uma nova meta financeira para o usuário autenticado
// @Tags         goals
// @Accept       json
// @Produce      json
// @Param        goal body contracts.GoalCreateRequest true "Dados da meta"
// @Success      201 {object} contracts.MessageResponse "Meta criada com sucesso"
// @Failure      400 {object} contracts.ErrorResponse "Erro de validação"
// @Failure      401 {object} contracts.ErrorResponse "Não autorizado"
// @Failure      500 {object} contracts.ErrorResponse "Erro interno do servidor"
// @Router       /api/goals [post]
// @Security     BearerAuth
func (h *Handler) CreateGoal(c *gin.Context) {
	var body contracts.GoalCreateRequest
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, contracts.ErrorResponse{Error: err.Error()})
		return
	}

	userID, err := h.GetUserIDFromContext(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, contracts.ErrorResponse{Error: err.Error()})
		return
	}

	req := goal.GoalCreateRequest{
		UserId:  userID,
		Name:    body.Name,
		Target:  body.Target,
		EndedAt: body.EndAt,
	}

	if err := h.GoalService.CreateGoal(&req); err != nil {
		c.JSON(http.StatusInternalServerError, contracts.ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusCreated, contracts.MessageResponse{Message: "Meta criada com sucesso"})
}

// UpdateGoal godoc
// @Summary      Atualizar meta financeira
// @Description  Atualiza uma meta financeira existente do usuário autenticado
// @Tags         goals
// @Accept       json
// @Produce      json
// @Param        id path string true "ID da meta"
// @Param        goal body contracts.GoalUpdateRequest true "Dados atualizados da meta"
// @Success      200 {object} contracts.MessageResponse "Meta atualizada com sucesso"
// @Failure      400 {object} contracts.ErrorResponse "Erro de validação"
// @Failure      401 {object} contracts.ErrorResponse "Não autorizado"
// @Failure      500 {object} contracts.ErrorResponse "Erro interno do servidor"
// @Router       /api/goals/{id} [patch]
// @Security     BearerAuth
func (h *Handler) UpdateGoal(c *gin.Context) {
	var body contracts.GoalUpdateRequest
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, contracts.ErrorResponse{Error: err.Error()})
		return
	}

	userID, err := h.GetUserIDFromContext(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, contracts.ErrorResponse{Error: err.Error()})
		return
	}

	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, contracts.ErrorResponse{Error: "id é obrigatório"})
		return
	}

	goalID, err := utils.ParseULID(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, contracts.ErrorResponse{Error: err.Error()})
		return
	}

	req := goal.GoalUpdateRequest{
		Id:      goalID,
		UserId:  userID,
		Name:    body.Name,
		Target:  body.Target,
		EndedAt: body.EndAt,
	}

	if err := h.GoalService.UpdateGoal(&req); err != nil {
		c.JSON(http.StatusInternalServerError, contracts.ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, contracts.MessageResponse{Message: "Meta atualizada com sucesso"})
}
