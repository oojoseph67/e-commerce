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

func (h *Handler) CreateProduct(ctx *gin.Context) {
	service := h.service.ProductService
	var req *dto.CreateProductRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		responses.BadRequestResponse(ctx, "Invalid request data", validators.FormatValidationError(err))
		return
	}

	category, err := service.CreateProduct(req)
	if err != nil {
		responses.InternalServerResponse(ctx, "Couldnt create product", err)
		return
	}

	responses.CreatedResponse(ctx, "Product created successful", category)
}

func (h *Handler) GetProducts(ctx *gin.Context) {
	// page := ctx.Query("page")
	// limit := ctx.Query("limit")
	service := h.service.ProductService

	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(ctx.DefaultQuery("limit", "10"))

	products, pagination, err := service.GetProducts(page, limit)
	if err != nil {
		responses.InternalServerResponse(ctx, "Failed to retrieve products", err)
		return
	}

	responses.PaginatedSuccessResponse(ctx, "Products retrieved successful", products, *pagination)
}

func (h *Handler) GetProduct(ctx *gin.Context) {
	service := h.service.ProductService
	id := ctx.Param("id")

	if id == "" {
		responses.BadRequestResponse(ctx, "Please provide id param", errors.New("please provide id param"))
	}

	product, err := service.GetProduct(id)
	if err != nil {
		responses.NotFoundResponse(ctx, "Couldnt retrieve product", err)
		return
	}

	responses.SuccessResponse(ctx, "Proudct retrieved successful", product)

}

func (h *Handler) UpdateProduct(ctx *gin.Context) {
	service := h.service.ProductService
	var req *dto.UpdateProductRequest

	id := ctx.Param("id")

	if id == "" {
		responses.BadRequestResponse(ctx, "Please provide id param", errors.New("please provide id param"))
	}

	if err := ctx.ShouldBindJSON(&req); err != nil {
		responses.BadRequestResponse(ctx, "Invalid request data", validators.FormatValidationError(err))
		return
	}

	category, err := service.UpdateProduct(id, req)
	if err != nil {
		responses.InternalServerResponse(ctx, "Couldnt update product", err)
		return
	}

	responses.SuccessResponse(ctx, "Product updated successful", category)
}

func (h *Handler) DeleteProduct(ctx *gin.Context) {
	service := h.service.ProductService
	id := ctx.Param("id")

	if id == "" {
		responses.BadRequestResponse(ctx, "Please provide id param", errors.New("please provide id param"))
	}

	if err := service.DeleteProduct(id); err != nil {
		responses.InternalServerResponse(ctx, "Couldnt delete product", err)
		return
	}

	responses.SuccessResponse(ctx, "Product deleted successful", nil)
}
