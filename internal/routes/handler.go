package routes

import (
	"Fynance/internal/domain/auth"
	"Fynance/internal/domain/goal"
	"Fynance/internal/domain/transaction"
	"Fynance/internal/domain/user"
	"Fynance/internal/utils"
	"errors"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type Handler struct {
	UserService        user.Service
	AuthService        auth.Service
	JwtService         *utils.JwtService
	TransactionService transaction.Service
	GoalService        goal.Service
}

func (h *Handler) GetUserIDFromContext(c *gin.Context) (uuid.UUID, error) {
	userIDStr, exists := c.Get("user_id")
	if !exists {
		return uuid.Nil, errors.New("user not authenticated")
	}

	// Se o claims.Sub já é string UUID
	userID, err := uuid.Parse(userIDStr.(string))
	if err != nil {
		return uuid.Nil, errors.New("invalid user ID format")
	}

	return userID, nil
}
