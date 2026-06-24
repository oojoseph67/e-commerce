package handlers

import (
	"errors"
	_ "fmt"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/oojoseph67/ecommerce/internal/dto"
	"github.com/oojoseph67/ecommerce/internal/utils/responses"
	"github.com/oojoseph67/ecommerce/internal/utils/validators"
)

/**** CATEGORY HANDLER */

func (h *Handler) CreateCategory(ctx *gin.Context) {
	var req *dto.CreateCategoryRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		responses.BadRequestResponse(ctx, "Invalid request data", validators.FormatValidationError(err))
		return
	}

	category, err := h.service.CreateCategory(req)
	if err != nil {
		responses.InternalServerResponse(ctx, "Couldnt create category", err)
		return
	}

	responses.CreatedResponse(ctx, "Category created successful", category)
}

func (h *Handler) GetCategories(ctx *gin.Context) {

	categories, err := h.service.GetCategories()
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

	category, err := h.service.UpdateCategory(id, req)
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

	if err := h.service.DeleteCategory(id); err != nil {
		responses.InternalServerResponse(ctx, "Couldnt delete category", err)
		return
	}

	responses.SuccessResponse(ctx, "Category deleted successful", nil)
}

/**** PRODUCT HANDLER */

func (h *Handler) CreateProduct(ctx *gin.Context) {
	var req *dto.CreateProductRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		responses.BadRequestResponse(ctx, "Invalid request data", validators.FormatValidationError(err))
		return
	}

	category, err := h.service.CreateProduct(req)
	if err != nil {
		responses.InternalServerResponse(ctx, "Couldnt create product", err)
		return
	}

	responses.CreatedResponse(ctx, "Product created successful", category)
}

func (h *Handler) GetProducts(ctx *gin.Context) {
	// page := ctx.Query("page")
	// limit := ctx.Query("limit")

	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(ctx.DefaultQuery("limit", "10"))

	products, pagination, err := h.service.GetProducts(page, limit)
	if err != nil {
		responses.InternalServerResponse(ctx, "Failed to retrieve products", err)
		return
	}

	responses.PaginatedSuccessResponse(ctx, "Products retrieved successful", products, *pagination)
}

func (h *Handler) GetProduct(ctx *gin.Context) {
	id := ctx.Param("id")

	if id == "" {
		responses.BadRequestResponse(ctx, "Please provide id param", errors.New("please provide id param"))
	}

	product, err := h.service.GetProduct(id)
	if err != nil {
		responses.NotFoundResponse(ctx, "Couldnt retrieve product", err)
		return
	}

	responses.SuccessResponse(ctx, "Proudct retrieved successful", product)

}

func (h *Handler) UpdateProduct(ctx *gin.Context) {
	var req *dto.UpdateProductRequest

	id := ctx.Param("id")

	if id == "" {
		responses.BadRequestResponse(ctx, "Please provide id param", errors.New("please provide id param"))
	}

	if err := ctx.ShouldBindJSON(&req); err != nil {
		responses.BadRequestResponse(ctx, "Invalid request data", validators.FormatValidationError(err))
		return
	}

	category, err := h.service.UpdateProduct(id, req)
	if err != nil {
		responses.InternalServerResponse(ctx, "Couldnt update product", err)
		return
	}

	responses.SuccessResponse(ctx, "Product updated successful", category)
}

func (h *Handler) DeleteProduct(ctx *gin.Context) {
	id := ctx.Param("id")

	if id == "" {
		responses.BadRequestResponse(ctx, "Please provide id param", errors.New("please provide id param"))
	}

	if err := h.service.DeleteProduct(id); err != nil {
		responses.InternalServerResponse(ctx, "Couldnt delete product", err)
		return
	}

	responses.SuccessResponse(ctx, "Product deleted successful", nil)
}
