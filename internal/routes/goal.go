package routes

import (
	"Fynance/internal/domain/goal"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *Handler) CreateGoal(c *gin.Context) {
	var goal goal.Goal
	if err := c.ShouldBindJSON(&goal); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID, err := h.GetUserIDFromContext(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}
	goal.UserId = userID.String()

	if err := h.GoalService.CreateGoal(&goal); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, "Goal created with success")
}
