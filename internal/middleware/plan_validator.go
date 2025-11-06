package middleware

import (
	"net/http"

	"Fynance/internal/domain/user"

	"github.com/gin-gonic/gin"
)

func RequirePlan(allowedPlans ...user.Plan) gin.HandlerFunc {
	plans := make(map[user.Plan]struct{}, len(allowedPlans))
	for _, plan := range allowedPlans {
		plans[plan] = struct{}{}
	}

	return func(c *gin.Context) {
		planValue, exists := c.Get("plan")
		if !exists {
			c.JSON(http.StatusForbidden, gin.H{"error": "Plano nao encontrado para o usuario autenticado"})
			c.Abort()
			return
		}

		plan, ok := planValue.(user.Plan)
		if !ok {
			c.JSON(http.StatusForbidden, gin.H{"error": "Plano invalido para o usuario autenticado"})
			c.Abort()
			return
		}

		if len(plans) == 0 {
			c.Next()
			return
		}

		if _, allowed := plans[plan]; !allowed {
			c.JSON(http.StatusForbidden, gin.H{"error": "Plano incompativel com a operacao solicitada"})
			c.Abort()
			return
		}

		c.Next()
	}
}
