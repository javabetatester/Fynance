package middleware

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func RequireOwnership() gin.HandlerFunc {
	return func(c *gin.Context) {
		userIDFromToken, exists := c.Get("user_id")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			c.Abort()
			return
		}

		tokenUserID := userIDFromToken.(string)

		for _, param := range c.Params {
			if strings.EqualFold(param.Key, "user_id") || strings.EqualFold(param.Key, "userid") {
				if !strings.EqualFold(tokenUserID, param.Value) {
					c.JSON(http.StatusForbidden, gin.H{"error": "Forbidden"})
					c.Abort()
					return
				}
			}
		}

		if userIDQuery := c.Query("user_id"); userIDQuery != "" {
			if !strings.EqualFold(tokenUserID, userIDQuery) {
				c.JSON(http.StatusForbidden, gin.H{"error": "Forbidden"})
				c.Abort()
				return
			}
		}

		if c.Request.Body == nil || c.Request.Body == http.NoBody {
			c.Next()
			return
		}

		bodyBytes, err := io.ReadAll(c.Request.Body)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Cannot read body"})
			c.Abort()
			return
		}

		c.Request.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

		if len(bodyBytes) == 0 {
			c.Next()
			return
		}

		var bodyData map[string]interface{}
		if err := json.Unmarshal(bodyBytes, &bodyData); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON"})
			c.Abort()
			return
		}

		for _, key := range []string{"user_id", "userId"} {
			if val, ok := bodyData[key]; ok {
				idStr, ok := val.(string)
				if !ok {
					c.JSON(http.StatusBadRequest, gin.H{"error": "user_id must be a string"})
					c.Abort()
					return
				}
				if !strings.EqualFold(tokenUserID, idStr) {
					c.JSON(http.StatusForbidden, gin.H{"error": "Forbidden"})
					c.Abort()
					return
				}
			}
		}

		c.Next()
	}
}
