package routes

import (
	"errors"
	"net/http"

	"Fynance/internal/contracts"
	"Fynance/internal/domain/goal"
	"Fynance/internal/utils"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

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

	ctx := c.Request.Context()
	if err := h.GoalService.CreateGoal(ctx, &req); err != nil {
		c.JSON(http.StatusInternalServerError, contracts.ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusCreated, contracts.MessageResponse{Message: "Meta criada com sucesso"})
}

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

	ctx := c.Request.Context()
	if err := h.GoalService.UpdateGoal(ctx, &req); err != nil {
		c.JSON(http.StatusInternalServerError, contracts.ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, contracts.MessageResponse{Message: "Meta atualizada com sucesso"})
}

func (h *Handler) ListGoals(c *gin.Context) {
	userID, err := h.GetUserIDFromContext(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, contracts.ErrorResponse{Error: err.Error()})
		return
	}

	ctx := c.Request.Context()
	goals, err := h.GoalService.GetGoalsByUserID(ctx, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, contracts.ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, contracts.GoalListResponse{Goals: goals, Total: len(goals)})
}

func (h *Handler) GetGoal(c *gin.Context) {
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

	userID, err := h.GetUserIDFromContext(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, contracts.ErrorResponse{Error: err.Error()})
		return
	}

	ctx := c.Request.Context()
	goalEntity, err := h.GoalService.GetGoalByID(ctx, goalID, userID)
	if err != nil {
		if err.Error() == "goal does not belong to user" || errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, contracts.ErrorResponse{Error: "Meta não encontrada"})
			return
		}
		c.JSON(http.StatusInternalServerError, contracts.ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, contracts.GoalResponse{Goal: goalEntity})
}

func (h *Handler) DeleteGoal(c *gin.Context) {
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

	userID, err := h.GetUserIDFromContext(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, contracts.ErrorResponse{Error: err.Error()})
		return
	}

	ctx := c.Request.Context()
	if err := h.GoalService.DeleteGoal(ctx, goalID, userID); err != nil {
		if err.Error() == "goal does not belong to user" || errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, contracts.ErrorResponse{Error: "Meta não encontrada"})
			return
		}
		c.JSON(http.StatusInternalServerError, contracts.ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, contracts.MessageResponse{Message: "Meta removida com sucesso"})
}
