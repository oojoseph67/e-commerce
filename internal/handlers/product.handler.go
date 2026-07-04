package handlers

import (
	"errors"

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

	product, err := h.productService.CreateProduct(req)
	if err != nil {
		responses.InternalServerResponse(ctx, "Couldnt create product", err)
		return
	}

	responses.CreatedResponse(ctx, "Product created successful", product)
}

func (h *Handler) GetProducts(ctx *gin.Context) {
	page, limit := getPaginationValues(ctx)

	category := ctx.Query("category")

	products, pagination, err := h.productService.GetProducts(page, limit, category)
	if err != nil {
		responses.InternalServerResponse(ctx, "Failed to retrieve products", err)
		return
	}

	responses.PaginatedSuccessResponse(ctx, "Products retrieved successful", products, *pagination)
}

func (h *Handler) GetProduct(ctx *gin.Context) {
	id, ok := getReqParam(ctx, "id")
	if !ok {
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
	id, ok := getReqParam(ctx, "id")
	if !ok {
		return
	}

	var req *dto.UpdateProductRequest
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
	id, ok := getReqParam(ctx, "id")
	if !ok {
		return
	}

	if err := h.productService.DeleteProduct(id); err != nil {
		responses.InternalServerResponse(ctx, "Couldnt delete product", err)
		return
	}

	responses.SuccessResponse(ctx, "Product deleted successful", nil)
}

func (h *Handler) UploadProductImage(ctx *gin.Context) {
	id, ok := getReqParam(ctx, "id")
	if !ok {
		return
	}

	// file, err := ctx.FormFile("productImage")
	// if err != nil {
	// 	responses.BadRequestResponse(ctx, "No image file uploaded", err)
	// 	return
	// }

	form, err := ctx.MultipartForm()
	if err != nil {
		responses.BadRequestResponse(ctx, "No image files uploaded", err)
		return
	}

	files := form.File["productImage"]
	if len(files) == 0 {
		responses.BadRequestResponse(ctx, "No image files uploaded", errors.New("no files provided"))
		return
	}

	// url, altText, err := h.uploadService.UploadProductImage(id, file)
	// if err != nil {
	// 	responses.InternalServerResponse(ctx, "Couldnt upload image", err)
	// 	return
	// }

	var results []dto.ProductImageUpload
	for _, file := range files {
		url, altText, err := h.uploadService.UploadProductImage(id, file)
		if err != nil {
			responses.InternalServerResponse(ctx, "Couldn't upload image", err)
			return
		}
		results = append(results, dto.ProductImageUpload{
			URL:     url,
			AltText: altText,
		})
	}

	if err := h.productService.AddProductImages(id, results); err != nil {
		responses.InternalServerResponse(ctx, "Couldn't add images to product", err)
		return
	}

	// err = h.productService.AddProductImage(id, url, altText)
	// if err != nil {
	// 	responses.InternalServerResponse(ctx, "Couldnt add image to product", err)
	// 	return
	// }

	responses.SuccessResponse(ctx, "Product images added successfully", nil)
}

func (h *Handler) DeleteProductImage(ctx *gin.Context) {
	id, ok := getReqParam(ctx, "id")
	if !ok {
		return
	}

	if err := h.productService.DeleteProductImage(id); err != nil {
		responses.InternalServerResponse(ctx, "Couldnt delete product image", err)
		return
	}

	responses.SuccessResponse(ctx, "Product image deleted successful", nil)
}
