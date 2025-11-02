package routes

import (
	"Fynance/internal/domain/investment"
	"Fynance/internal/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *Handler) CreateInvestment(c *gin.Context) {
	var req investment.CreateInvestmentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID, err := h.GetUserIDFromContext(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	req.UserId = userID

	inv, err := h.InvestmentService.CreateInvestment(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message":    "Investment created successfully",
		"investment": inv,
	})
}

func (h *Handler) ListInvestments(c *gin.Context) {
	userID, err := h.GetUserIDFromContext(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	investments, err := h.InvestmentService.ListInvestments(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"investments": investments,
		"total":       len(investments),
	})
}

func (h *Handler) GetInvestment(c *gin.Context) {
	investmentID, err := utils.ParseULID(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid investment id"})
		return
	}

	userID, err := h.GetUserIDFromContext(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	inv, err := h.InvestmentService.GetInvestment(investmentID, userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "investment not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"investment": inv,
	})
}

func (h *Handler) MakeContribution(c *gin.Context) {
	investmentID, err := utils.ParseULID(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid investment id"})
		return
	}

	userID, err := h.GetUserIDFromContext(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	var req investment.ContributionRequest
	if errs := c.ShouldBindJSON(&req); errs != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": errs.Error()})
		return
	}

	err = h.InvestmentService.MakeContribution(investmentID, userID, req.Amount, req.CategoryId, req.Description)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Contribution made successfully",
	})
}

func (h *Handler) MakeWithdrawal(c *gin.Context) {
	investmentID, err := utils.ParseULID(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid investment id"})
		return
	}

	userID, err := h.GetUserIDFromContext(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	var req investment.WithdralRequest
	if errs := c.ShouldBindJSON(&req); errs != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": errs.Error()})
		return
	}

	err = h.InvestmentService.MakeWithdrawal(investmentID, userID, req.Amount, req.CategoryId, req.Description)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Withdrawal made successfully",
	})
}

func (h *Handler) GetInvestmentReturn(c *gin.Context) {
	investmentID, err := utils.ParseULID(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid investment id"})
		return
	}

	userID, err := h.GetUserIDFromContext(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	profit, returnPercentage, err := h.InvestmentService.CalculateReturn(investmentID, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"profit":            profit,
		"return_percentage": returnPercentage,
	})
}

func (h *Handler) DeleteInvestment(c *gin.Context) {
	investmentID, err := utils.ParseULID(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid investment id"})
		return
	}

	userID, err := h.GetUserIDFromContext(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	err = h.InvestmentService.DeleteInvestment(investmentID, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Investment deleted successfully",
	})
}
