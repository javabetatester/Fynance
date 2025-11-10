package routes

import (
	"net/http"

	"Fynance/internal/contracts"
	"Fynance/internal/domain/transaction"
	appErrors "Fynance/internal/errors"
	"Fynance/internal/pkg"

	"github.com/gin-gonic/gin"
)

func (h *Handler) CreateCategory(c *gin.Context) {
	var body contracts.CategoryCreateRequest
	if err := c.ShouldBindJSON(&body); err != nil {
		h.respondError(c, appErrors.ErrBadRequest.WithError(err))
		return
	}

	userID, err := h.GetUserIDFromContext(c)
	if err != nil {
		h.respondError(c, err)
		return
	}

	category := transaction.Category{
		UserId: userID,
		Name:   body.Name,
		Icon:   body.Icon,
	}

	ctx := c.Request.Context()
	if err := h.TransactionService.CreateCategory(ctx, &category); err != nil {
		h.respondError(c, err)
		return
	}

	c.JSON(http.StatusCreated, contracts.MessageResponse{Message: "Categoria criada com sucesso"})
}

func (h *Handler) ListCategories(c *gin.Context) {
	userID, err := h.GetUserIDFromContext(c)
	if err != nil {
		h.respondError(c, err)
		return
	}

	ctx := c.Request.Context()
	categories, err := h.TransactionService.GetAllCategories(ctx, userID)
	if err != nil {
		h.respondError(c, err)
		return
	}

	c.JSON(http.StatusOK, contracts.CategoryListResponse{Categories: categories, Total: len(categories)})
}

func (h *Handler) UpdateCategory(c *gin.Context) {
	categoryID, err := pkg.ParseULID(c.Param("id"))
	if err != nil {
		h.respondError(c, appErrors.NewValidationError("id", "formato inválido"))
		return
	}

	userID, err := h.GetUserIDFromContext(c)
	if err != nil {
		h.respondError(c, err)
		return
	}

	var body contracts.CategoryUpdateRequest
	if err := c.ShouldBindJSON(&body); err != nil {
		h.respondError(c, appErrors.ErrBadRequest.WithError(err))
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
		h.respondError(c, err)
		return
	}

	c.JSON(http.StatusOK, contracts.MessageResponse{Message: "Categoria atualizada com sucesso"})
}

func (h *Handler) DeleteCategory(c *gin.Context) {
	categoryID, err := pkg.ParseULID(c.Param("id"))
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
	if err := h.TransactionService.DeleteCategory(ctx, categoryID, userID); err != nil {
		h.respondError(c, err)
		return
	}

	c.JSON(http.StatusOK, contracts.MessageResponse{Message: "Categoria removida com sucesso"})
}
