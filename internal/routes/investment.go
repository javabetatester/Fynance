package routes

import (
	"net/http"

	"Fynance/internal/contracts"
	"Fynance/internal/domain/investment"
	appErrors "Fynance/internal/errors"
	"Fynance/internal/pkg"

	"github.com/gin-gonic/gin"
)

func (h *Handler) CreateInvestment(c *gin.Context) {
	var body contracts.InvestmentCreateRequest
	if err := c.ShouldBindJSON(&body); err != nil {
		h.respondError(c, appErrors.ErrBadRequest.WithError(err))
		return
	}

	userID, err := h.GetUserIDFromContext(c)
	if err != nil {
		h.respondError(c, err)
		return
	}

	req := investment.CreateInvestmentRequest{
		UserId:        userID,
		Type:          investment.Types(body.Type),
		Name:          body.Name,
		InitialAmount: body.InitialAmount,
		ReturnRate:    body.ReturnRate,
	}

	ctx := c.Request.Context()
	inv, err := h.InvestmentService.CreateInvestment(ctx, req)
	if err != nil {
		h.respondError(c, err)
		return
	}

	c.JSON(http.StatusCreated, contracts.InvestmentCreateResponse{
		Message:    "Investimento criado com sucesso",
		Investment: *inv,
	})
}

func (h *Handler) ListInvestments(c *gin.Context) {
	userID, err := h.GetUserIDFromContext(c)
	if err != nil {
		h.respondError(c, err)
		return
	}

	ctx := c.Request.Context()
	investments, err := h.InvestmentService.ListInvestments(ctx, userID)
	if err != nil {
		h.respondError(c, err)
		return
	}

	c.JSON(http.StatusOK, contracts.InvestmentListResponse{
		Total:       len(investments),
		Investments: investments,
	})
}

func (h *Handler) GetInvestment(c *gin.Context) {
	investmentID, err := pkg.ParseULID(c.Param("id"))
	if err != nil {
		h.respondError(c, appErrors.NewValidationError("id", "formato inválido"))
		return
	}

	userID, err := h.GetUserIDFromContext(c)
	if err != nil {
		h.respondError(c, err)
		return
	}

	ctx := c.Request.Context()
	inv, err := h.InvestmentService.GetInvestment(ctx, investmentID, userID)
	if err != nil {
		h.respondError(c, err)
		return
	}

	c.JSON(http.StatusOK, contracts.InvestmentSingleResponse{Investment: inv})
}

func (h *Handler) MakeContribution(c *gin.Context) {
	investmentID, err := pkg.ParseULID(c.Param("id"))
	if err != nil {
		h.respondError(c, appErrors.NewValidationError("id", "formato inválido"))
		return
	}

	userID, err := h.GetUserIDFromContext(c)
	if err != nil {
		h.respondError(c, err)
		return
	}

	var body contracts.InvestmentContributionRequest
	if errs := c.ShouldBindJSON(&body); errs != nil {
		h.respondError(c, appErrors.ErrBadRequest.WithError(errs))
		return
	}

	ctx := c.Request.Context()
	if err := h.InvestmentService.MakeContribution(ctx, investmentID, userID, body.Amount, body.Description); err != nil {
		h.respondError(c, err)
		return
	}

	c.JSON(http.StatusOK, contracts.MessageResponse{Message: "Aporte registrado com sucesso"})
}

func (h *Handler) MakeWithdraw(c *gin.Context) {
	investmentID, err := pkg.ParseULID(c.Param("id"))
	if err != nil {
		h.respondError(c, appErrors.NewValidationError("id", "formato inválido"))
		return
	}

	userID, err := h.GetUserIDFromContext(c)
	if err != nil {
		h.respondError(c, err)
		return
	}

	var body contracts.InvestmentWithdrawRequest
	if errs := c.ShouldBindJSON(&body); errs != nil {
		h.respondError(c, appErrors.ErrBadRequest.WithError(errs))
		return
	}

	ctx := c.Request.Context()
	if err := h.InvestmentService.MakeWithdraw(ctx, investmentID, userID, body.Amount, body.Description); err != nil {
		h.respondError(c, err)
		return
	}

	c.JSON(http.StatusOK, contracts.MessageResponse{Message: "Resgate realizado com sucesso"})
}

func (h *Handler) GetInvestmentReturn(c *gin.Context) {
	investmentID, err := pkg.ParseULID(c.Param("id"))
	if err != nil {
		h.respondError(c, appErrors.NewValidationError("id", "formato inválido"))
		return
	}

	userID, err := h.GetUserIDFromContext(c)
	if err != nil {
		h.respondError(c, err)
		return
	}

	ctx := c.Request.Context()
	profit, returnPercentage, err := h.InvestmentService.CalculateReturn(ctx, investmentID, userID)
	if err != nil {
		h.respondError(c, err)
		return
	}

	c.JSON(http.StatusOK, contracts.InvestmentReturnResponse{
		Profit:           profit,
		ReturnPercentage: returnPercentage,
	})
}

func (h *Handler) DeleteInvestment(c *gin.Context) {
	investmentID, err := pkg.ParseULID(c.Param("id"))
	if err != nil {
		h.respondError(c, appErrors.NewValidationError("id", "formato inválido"))
		return
	}

	userID, err := h.GetUserIDFromContext(c)
	if err != nil {
		h.respondError(c, err)
		return
	}

	ctx := c.Request.Context()
	if err := h.InvestmentService.DeleteInvestment(ctx, investmentID, userID); err != nil {
		h.respondError(c, err)
		return
	}

	c.JSON(http.StatusOK, contracts.MessageResponse{Message: "Investimento excluído com sucesso"})
}

func (h *Handler) UpdateInvestment(c *gin.Context) {
	investmentID, err := pkg.ParseULID(c.Param("id"))
	if err != nil {
		h.respondError(c, appErrors.NewValidationError("id", "formato inválido"))
		return
	}

	userID, err := h.GetUserIDFromContext(c)
	if err != nil {
		h.respondError(c, err)
		return
	}

	var body contracts.InvestmentUpdateRequest
	if err := c.ShouldBindJSON(&body); err != nil {
		h.respondError(c, appErrors.ErrBadRequest.WithError(err))
		return
	}

	updateReq := investment.UpdateInvestmentRequest{
		UserId: userID,
		Id:     investmentID,
	}

	if body.Name != nil {
		updateReq.Name = body.Name
	}
	if body.Type != nil {
		updateReq.Type = body.Type
	}
	if body.ReturnRate != nil {
		updateReq.ReturnRate = body.ReturnRate
	}

	ctx := c.Request.Context()
	if err := h.InvestmentService.UpdateInvestment(ctx, investmentID, userID, updateReq); err != nil {
		h.respondError(c, err)
		return
	}

	c.JSON(http.StatusOK, contracts.MessageResponse{Message: "Investimento atualizado com sucesso"})
}
