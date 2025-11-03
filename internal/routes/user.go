package routes

import (
	"Fynance/internal/domain/user"
	"Fynance/internal/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

// CreateUser godoc
// @Summary      Criar usuário
// @Description  Cria um novo usuário no sistema
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        user body object true "Dados do usuário"
// @Success      201 {string} string "Usuário criado com sucesso"
// @Failure      400 {object} map[string]string "Erro de validação"
// @Failure      500 {object} map[string]string "Erro interno do servidor"
// @Router       /api/users [post]
func (h *Handler) CreateUser(c *gin.Context) {
	var user user.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := h.UserService.Create(&user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, "User created with success")
}

// GetUserByID godoc
// @Summary      Obter usuário por ID
// @Description  Busca um usuário pelo seu ID
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        id path string true "ID do usuário"
// @Success      200 {object} object "Dados do usuário"
// @Failure      500 {object} map[string]string "Erro interno do servidor"
// @Router       /api/users/{id} [get]
// @Security     BearerAuth
func (h *Handler) GetUserByID(c *gin.Context) {
	id := c.Param("id")

	user, err := h.UserService.GetByID(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, user)
}

// GetUserByEmail godoc
// @Summary      Obter usuário por email
// @Description  Busca um usuário pelo seu email
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        email query string true "Email do usuário"
// @Success      200 {object} object "Dados do usuário"
// @Failure      500 {object} map[string]string "Erro interno do servidor"
// @Router       /api/users/email [get]
// @Security     BearerAuth
func (h *Handler) GetUserByEmail(c *gin.Context) {
	email := c.Param("email")
	user, err := h.UserService.GetByEmail(email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, user)
}

// UpdateUser godoc
// @Summary      Atualizar usuário
// @Description  Atualiza os dados de um usuário (apenas o próprio usuário pode atualizar)
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        id path string true "ID do usuário"
// @Param        user body user.User true "Dados atualizados do usuário"
// @Success      200 {object} object "Dados do usuário atualizados"
// @Failure      400 {object} map[string]string "Erro de validação"
// @Failure      403 {object} map[string]string "Acesso negado"
// @Failure      500 {object} map[string]string "Erro interno do servidor"
// @Router       /api/users/{id} [put]
// @Security     BearerAuth
func (h *Handler) UpdateUser(c *gin.Context) {
	idStr := c.Param("id")
	id, err := utils.ParseULID(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userIDFromToken, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	if userIDFromToken.(string) != id.String() {
		c.JSON(http.StatusForbidden, gin.H{"error": "You can only update your own profile"})
		return
	}

	var user user.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if user.Id != "" && user.Id != id.String() {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID mismatch"})
		return
	}

	user.Id = id.String()

	if err := h.UserService.Update(&user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, user)
}

// DeleteUser godoc
// @Summary      Excluir usuário
// @Description  Exclui um usuário do sistema (apenas o próprio usuário pode se excluir)
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        id path string true "ID do usuário"
// @Success      200 {object} map[string]string "Usuário excluído com sucesso"
// @Failure      400 {object} map[string]string "ID inválido"
// @Failure      403 {object} map[string]string "Acesso negado"
// @Failure      500 {object} map[string]string "Erro interno do servidor"
// @Router       /api/users/{id} [delete]
// @Security     BearerAuth
func (h *Handler) DeleteUser(c *gin.Context) {
	idStr := c.Param("id")
	id, err := utils.ParseULID(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userIDFromToken, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	if userIDFromToken.(string) != id.String() {
		c.JSON(http.StatusForbidden, gin.H{"error": "You can only delete your own profile"})
		return
	}

	if err := h.UserService.Delete(id.String()); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User deleted successfully"})
}
