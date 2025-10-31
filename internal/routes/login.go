package routes

import (
	"Fynance/internal/domain/login"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *Handler) Login(c *gin.Context) {
	var login login.Login
	if err := c.ShouldBindJSON(&login); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := h.LoginService.Login(login)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	token, err := h.JwtService.GenerateToken(user.Id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Login successful",
		"user":    user.Name,
		"token":   token,
	})
}
