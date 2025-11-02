package routes

import (
	"Fynance/internal/domain/investment"
	"Fynance/internal/utils"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/oklog/ulid/v2"
)

func (h *Handler) CreateInvestment(c *gin.Context) {
	type bodyDTO struct {
		Type             investment.Types `json:"type" binding:"required"`
		Name             string           `json:"name" binding:"required"`
		InitialAmount    float64          `json:"initial_amount"`
		InitialAmountAlt float64          `json:"initialAmount"`
		ReturnRate       float64          `json:"return_rate"`
		ReturnRateAlt    float64          `json:"returnRate"`
		CategoryId       string           `json:"category_id"`
		CategoryIdAlt    string           `json:"categoryId"`
	}

	var body bodyDTO
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID, err := h.GetUserIDFromContext(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	initial := body.InitialAmount
	if initial == 0 && body.InitialAmountAlt != 0 {
		initial = body.InitialAmountAlt
	}
	if initial <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "initial amount must be greater than 0"})
		return
	}

	ret := body.ReturnRate
	if ret == 0 && body.ReturnRateAlt != 0 {
		ret = body.ReturnRateAlt
	}

	catIDStr := body.CategoryId
	if catIDStr == "" {
		catIDStr = body.CategoryIdAlt
	}
	var categoryID ulid.ULID
	if catIDStr == "" {
		categoryID, err = h.TransactionService.EnsureDefaultInvestmentCategory(userID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	} else {
		categoryID, err = utils.ParseULID(catIDStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid category id"})
			return
		}
		if errs := h.TransactionService.CategoryValidation(categoryID, userID); errs != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": errs.Error()})
			return
		}
	}

	req := investment.CreateInvestmentRequest{
		UserId:        userID,
		Type:          body.Type,
		Name:          body.Name,
		InitialAmount: initial,
		ReturnRate:    ret,
		CategoryId:    categoryID,
	}

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

	// if errs := h.TransactionService.CategoryValidation(req.CategoryId, userID); errs != nil {
	// 	c.JSON(http.StatusBadRequest, gin.H{"error": errs.Error()})
	// 	return
	// }

	err = h.InvestmentService.MakeContribution(investmentID, userID, req.Amount, req.CategoryId, req.Description)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Contribution made successfully",
	})
}

func (h *Handler) MakeWithdraw(c *gin.Context) {
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

	// if errs := h.TransactionService.CategoryValidation(req.CategoryId, userID); errs != nil {
	// 	c.JSON(http.StatusBadRequest, gin.H{"error": errs.Error()})
	// 	return
	// }

	err = h.InvestmentService.MakeWithdraw(investmentID, userID, req.Amount, req.CategoryId, req.Description)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Withdraw made successfully",
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
