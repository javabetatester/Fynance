package middleware

import (
	"net/http"
	"strings"

	"Fynance/internal/logger"

	"github.com/gin-gonic/gin"
)

func AuthMiddleware(jwtService *JwtService) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			logger.Warn().Str("path", c.FullPath()).Msg("authorization header ausente")
			c.JSON(http.StatusUnauthorized, gin.H{"error": "auth_required"})
			c.Abort()
			return
		}

		if !strings.HasPrefix(authHeader, "Bearer ") {
			logger.Warn().Str("path", c.FullPath()).Msg("authorization header inválido")
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid_token"})
			c.Abort()
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

		claims, err := jwtService.ParseToken(tokenString)
		if err != nil {
			logger.Warn().Str("path", c.FullPath()).Msg("token inválido")
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid_token"})
			c.Abort()
			return
		}

		c.Set("user_id", claims.Sub)
		c.Set("plan", claims.Plan)

		c.Next()
	}
}
