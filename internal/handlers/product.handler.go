package handlers

import (
	"errors"

	"github.com/gin-gonic/gin"
	"github.com/oojoseph67/ecommerce/internal/dto"
	"github.com/oojoseph67/ecommerce/internal/utils/responses"
	"github.com/oojoseph67/ecommerce/internal/utils/validators"
)

// CreateProduct godoc
// @Summary Create a product
// @Description Create a new product (admin only)
// @Tags Products
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body dto.CreateProductRequest true "Product data"
// @Success 201 {object} responses.Response{data=dto.ProductResponse}
// @Failure 400 {object} responses.Response
// @Failure 401 {object} responses.Response
// @Failure 403 {object} responses.Response
// @Router /products/ [post]
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

// GetProducts godoc
// @Summary Get all products
// @Description Retrieve a paginated list of products with optional category filter
// @Tags Products
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Items per page" default(10)
// @Param category query string false "Category ID filter"
// @Success 200 {object} responses.PaginatedResponse{response=responses.Response{data=[]dto.ProductResponse},meta=responses.PaginationMeta}
// @Failure 401 {object} responses.Response
// @Failure 500 {object} responses.Response
// @Router /products/ [get]
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

// GetProduct godoc
// @Summary Get a product
// @Description Retrieve a single product by ID
// @Tags Products
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Product ID"
// @Success 200 {object} responses.Response{data=dto.ProductResponse}
// @Failure 401 {object} responses.Response
// @Failure 404 {object} responses.Response
// @Router /products/{id} [get]
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

// UpdateProduct godoc
// @Summary Update a product
// @Description Update an existing product (admin only)
// @Tags Products
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Product ID"
// @Param request body dto.UpdateProductRequest true "Product update data"
// @Success 200 {object} responses.Response{data=dto.ProductResponse}
// @Failure 400 {object} responses.Response
// @Failure 401 {object} responses.Response
// @Failure 403 {object} responses.Response
// @Router /products/{id} [put]
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

// DeleteProduct godoc
// @Summary Delete a product
// @Description Delete a product (admin only)
// @Tags Products
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Product ID"
// @Success 200 {object} responses.Response
// @Failure 401 {object} responses.Response
// @Failure 403 {object} responses.Response
// @Failure 500 {object} responses.Response
// @Router /products/{id} [delete]
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

// UploadProductImage godoc
// @Summary Upload product images
// @Description Upload one or more images for a product (admin only)
// @Tags Products
// @Accept multipart/form-data
// @Produce json
// @Security BearerAuth
// @Param id path string true "Product ID"
// @Param productImage formData file true "Image files to upload"
// @Success 200 {object} responses.Response
// @Failure 400 {object} responses.Response
// @Failure 401 {object} responses.Response
// @Failure 403 {object} responses.Response
// @Router /products/{id}/upload [post]
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

// DeleteProductImage godoc
// @Summary Delete product image
// @Description Delete a product image (admin only)
// @Tags Products
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Product ID"
// @Success 200 {object} responses.Response
// @Failure 401 {object} responses.Response
// @Failure 403 {object} responses.Response
// @Failure 500 {object} responses.Response
// @Router /products/{id}/upload [delete]
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
