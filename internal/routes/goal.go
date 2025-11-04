package routes

import (
	"Fynance/internal/contracts"
	"Fynance/internal/domain/goal"
	"Fynance/internal/utils"
	"net/http"

	"github.com/gin-gonic/gin"
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

	if err := h.GoalService.CreateGoal(&req); err != nil {
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

	if err := h.GoalService.UpdateGoal(&req); err != nil {
		c.JSON(http.StatusInternalServerError, contracts.ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, contracts.MessageResponse{Message: "Meta atualizada com sucesso"})
}
