package routes

import (
	"Fynance/internal/domain/auth"
	"Fynance/internal/domain/goal"
	"Fynance/internal/domain/investment"
	"Fynance/internal/domain/transaction"
	"Fynance/internal/domain/user"
	appErrors "Fynance/internal/errors"
	"Fynance/internal/logger"
	"Fynance/internal/middleware"
	"Fynance/internal/pkg"

	"github.com/gin-gonic/gin"
	"github.com/oklog/ulid/v2"
)

type Handler struct {
	UserService        user.Service
	AuthService        auth.Service
	JwtService         *middleware.JwtService
	TransactionService transaction.Service
	GoalService        goal.Service
	InvestmentService  investment.Service
}

func (h *Handler) GetUserIDFromContext(c *gin.Context) (ulid.ULID, error) {
	userIDStr, exists := c.Get("user_id")
	if !exists {
		return ulid.ULID{}, appErrors.ErrUnauthorized
	}

	userID, err := pkg.ParseULID(userIDStr.(string))
	if err != nil {
		return ulid.ULID{}, appErrors.ErrUnauthorized.WithError(err)
	}

	return userID, nil
}

func (h *Handler) respondError(c *gin.Context, err error) {
	appErr := appErrors.FromError(err)
	event := logger.Error().Str("code", appErr.Code).Str("path", c.FullPath())
	if appErr.Err != nil {
		event = event.Err(appErr.Err)
	}
	event.Msg("request_error")
	payload := gin.H{
		"error":   appErr.Code,
		"message": appErr.Message,
	}
	if len(appErr.Details) > 0 {
		payload["details"] = appErr.Details
	}
	c.JSON(appErr.StatusCode, payload)
}
