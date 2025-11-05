package routes

import (
	"errors"
	"net/http"

	"Fynance/internal/contracts"
	"Fynance/internal/domain/investment"
	"Fynance/internal/utils"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func (h *Handler) CreateInvestment(c *gin.Context) {
	var body contracts.InvestmentCreateRequest
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, contracts.ErrorResponse{Error: err.Error()})
		return
	}

	userID, err := h.GetUserIDFromContext(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, contracts.ErrorResponse{Error: err.Error()})
		return
	}

	req := investment.CreateInvestmentRequest{
		UserId:        userID,
		Type:          investment.Types(body.Type),
		Name:          body.Name,
		InitialAmount: body.InitialAmount,
		ReturnRate:    body.ReturnRate,
	}

	inv, err := h.InvestmentService.CreateInvestment(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, contracts.ErrorResponse{Error: err.Error()})
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
		c.JSON(http.StatusUnauthorized, contracts.ErrorResponse{Error: err.Error()})
		return
	}

	investments, err := h.InvestmentService.ListInvestments(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, contracts.ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, contracts.InvestmentListResponse{
		Total:       len(investments),
		Investments: investments,
	})
}

func (h *Handler) GetInvestment(c *gin.Context) {
	investmentID, err := utils.ParseULID(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, contracts.ErrorResponse{Error: "id de investimento inválido"})
		return
	}

	userID, err := h.GetUserIDFromContext(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, contracts.ErrorResponse{Error: err.Error()})
		return
	}

	inv, err := h.InvestmentService.GetInvestment(investmentID, userID)
	if err != nil {
		c.JSON(http.StatusNotFound, contracts.ErrorResponse{Error: "investimento não encontrado"})
		return
	}

	c.JSON(http.StatusOK, contracts.InvestmentSingleResponse{Investment: inv})
}

func (h *Handler) MakeContribution(c *gin.Context) {
	investmentID, err := utils.ParseULID(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, contracts.ErrorResponse{Error: "id de investimento inválido"})
		return
	}

	userID, err := h.GetUserIDFromContext(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, contracts.ErrorResponse{Error: err.Error()})
		return
	}

	var body contracts.InvestmentContributionRequest
	if errs := c.ShouldBindJSON(&body); errs != nil {
		c.JSON(http.StatusBadRequest, contracts.ErrorResponse{Error: errs.Error()})
		return
	}

	if err := h.InvestmentService.MakeContribution(investmentID, userID, body.Amount, body.Description); err != nil {
		switch err.Error() {
		case "investment not found":
			c.JSON(http.StatusNotFound, contracts.ErrorResponse{Error: "investimento não encontrado"})
			return
		case "amount must be greater than 0":
			c.JSON(http.StatusBadRequest, contracts.ErrorResponse{Error: "valor do aporte deve ser maior que zero"})
			return
		}
		c.JSON(http.StatusInternalServerError, contracts.ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, contracts.MessageResponse{Message: "Aporte registrado com sucesso"})
}

func (h *Handler) MakeWithdraw(c *gin.Context) {
	investmentID, err := utils.ParseULID(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, contracts.ErrorResponse{Error: "id de investimento inválido"})
		return
	}

	userID, err := h.GetUserIDFromContext(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, contracts.ErrorResponse{Error: err.Error()})
		return
	}

	var body contracts.InvestmentWithdrawRequest
	if errs := c.ShouldBindJSON(&body); errs != nil {
		c.JSON(http.StatusBadRequest, contracts.ErrorResponse{Error: errs.Error()})
		return
	}

	if err := h.InvestmentService.MakeWithdraw(investmentID, userID, body.Amount, body.Description); err != nil {
		switch err.Error() {
		case "investment not found":
			c.JSON(http.StatusNotFound, contracts.ErrorResponse{Error: "investimento não encontrado"})
			return
		case "insufficient balance in investment":
			c.JSON(http.StatusBadRequest, contracts.ErrorResponse{Error: "saldo insuficiente no investimento"})
			return
		case "amount must be greater than 0":
			c.JSON(http.StatusBadRequest, contracts.ErrorResponse{Error: "valor deve ser maior que zero"})
			return
		}
		c.JSON(http.StatusInternalServerError, contracts.ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, contracts.MessageResponse{Message: "Resgate realizado com sucesso"})
}

func (h *Handler) GetInvestmentReturn(c *gin.Context) {
	investmentID, err := utils.ParseULID(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, contracts.ErrorResponse{Error: "id de investimento inválido"})
		return
	}

	userID, err := h.GetUserIDFromContext(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, contracts.ErrorResponse{Error: err.Error()})
		return
	}

	profit, returnPercentage, err := h.InvestmentService.CalculateReturn(investmentID, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, contracts.ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, contracts.InvestmentReturnResponse{
		Profit:           profit,
		ReturnPercentage: returnPercentage,
	})
}

func (h *Handler) DeleteInvestment(c *gin.Context) {
	investmentID, err := utils.ParseULID(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, contracts.ErrorResponse{Error: "id de investimento inválido"})
		return
	}

	userID, err := h.GetUserIDFromContext(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, contracts.ErrorResponse{Error: err.Error()})
		return
	}

	if err := h.InvestmentService.DeleteInvestment(investmentID, userID); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) || err.Error() == "investment not found" {
			c.JSON(http.StatusNotFound, contracts.ErrorResponse{Error: "investimento não encontrado"})
			return
		}
		if err.Error() == "cannot delete investment with balance" {
			c.JSON(http.StatusBadRequest, contracts.ErrorResponse{Error: "não é possível excluir investimento com saldo"})
			return
		}
		c.JSON(http.StatusInternalServerError, contracts.ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, contracts.MessageResponse{Message: "Investimento excluído com sucesso"})
}

func (h *Handler) UpdateInvestment(c *gin.Context) {
	investmentID, err := utils.ParseULID(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, contracts.ErrorResponse{Error: "id de investimento inválido"})
		return
	}

	userID, err := h.GetUserIDFromContext(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, contracts.ErrorResponse{Error: err.Error()})
		return
	}

	var body contracts.InvestmentUpdateRequest
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, contracts.ErrorResponse{Error: err.Error()})
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

	if err := h.InvestmentService.UpdateInvestment(investmentID, userID, updateReq); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) || err.Error() == "investment not found" {
			c.JSON(http.StatusNotFound, contracts.ErrorResponse{Error: "investimento não encontrado"})
			return
		}
		if err.Error() == "name is required" {
			c.JSON(http.StatusBadRequest, contracts.ErrorResponse{Error: "nome é obrigatório"})
			return
		}
		c.JSON(http.StatusInternalServerError, contracts.ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, contracts.MessageResponse{Message: "Investimento atualizado com sucesso"})
}
