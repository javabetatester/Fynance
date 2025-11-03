package routes

import (
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
// @Param        login body auth.Login true "Credenciais de login"
// @Success      200 {object} map[string]interface{} "Login bem-sucedido"
// @Failure      400 {object} map[string]string "Erro de validação"
// @Failure      401 {object} map[string]string "Credenciais inválidas"
// @Router       /api/auth/login [post]
func (h *Handler) Authenticate(c *gin.Context) {
	var login auth.Login
	if err := c.ShouldBindJSON(&login); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := h.AuthService.Login(login)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	userID, err := utils.ParseULID(user.Id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	token, err := h.JwtService.GenerateToken(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Login successful",
		"user":    user.Name,
		"token":   token,
	})
}

// Registration godoc
// @Summary      Registro de novo usuário
// @Description  Cria um novo usuário no sistema
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        user body object true "Dados do usuário"
// @Success      201 {string} string "Usuário registrado com sucesso"
// @Failure      400 {object} map[string]string "Erro de validação"
// @Failure      500 {object} map[string]string "Erro interno do servidor"
// @Router       /api/auth/register [post]
func (h *Handler) Registration(c *gin.Context) {
	var user user.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.AuthService.Register(&user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, "User registered with success")
}
