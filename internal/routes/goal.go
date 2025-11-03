package routes

import (
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
// @Param        goal body object true "Dados da meta"
// @Success      201 {string} string "Meta criada com sucesso"
// @Failure      400 {object} map[string]string "Erro de validação"
// @Failure      401 {object} map[string]string "Não autorizado"
// @Failure      500 {object} map[string]string "Erro interno do servidor"
// @Router       /api/goals [post]
// @Security     BearerAuth
func (h *Handler) CreateGoal(c *gin.Context) {
	var goal goal.GoalCreateRequest
	if err := c.ShouldBindJSON(&goal); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

    userID, err := h.GetUserIDFromContext(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}
    goal.UserId = userID

	if err := h.GoalService.CreateGoal(&goal); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, "Goal created with success")
}

// UpdateGoal godoc
// @Summary      Atualizar meta financeira
// @Description  Atualiza uma meta financeira existente do usuário autenticado
// @Tags         goals
// @Accept       json
// @Produce      json
// @Param        id path string true "ID da meta"
// @Param        goal body object true "Dados atualizados da meta"
// @Success      200 {string} string "Meta atualizada com sucesso"
// @Failure      400 {object} map[string]string "Erro de validação"
// @Failure      401 {object} map[string]string "Não autorizado"
// @Failure      500 {object} map[string]string "Erro interno do servidor"
// @Router       /api/goals/{id} [patch]
// @Security     BearerAuth
func (h *Handler) UpdateGoal(c *gin.Context) {
	var goal goal.GoalUpdateRequest
	if err := c.ShouldBindJSON(&goal); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

    userID, err := h.GetUserIDFromContext(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}
    goal.UserId = userID

	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id is required"})
		return
	}

	goal.Id, err = utils.ParseULID(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	
	if err := h.GoalService.UpdateGoal(&goal); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, "Goal created with success")
}