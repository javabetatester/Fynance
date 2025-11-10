package routes

import (
	"Fynance/internal/contracts"
	"Fynance/internal/domain/goal"
	appErrors "Fynance/internal/errors"
	"Fynance/internal/pkg"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *Handler) CreateGoal(c *gin.Context) {
	var body contracts.GoalCreateRequest
	if err := c.ShouldBindJSON(&body); err != nil {
		h.respondError(c, appErrors.ErrBadRequest.WithError(err))
		return
	}

	userID, err := h.GetUserIDFromContext(c)
	if err != nil {
		h.respondError(c, err)
		return
	}

	req := goal.GoalCreateRequest{
		UserId:  userID,
		Name:    body.Name,
		Target:  body.Target,
		EndedAt: body.EndAt,
	}

	ctx := c.Request.Context()
	if err := h.GoalService.CreateGoal(ctx, &req); err != nil {
		h.respondError(c, err)
		return
	}

	c.JSON(http.StatusCreated, contracts.MessageResponse{Message: "Meta criada com sucesso"})
}

func (h *Handler) UpdateGoal(c *gin.Context) {
	var body contracts.GoalUpdateRequest
	if err := c.ShouldBindJSON(&body); err != nil {
		h.respondError(c, appErrors.ErrBadRequest.WithError(err))
		return
	}

	userID, err := h.GetUserIDFromContext(c)
	if err != nil {
		h.respondError(c, err)
		return
	}

	id := c.Param("id")
	if id == "" {
		h.respondError(c, appErrors.NewValidationError("id", "é obrigatório"))
		return
	}

	goalID, err := pkg.ParseULID(id)
	if err != nil {
		h.respondError(c, appErrors.NewValidationError("id", "formato inválido"))
		return
	}

	req := goal.GoalUpdateRequest{
		Id:      goalID,
		UserId:  userID,
		Name:    body.Name,
		Target:  body.Target,
		EndedAt: body.EndAt,
	}

	ctx := c.Request.Context()
	if err := h.GoalService.UpdateGoal(ctx, &req); err != nil {
		h.respondError(c, err)
		return
	}

	c.JSON(http.StatusOK, contracts.MessageResponse{Message: "Meta atualizada com sucesso"})
}

func (h *Handler) ListGoals(c *gin.Context) {
	userID, err := h.GetUserIDFromContext(c)
	if err != nil {
		h.respondError(c, err)
		return
	}

	ctx := c.Request.Context()
	goals, err := h.GoalService.GetGoalsByUserID(ctx, userID)
	if err != nil {
		h.respondError(c, err)
		return
	}

	c.JSON(http.StatusOK, contracts.GoalListResponse{Goals: goals, Total: len(goals)})
}

func (h *Handler) GetGoal(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		h.respondError(c, appErrors.NewValidationError("id", "é obrigatório"))
		return
	}

	goalID, err := pkg.ParseULID(id)
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
	goalEntity, err := h.GoalService.GetGoalByID(ctx, goalID, userID)
	if err != nil {
		h.respondError(c, err)
		return
	}

	c.JSON(http.StatusOK, contracts.GoalResponse{Goal: goalEntity})
}

func (h *Handler) DeleteGoal(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		h.respondError(c, appErrors.NewValidationError("id", "é obrigatório"))
		return
	}

	goalID, err := pkg.ParseULID(id)
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
	if err := h.GoalService.DeleteGoal(ctx, goalID, userID); err != nil {
		h.respondError(c, err)
		return
	}

	c.JSON(http.StatusOK, contracts.MessageResponse{Message: "Meta removida com sucesso"})
}
