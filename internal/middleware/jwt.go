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

		hasValidation := false

		idFromURL := c.Param("id")
		if idFromURL != "" {

			if !strings.EqualFold(tokenUserID, idFromURL) {
				c.JSON(http.StatusNotFound, gin.H{"error": "Resource not found"})
				c.Abort()
				return
			}
			hasValidation = true
		}

		if c.Request.Body != nil {
			bodyBytes, err := io.ReadAll(c.Request.Body)
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": "Cannot read body"})
				c.Abort()
				return
			}

			c.Request.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

			if len(bodyBytes) > 0 {
				var bodyData map[string]interface{}
				if err := json.Unmarshal(bodyBytes, &bodyData); err != nil {
					if hasValidation {
						c.Next()
						return
					}
					c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON"})
					c.Abort()
					return
				}

				var idFromBody interface{}
				var foundIDField bool

				if val, exists := bodyData["id"]; exists {
					idFromBody = val
					foundIDField = true
				} else if val, exists := bodyData["Id"]; exists {
					idFromBody = val
					foundIDField = true
				} else if val, exists := bodyData["user_id"]; exists {
					idFromBody = val
					foundIDField = true
				}

				if foundIDField {
					idStr, ok := idFromBody.(string)
					if !ok {
						c.JSON(http.StatusBadRequest, gin.H{"error": "ID must be a string"})
						c.Abort()
						return
					}

					if tokenUserID != idStr {
						c.JSON(http.StatusNotFound, gin.H{"error": "Resource not found"})
						c.Abort()
						return
					}
				}
			}
		}

		c.Next()
	}
}
