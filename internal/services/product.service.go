package services

import (
	"errors"

	"github.com/oojoseph67/ecommerce/internal/dto"
	"github.com/oojoseph67/ecommerce/internal/models"
	"github.com/oojoseph67/ecommerce/internal/utils/responses"
	"gorm.io/gorm"
)

func (s *ProductService) CreateProduct(req *dto.CreateProductRequest) (*dto.ProductResponse, error) {

	// check if category exists
	categoryModel, err := s.categoryService.getCategory(req.CategoryID)
	if err != nil {
		return nil, err
	}

	// check if sku exists
	var existingSku models.Product
	err = s.db.Where("sku = ?", req.SKU).First(&existingSku).Error
	if err == nil {
		return nil, errors.New("SKU value already exists")
	}

	// declare model
	productModel := models.Product{
		CategoryID:  categoryModel.ID,
		Name:        req.Name,
		Description: req.Description,
		Price:       req.Price,
		Stock:       int(req.Stock),
		SKU:         req.SKU,
	}

	// save model
	if err := s.db.Create(&productModel).Error; err != nil {
		s.internalLogger("product:create").Warn().Err(err).Msg("error creating product")
		return nil, err
	}

	product := s.convertToProductResponse(&productModel)

	return product, nil
}

func (s *ProductService) GetProducts(page, limit int) ([]dto.ProductResponse, *responses.PaginationMeta, error) {

	if page < 1 {
		page = 1
	}

	if limit < 1 {
		limit = 10
	}

	offset := (page - 1) * limit

	var productsModel []models.Product
	var total int64

	// db search // we use Model when we want to count, update or pluck || when we want to trigger model hooks BeforeFind, AfterFind|| Scopes() or Select()
	s.db.Where("is_active = ?", true).Find(&productsModel)
	s.db.Model(&models.Product{}).Where("is_active = ?", true).Count(&total)

	err := s.db.Preload("Category", "is_active = ?", true).
		Preload("Images").
		Where("is_active = ?", true).
		Offset(offset).
		Limit(limit).
		Find(&productsModel).
		Error

	if err != nil {
		s.internalLogger("product:get_products").Warn().Err(err).Msg("error getting products")
		return nil, nil, err
	}

	products := make([]dto.ProductResponse, len(productsModel))
	for i := range productsModel {
		products[i] = *s.convertToProductResponse(&productsModel[i])
	}

	totalPages := int(total+int64(limit)-1) / int(limit)
	meta := &responses.PaginationMeta{
		Page:       page,
		Limit:      limit,
		Total:      total,
		TotalPages: totalPages,
	}

	return products, meta, nil
}

func (s *ProductService) GetProduct(id string) (*dto.ProductResponse, error) {
	var productModel models.Product
	if err := s.db.Where("id = ?", id).First(&productModel).Error; err != nil {
		s.internalLogger("product:get_product").Warn().Str("product_id", id).Err(err).Msg("error finding product")
		return nil, errors.New("product not found")
	}

	response := s.convertToProductResponse(&productModel)

	return response, nil

	// return &productModel, nil
}

func (s *ProductService) UpdateProduct(id string, req *dto.UpdateProductRequest) (*dto.ProductResponse, error) {

	// check if the id exists
	var productModel models.Product
	if err := s.db.First(&productModel).Error; err != nil {
		return nil, err
	}

	// update values
	productModel.Name = coalesce(req.Name, productModel.Name)
	productModel.Description = coalesce(req.Description, productModel.Description)
	productModel.Price = coalesce(req.Price, productModel.Price)
	productModel.Stock = coalesce(req.Stock, productModel.Stock)

	if req.CategoryID != "" {

		if _, err := s.categoryService.getCategory(req.CategoryID); err != nil {
			return nil, err
		}

		productModel.CategoryID = coalesce(req.CategoryID, productModel.CategoryID)
	}

	if req.IsActive != nil {
		productModel.IsActive = coalescePtr(req.IsActive, productModel.IsActive)
	}

	// update the data
	err := s.db.Model(productModel).Updates(map[string]interface{}{
		"name":        productModel.Name,
		"description": productModel.Description,
		"price":       productModel.Price,
		"stock":       productModel.Stock,
		"category_id": productModel.CategoryID,
		"is_active":   productModel.IsActive,
	}).Error

	if err != nil {
		s.internalLogger("product:update_product").Warn().Str("product_id", id).Err(err).Msg("error updating product")
		return nil, err
	}

	product := s.convertToProductResponse(&productModel)

	return product, nil
}

func (s *ProductService) DeleteProduct(id string) error {
	// get products
	// productModel, err := s.GetProduct(id)
	// if err != nil {
	// 	return err
	// }

	var productModel models.Product
	err := s.db.Where("id = ?", id).First(&productModel).Error
	if err != nil {
		return errors.New("product not found")
	}

	// delete products
	// if err := s.db.Where("id = ?", productModel.ID).Delete(&productModel).Error; err != nil {
	// 	return errors.New("error deleting product")
	// }

	err = s.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&productModel).Update("is_active", false).Error; err != nil {
			s.internalLogger("product:delete_product").Warn().Str("product_id", id).Err(err).Msg("error deleting products")
			return err
		}
		// return tx.Unscoped().Delete(&categoryModel).Error // Unscoped = hard delete

		return tx.Delete(&productModel).Error
	})

	if err != nil {
		s.internalLogger("product:delete_product").Warn().Str("product_id", id).Err(err).Msg("error deleting products")
		return err
	}

	return nil
}

func (s *ProductService) convertToProductResponse(product *models.Product) *dto.ProductResponse {

	images := make([]dto.ProductImageResponse, len(product.Images))
	for i := range product.Images {
		images[i] = dto.ProductImageResponse{
			ID:        product.Images[i].ID,
			URL:       product.Images[i].URL,
			AltText:   product.Images[i].AltText,
			IsPrimary: product.Images[i].IsPrimary,
		}
	}

	category := dto.CategoryResponse{
		ID:          product.Category.ID,
		Name:        product.Category.Name,
		Description: product.Category.Description,
		IsActive:    product.Category.IsActive,
	}

	response := dto.ProductResponse{
		ID:          product.ID,
		Name:        product.Name,
		Description: product.Description,
		CategoryID:  product.CategoryID,
		Price:       product.Price,
		Stock:       product.Stock,
		SKU:         product.SKU,
		IsActive:    product.IsActive,
		Category:    category,
		Images:      images,
	}

	return &response
}
