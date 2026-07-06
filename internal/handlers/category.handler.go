package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/oojoseph67/ecommerce/internal/dto"
	"github.com/oojoseph67/ecommerce/internal/utils/responses"
	"github.com/oojoseph67/ecommerce/internal/utils/validators"
)

// CreateCategory godoc
// @Summary Create a category
// @Description Create a new product category (admin only)
// @Tags Categories
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body dto.CreateCategoryRequest true "Category data"
// @Success 201 {object} responses.Response{data=dto.CategoryResponse}
// @Failure 400 {object} responses.Response
// @Failure 401 {object} responses.Response
// @Failure 403 {object} responses.Response
// @Router /categories/ [post]
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

// GetCategories godoc
// @Summary Get all categories
// @Description Retrieve all product categories
// @Tags Categories
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} responses.Response{data=dto.GetCategoriesResponse}
// @Failure 401 {object} responses.Response
// @Failure 500 {object} responses.Response
// @Router /categories/ [get]
func (h *Handler) GetCategories(ctx *gin.Context) {
	categories, err := h.categoryService.GetCategories()
	if err != nil {
		responses.InternalServerResponse(ctx, "Couldnt get categories", err)
		return
	}

	responses.SuccessResponse(ctx, "Category retrieved", categories)
}

// UpdateCategory godoc
// @Summary Update a category
// @Description Update an existing product category (admin only)
// @Tags Categories
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Category ID"
// @Param request body dto.UpdateCategoryRequest true "Category update data"
// @Success 200 {object} responses.Response{data=dto.CategoryResponse}
// @Failure 400 {object} responses.Response
// @Failure 401 {object} responses.Response
// @Failure 403 {object} responses.Response
// @Router /categories/{id} [put]
func (h *Handler) UpdateCategory(ctx *gin.Context) {
	id, ok := getReqParam(ctx, "id")
	if !ok {
		return
	}

	var req *dto.UpdateCategoryRequest
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

// DeleteCategory godoc
// @Summary Delete a category
// @Description Delete a product category (admin only)
// @Tags Categories
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Category ID"
// @Success 200 {object} responses.Response
// @Failure 401 {object} responses.Response
// @Failure 403 {object} responses.Response
// @Failure 500 {object} responses.Response
// @Router /categories/{id} [delete]
func (h *Handler) DeleteCategory(ctx *gin.Context) {
	id, ok := getReqParam(ctx, "id")
	if !ok {
		return
	}

	if err := h.categoryService.DeleteCategory(id); err != nil {
		responses.InternalServerResponse(ctx, "Couldnt delete category", err)
		return
	}

	responses.SuccessResponse(ctx, "Category deleted successful", nil)
}
