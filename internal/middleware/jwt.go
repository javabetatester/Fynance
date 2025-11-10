package middleware

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"strings"

	appErrors "Fynance/internal/errors"

	"github.com/gin-gonic/gin"
)

func respondOwnership(c *gin.Context, err *appErrors.AppError) {
	payload := gin.H{
		"error":   err.Code,
		"message": err.Message,
	}
	if len(err.Details) > 0 {
		payload["details"] = err.Details
	}
	c.JSON(err.StatusCode, payload)
	c.Abort()
}

func RequireOwnership() gin.HandlerFunc {
	return func(c *gin.Context) {
		userIDFromToken, exists := c.Get("user_id")
		if !exists {
			respondOwnership(c, appErrors.ErrUnauthorized)
			return
		}

		tokenUserID, ok := userIDFromToken.(string)
		if !ok || tokenUserID == "" {
			err := appErrors.ErrUnauthorized.WithDetails(map[string]interface{}{"reason": "user_id_invalid"})
			respondOwnership(c, err)
			return
		}

		for _, param := range c.Params {
			if strings.EqualFold(param.Key, "user_id") || strings.EqualFold(param.Key, "userid") {
				if !strings.EqualFold(tokenUserID, param.Value) {
					respondOwnership(c, appErrors.ErrResourceNotOwned)
					return
				}
			}
		}

		if userIDQuery := c.Query("user_id"); userIDQuery != "" {
			if !strings.EqualFold(tokenUserID, userIDQuery) {
				respondOwnership(c, appErrors.ErrResourceNotOwned)
				return
			}
		}

		if c.Request.Body == nil || c.Request.Body == http.NoBody {
			c.Next()
			return
		}

		bodyBytes, err := io.ReadAll(c.Request.Body)
		if err != nil {
			respondOwnership(c, appErrors.ErrBadRequest.WithError(err))
			return
		}

		c.Request.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

		if len(bodyBytes) == 0 {
			c.Next()
			return
		}

		var bodyData map[string]interface{}
		if err := json.Unmarshal(bodyBytes, &bodyData); err != nil {
			respondOwnership(c, appErrors.ErrBadRequest.WithError(err))
			return
		}

		for _, key := range []string{"user_id", "userId"} {
			if val, ok := bodyData[key]; ok {
				idStr, ok := val.(string)
				if !ok {
					err := appErrors.NewValidationError("user_id", "deve ser uma string")
					respondOwnership(c, err)
					return
				}
				if !strings.EqualFold(tokenUserID, idStr) {
					respondOwnership(c, appErrors.ErrResourceNotOwned)
					return
				}
			}
		}

		c.Next()
	}
}
