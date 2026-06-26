package handlers

import (
	"errors"
	"fmt"
	_ "fmt"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/oojoseph67/ecommerce/internal/dto"
	"github.com/oojoseph67/ecommerce/internal/utils/responses"
	"github.com/oojoseph67/ecommerce/internal/utils/validators"
)

func (h *Handler) CreateProduct(ctx *gin.Context) {
	var req *dto.CreateProductRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		responses.BadRequestResponse(ctx, "Invalid request data", validators.FormatValidationError(err))
		return
	}

	category, err := h.productService.CreateProduct(req)
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

	products, pagination, err := h.productService.GetProducts(page, limit)
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
		return
	}

	product, err := h.productService.GetProduct(id)
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
		return
	}

	if err := ctx.ShouldBindJSON(&req); err != nil {
		responses.BadRequestResponse(ctx, "Invalid request data", validators.FormatValidationError(err))
		return
	}

	category, err := h.productService.UpdateProduct(id, req)
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

	if err := h.productService.DeleteProduct(id); err != nil {
		responses.InternalServerResponse(ctx, "Couldnt delete product", err)
		return
	}

	responses.SuccessResponse(ctx, "Product deleted successful", nil)
}

func (h *Handler) UploadProductImage(ctx *gin.Context) {
	id := ctx.Param("id")
	if id == "" || id == ":id" {
		responses.BadRequestResponse(ctx, "Please provide id param", errors.New("please provide id param"))
		return
	}

	fmt.Println("id:", id)
	fmt.Println("Content-Type header:", ctx.GetHeader("Content-Type"))

	file, err := ctx.FormFile("productImage")
	if err != nil {
		responses.BadRequestResponse(ctx, "No image file uploaded", err)
		return
	}

	url, altText, err := h.uploadService.UploadProductImage(id, file)
	if err != nil {
		responses.InternalServerResponse(ctx, "Couldnt upload image", err)
		return
	}

	fmt.Println(url, altText)

	err = h.productService.AddProductImage(id, url, altText)
	if err != nil {
		responses.InternalServerResponse(ctx, "Couldnt add image to product", err)
		return
	}

	responses.SuccessResponse(ctx, "Product image added successfully", nil)
}
