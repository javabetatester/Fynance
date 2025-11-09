package routes

import (
	"errors"
	"net/http"

	"Fynance/internal/contracts"
	"Fynance/internal/domain/transaction"
	"Fynance/internal/pkg"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func (h *Handler) CreateCategory(c *gin.Context) {
	var body contracts.CategoryCreateRequest
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, contracts.ErrorResponse{Error: err.Error()})
		return
	}

	userID, err := h.GetUserIDFromContext(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, contracts.ErrorResponse{Error: err.Error()})
		return
	}

	category := transaction.Category{
		UserId: userID,
		Name:   body.Name,
		Icon:   body.Icon,
	}

	ctx := c.Request.Context()
	if err := h.TransactionService.CreateCategory(ctx, &category); err != nil {
		c.JSON(http.StatusInternalServerError, contracts.ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusCreated, contracts.MessageResponse{Message: "Categoria criada com sucesso"})
}

func (h *Handler) ListCategories(c *gin.Context) {
	userID, err := h.GetUserIDFromContext(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, contracts.ErrorResponse{Error: err.Error()})
		return
	}

	ctx := c.Request.Context()
	categories, err := h.TransactionService.GetAllCategories(ctx, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, contracts.ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, contracts.CategoryListResponse{Categories: categories, Total: len(categories)})
}

func (h *Handler) UpdateCategory(c *gin.Context) {
	categoryID, err := pkg.ParseULID(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, contracts.ErrorResponse{Error: "id de categoria inválido"})
		return
	}

	userID, err := h.GetUserIDFromContext(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, contracts.ErrorResponse{Error: err.Error()})
		return
	}

	var body contracts.CategoryUpdateRequest
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, contracts.ErrorResponse{Error: err.Error()})
		return
	}

	category := transaction.Category{
		Id:     categoryID,
		UserId: userID,
		Name:   body.Name,
		Icon:   body.Icon,
	}

	ctx := c.Request.Context()
	if err := h.TransactionService.UpdateCategory(ctx, &category); err != nil {
		if err.Error() == "category not found" || errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, contracts.ErrorResponse{Error: "Categoria não encontrada"})
			return
		}
		if err.Error() == "category already exists" {
			c.JSON(http.StatusConflict, contracts.ErrorResponse{Error: "Categoria já existente"})
			return
		}
		c.JSON(http.StatusInternalServerError, contracts.ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, contracts.MessageResponse{Message: "Categoria atualizada com sucesso"})
}

func (h *Handler) DeleteCategory(c *gin.Context) {
	categoryID, err := pkg.ParseULID(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, contracts.ErrorResponse{Error: "id de categoria inválido"})
		return
	}

	userID, err := h.GetUserIDFromContext(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, contracts.ErrorResponse{Error: err.Error()})
		return
	}

	ctx := c.Request.Context()
	if err := h.TransactionService.DeleteCategory(ctx, categoryID, userID); err != nil {
		if err.Error() == "category not found" || errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, contracts.ErrorResponse{Error: "Categoria não encontrada"})
			return
		}
		c.JSON(http.StatusInternalServerError, contracts.ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, contracts.MessageResponse{Message: "Categoria removida com sucesso"})
}
