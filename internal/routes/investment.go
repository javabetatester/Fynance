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

// CreateInvestment godoc
// @Summary      Criar investimento
// @Description  Cria um novo investimento para o usuário autenticado
// @Tags         investments
// @Accept       json
// @Produce      json
// @Param        investment body contracts.InvestmentCreateRequest true "Dados do investimento"
// @Success      201 {object} contracts.InvestmentCreateResponse "Investimento criado com sucesso"
// @Failure      400 {object} contracts.ErrorResponse "Erro de validação"
// @Failure      401 {object} contracts.ErrorResponse "Não autorizado"
// @Failure      500 {object} contracts.ErrorResponse "Erro interno do servidor"
// @Router       /api/investments [post]
// @Security     BearerAuth
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

// ListInvestments godoc
// @Summary      Listar investimentos
// @Description  Lista todos os investimentos do usuário autenticado
// @Tags         investments
// @Accept       json
// @Produce      json
// @Success      200 {object} contracts.InvestmentListResponse "Lista de investimentos"
// @Failure      401 {object} contracts.ErrorResponse "Não autorizado"
// @Failure      500 {object} contracts.ErrorResponse "Erro interno do servidor"
// @Router       /api/investments [get]
// @Security     BearerAuth
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

// GetInvestment godoc
// @Summary      Obter investimento
// @Description  Recupera os dados de um investimento específico do usuário autenticado
// @Tags         investments
// @Accept       json
// @Produce      json
// @Param        id path string true "ID do investimento"
// @Success      200 {object} contracts.InvestmentSingleResponse "Dados do investimento"
// @Failure      400 {object} contracts.ErrorResponse "ID inválido"
// @Failure      401 {object} contracts.ErrorResponse "Não autorizado"
// @Failure      404 {object} contracts.ErrorResponse "Investimento não encontrado"
// @Router       /api/investments/{id} [get]
// @Security     BearerAuth
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

// MakeContribution godoc
// @Summary      Realizar aporte
// @Description  Registra um aporte financeiro em um investimento existente
// @Tags         investments
// @Accept       json
// @Produce      json
// @Param        id path string true "ID do investimento"
// @Param        contribution body contracts.InvestmentContributionRequest true "Dados do aporte"
// @Success      200 {object} contracts.MessageResponse "Aporte registrado"
// @Failure      400 {object} contracts.ErrorResponse "Dados inválidos"
// @Failure      401 {object} contracts.ErrorResponse "Não autorizado"
// @Failure      404 {object} contracts.ErrorResponse "Investimento não encontrado"
// @Failure      500 {object} contracts.ErrorResponse "Erro interno do servidor"
// @Router       /api/investments/{id}/contribution [post]
// @Security     BearerAuth
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
		if err.Error() == "investment not found" {
			c.JSON(http.StatusNotFound, contracts.ErrorResponse{Error: "investimento não encontrado"})
			return
		}
		c.JSON(http.StatusInternalServerError, contracts.ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, contracts.MessageResponse{Message: "Aporte registrado com sucesso"})
}

// MakeWithdraw godoc
// @Summary      Realizar resgate
// @Description  Registra um resgate financeiro de um investimento existente
// @Tags         investments
// @Accept       json
// @Produce      json
// @Param        id path string true "ID do investimento"
// @Param        withdraw body contracts.InvestmentWithdrawRequest true "Dados do resgate"
// @Success      200 {object} contracts.MessageResponse "Resgate realizado"
// @Failure      400 {object} contracts.ErrorResponse "Dados inválidos"
// @Failure      401 {object} contracts.ErrorResponse "Não autorizado"
// @Failure      404 {object} contracts.ErrorResponse "Investimento não encontrado"
// @Failure      500 {object} contracts.ErrorResponse "Erro interno do servidor"
// @Router       /api/investments/{id}/withdraw [post]
// @Security     BearerAuth
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
		if err.Error() == "investment not found" {
			c.JSON(http.StatusNotFound, contracts.ErrorResponse{Error: "investimento não encontrado"})
			return
		}
		c.JSON(http.StatusInternalServerError, contracts.ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, contracts.MessageResponse{Message: "Resgate realizado com sucesso"})
}

// GetInvestmentReturn godoc
// @Summary      Calcular retorno de investimento
// @Description  Obtém lucro absoluto e percentual acumulado do investimento
// @Tags         investments
// @Accept       json
// @Produce      json
// @Param        id path string true "ID do investimento"
// @Success      200 {object} contracts.InvestmentReturnResponse "Retorno calculado"
// @Failure      400 {object} contracts.ErrorResponse "ID inválido"
// @Failure      401 {object} contracts.ErrorResponse "Não autorizado"
// @Failure      500 {object} contracts.ErrorResponse "Erro interno do servidor"
// @Router       /api/investments/{id}/return [get]
// @Security     BearerAuth
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

// DeleteInvestment godoc
// @Summary      Excluir investimento
// @Description  Remove um investimento do usuário quando não há saldo pendente
// @Tags         investments
// @Accept       json
// @Produce      json
// @Param        id path string true "ID do investimento"
// @Success      200 {object} contracts.MessageResponse "Investimento excluído"
// @Failure      400 {object} contracts.ErrorResponse "Requisição inválida"
// @Failure      401 {object} contracts.ErrorResponse "Não autorizado"
// @Failure      404 {object} contracts.ErrorResponse "Investimento não encontrado"
// @Failure      500 {object} contracts.ErrorResponse "Erro interno do servidor"
// @Router       /api/investments/{id} [delete]
// @Security     BearerAuth
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

	var body contracts.InvestimentUpdateRequest
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, contracts.ErrorResponse{Error: err.Error()})
		return
	}

	updateReq := investment.UpdateInvestmentRequest{
		UserId:        userID,
		Id:            investmentID,
		Name:          body.Name,
		Type:          body.Type,
		InitialAmount: body.InitialAmount,
		ReturnRate:    body.ReturnRate,
	}

	if err := h.InvestmentService.UpdateInvestment(investmentID, userID, updateReq); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) || err.Error() == "investment not found" {
			c.JSON(http.StatusNotFound, contracts.ErrorResponse{Error: "investimento não encontrado"})
			return
		}
		c.JSON(http.StatusInternalServerError, contracts.ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, contracts.MessageResponse{Message: "Investimento atualizado com sucesso"})
}
