package handlers

import (
	"errors"

	"github.com/gin-gonic/gin"
	"github.com/oojoseph67/ecommerce/internal/dto"
	"github.com/oojoseph67/ecommerce/internal/utils/responses"
	"github.com/oojoseph67/ecommerce/internal/utils/validators"
)

func (h *Handler) CreateCategory(ctx *gin.Context) {
	var req *dto.CreateCategoryRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		responses.BadRequestResponse(ctx, "Invalid request data", validators.FormatValidationError(err))
		return
	}

	category, err := h.categoryService.CreateCategory(req)
	if err != nil {
		responses.InternalServerResponse(ctx, "Couldnt create category", err)
		return
	}

	responses.CreatedResponse(ctx, "Category created successful", category)
}

func (h *Handler) GetCategories(ctx *gin.Context) {
	categories, err := h.categoryService.GetCategories()
	if err != nil {
		responses.InternalServerResponse(ctx, "Couldnt get categories", err)
		return
	}

	responses.SuccessResponse(ctx, "Category retrieved", categories)
}

func (h *Handler) UpdateCategory(ctx *gin.Context) {
	var req *dto.UpdateCategoryRequest

	id := ctx.Param("id")
	// strconv.ParseUint(c.Param("id"), 10, 32)

	if id == "" {
		responses.BadRequestResponse(ctx, "Please provide id param", errors.New("please provide id param"))
	}

	if err := ctx.ShouldBindJSON(&req); err != nil {
		responses.BadRequestResponse(ctx, "Invalid request data", validators.FormatValidationError(err))
		return
	}

	category, err := h.categoryService.UpdateCategory(id, req)
	if err != nil {
		responses.InternalServerResponse(ctx, "Couldnt update category", err)
		return
	}

	responses.SuccessResponse(ctx, "Category updated successful", category)
}

func (h *Handler) DeleteCategory(ctx *gin.Context) {
	id := ctx.Param("id")

	if id == "" {
		responses.BadRequestResponse(ctx, "Please provide id param", errors.New("please provide id param"))
	}

	if err := h.categoryService.DeleteCategory(id); err != nil {
		responses.InternalServerResponse(ctx, "Couldnt delete category", err)
		return
	}

	responses.SuccessResponse(ctx, "Category deleted successful", nil)
}
