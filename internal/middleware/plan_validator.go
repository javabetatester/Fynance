package middleware

import (
	"net/http"

	"Fynance/internal/domain/user"
	appErrors "Fynance/internal/errors"

	"github.com/gin-gonic/gin"
)

func respondPlan(c *gin.Context, err *appErrors.AppError) {
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

func RequirePlan(allowedPlans ...user.Plan) gin.HandlerFunc {
	plans := make(map[user.Plan]struct{}, len(allowedPlans))
	for _, plan := range allowedPlans {
		plans[plan] = struct{}{}
	}

	return func(c *gin.Context) {
		planValue, exists := c.Get("plan")
		if !exists {
			err := appErrors.WrapError(nil, appErrors.ErrForbidden.Code, "Plano não encontrado para o usuário autenticado", http.StatusForbidden)
			respondPlan(c, err)
			return
		}

		plan, ok := planValue.(user.Plan)
		if !ok {
			err := appErrors.WrapError(nil, appErrors.ErrForbidden.Code, "Plano inválido para o usuário autenticado", http.StatusForbidden)
			respondPlan(c, err)
			return
		}

		if len(plans) == 0 {
			c.Next()
			return
		}

		if _, allowed := plans[plan]; !allowed {
			err := appErrors.WrapError(nil, appErrors.ErrForbidden.Code, "Plano incompatível com a operação solicitada", http.StatusForbidden)
			respondPlan(c, err)
			return
		}

		c.Next()
	}
}
