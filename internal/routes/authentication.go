package routes

import (
	"Fynance/internal/contracts"
	"Fynance/internal/domain/auth"
	"Fynance/internal/domain/user"
	"Fynance/internal/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

// Authenticate godoc
// @Summary      Login de usuário
// @Description  Autentica um usuário e retorna um token JWT
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        login body contracts.AuthLoginRequest true "Credenciais de login"
// @Success      200 {object} contracts.AuthLoginResponse "Login bem-sucedido"
// @Failure      400 {object} contracts.ErrorResponse "Erro de validação"
// @Failure      401 {object} contracts.ErrorResponse "Credenciais inválidas"
// @Router       /api/auth/login [post]
func (h *Handler) Authenticate(c *gin.Context) {
	var body contracts.AuthLoginRequest
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, contracts.ErrorResponse{Error: err.Error()})
		return
	}

	login := auth.Login{
		Email:    body.Email,
		Password: body.Password,
	}

	userEntity, err := h.AuthService.Login(login)
	if err != nil {
		c.JSON(http.StatusUnauthorized, contracts.ErrorResponse{Error: err.Error()})
		return
	}

	userID, err := utils.ParseULID(userEntity.Id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, contracts.ErrorResponse{Error: err.Error()})
		return
	}

	token, err := h.JwtService.GenerateToken(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, contracts.ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, contracts.AuthLoginResponse{
		Message: "Login realizado com sucesso",
		User:    userEntity.Name,
		Token:   token,
	})
}

// Registration godoc
// @Summary      Registro de novo usuário
// @Description  Cria um novo usuário no sistema
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        user body contracts.AuthRegisterRequest true "Dados do usuário"
// @Success      201 {object} contracts.MessageResponse "Usuário registrado com sucesso"
// @Failure      400 {object} contracts.ErrorResponse "Erro de validação"
// @Failure      500 {object} contracts.ErrorResponse "Erro interno do servidor"
// @Router       /api/auth/register [post]
func (h *Handler) Registration(c *gin.Context) {
	var body contracts.AuthRegisterRequest
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, contracts.ErrorResponse{Error: err.Error()})
		return
	}

	userEntity := user.User{
		Name:     body.Name,
		Email:    body.Email,
		Password: body.Password,
	}

	if err := h.AuthService.Register(&userEntity); err != nil {
		c.JSON(http.StatusInternalServerError, contracts.ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusCreated, contracts.MessageResponse{Message: "Usuário registrado com sucesso"})
}
