package routes

import (
	"Fynance/internal/domain/goal"
	"Fynance/internal/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

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