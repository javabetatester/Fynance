package routes

import (
	"Fynance/internal/domain/auth"
	"Fynance/internal/domain/goal"
	"Fynance/internal/domain/transaction"
	"Fynance/internal/domain/user"
	"Fynance/internal/middleware"
	"Fynance/internal/utils"
	"errors"

	"github.com/gin-gonic/gin"
	"github.com/oklog/ulid/v2"
)

type Handler struct {
	UserService        user.Service
	AuthService        auth.Service
	JwtService         *middleware.JwtService
	TransactionService transaction.Service
	GoalService        goal.Service
}

func (h *Handler) GetUserIDFromContext(c *gin.Context) (ulid.ULID, error) {
	userIDStr, exists := c.Get("user_id")
	if !exists {
		return ulid.ULID{}, errors.New("user not authenticated")
	}

	userID, err := utils.ParseULID(userIDStr.(string))
	if err != nil {
		return ulid.ULID{}, err
	}

	return userID, nil
}
