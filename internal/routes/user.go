package routes

import (
	"Fynance/internal/domain/user"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

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
	c.JSON(http.StatusCreated, "user created with success")
}

func (h *Handler) GetUserByID(c *gin.Context) {
	id := c.Param("id")

	user, err := h.UserService.GetByID(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, user)
}

func (h *Handler) GetUserByEmail(c *gin.Context) {
	email := c.Param("email")
	user, err := h.UserService.GetByEmail(email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, user)
}

func (h *Handler) UpdateUser(c *gin.Context) {

	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid UUID"})
		return
	}

	var user user.User

	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if user.Id != uuid.Nil && user.Id != id {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID mismatch"})
		return
	}

	user.Id = id

	if err := h.UserService.Update(&user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, user)
}

func (h *Handler) DeleteUser(c *gin.Context) {
	id := c.Param("id")

	if err := h.UserService.Delete(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User deleted"})
}
