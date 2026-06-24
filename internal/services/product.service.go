package services

import (
	"errors"

	"github.com/oojoseph67/ecommerce/internal/dto"
	"github.com/oojoseph67/ecommerce/internal/models"
	"github.com/oojoseph67/ecommerce/internal/utils/responses"
)

/**** CATEGORY SERVICES */

func (s *Services) CreateCategory(req *dto.CreateCategoryRequest) (*dto.CategoryResponse, error) {

	// creating the model
	categoryModel := models.Category{
		Name:        req.Name,
		Description: req.Description,
	}

	// db create
	if err := s.db.Create(&categoryModel).Error; err != nil {
		s.internalLogger("category:create").Warn().Err(err).Msg("error creating new category")
		return nil, err
	}

	category := &dto.CategoryResponse{
		ID:          categoryModel.ID,
		Name:        categoryModel.Name,
		Description: categoryModel.Description,
		IsActive:    categoryModel.IsActive,
	}

	return category, nil
}

func (s *Services) GetCategories() ([]dto.CategoryResponse, error) {

	var categoriesModel []models.Category

	// db search
	if err := s.db.Where("is_active = ?", true).Find(&categoriesModel).Error; err != nil {
		s.internalLogger("category:get_categories").Warn().Err(err).Msg("error getting active categories")
		return nil, err
	}

	categories := make([]dto.CategoryResponse, len(categoriesModel)) // to keep the length of the array fixed

	for i := range categoriesModel {
		categories[i] = dto.CategoryResponse{
			ID:          categoriesModel[i].ID,
			Name:        categoriesModel[i].Name,
			Description: categoriesModel[i].Description,
			IsActive:    categoriesModel[i].IsActive,
		}
	}

	// this is memory expensive because it copies data instead
	// for i, category := range response {
	// 	response[i] = dto.CategoryResponse{
	// 		ID:          category.ID,
	// 		Name:        category.Name,
	// 		Description: category.Description,
	// 		IsActive:    category.IsActive,
	// 	}
	// }

	return categories, nil
}

func (s *Services) UpdateCategory(id string, req *dto.UpdateCategoryRequest) (*dto.CategoryResponse, error) {

	// check if the id exists
	categoryModel, err := s.getCategory(id)
	if err != nil {
		return nil, err
	}

	// update values
	categoryModel.Name = coalesce(req.Name, categoryModel.Name)
	categoryModel.Description = coalesce(req.Description, categoryModel.Description)
	if req.IsActive != nil {
		categoryModel.IsActive = coalesce(*req.IsActive, categoryModel.IsActive)
	}

	// update the data
	err = s.db.Model(categoryModel).Updates(map[string]interface{}{
		"name":        categoryModel.Name,
		"description": categoryModel.Description,
		"is_active":   categoryModel.IsActive,
	}).Error

	if err != nil {
		s.internalLogger("category:update_category").Warn().Str("category_id", id).Err(err).Msg("error updating category")
		return nil, err
	}

	category := &dto.CategoryResponse{
		ID:          categoryModel.ID,
		Name:        categoryModel.Name,
		Description: categoryModel.Description,
		IsActive:    categoryModel.IsActive,
	}

	return category, nil
}

func (s *Services) DeleteCategory(id string) error {
	// get category
	categoryModel, err := s.getCategory(id)
	if err != nil {
		return err
	}

	// delete category
	if err := s.db.Where("id = ?", categoryModel.ID).Delete(&categoryModel).Error; err != nil {
		s.internalLogger("category:delete_category").Warn().Str("category_id", id).Err(err).Msg("error deleting category")
		return errors.New("error deleting category")
	}

	return nil
}

func (s *Services) getCategory(id string) (*models.Category, error) {
	var categoryModel models.Category
	if err := s.db.Where("id = ?", id).First(&categoryModel).Error; err != nil {
		s.internalLogger("category:get_category").Warn().Str("category_id", id).Err(err).Msg("error finding category")
		return nil, errors.New("category not found")
	}

	return &categoryModel, nil
}

/**** PRODUCT SERVICES */

func (s *Services) CreateProduct(req *dto.CreateProductRequest) (*dto.ProductResponse, error) {

	// check if category exists
	categoryModel, err := s.getCategory(req.CategoryID)
	if err != nil {
		return nil, err
	}

	// declare model
	productModel := models.Product{
		CategoryID:  categoryModel.ID,
		Name:        req.Name,
		Description: req.Description,
		Price:       req.Price,
		Stock:       int(req.Price),
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

func (s *Services) GetProducts(page, limit int) ([]dto.ProductResponse, *responses.PaginationMeta, error) {

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

func (s *Services) GetProduct(id string) (*dto.ProductResponse, error) {
	var productModel models.Product
	if err := s.db.Where("id = ?", id).First(&productModel).Error; err != nil {
		s.internalLogger("product:get_product").Warn().Str("product_id", id).Err(err).Msg("error finding product")
		return nil, errors.New("product not found")
	}

	response := s.convertToProductResponse(&productModel)

	return response, nil

	// return &productModel, nil
}

func (s *Services) UpdateProduct(id string, req *dto.UpdateProductRequest) (*dto.ProductResponse, error) {

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

	if productModel.CategoryID != "" {

		if _, err := s.getCategory(req.CategoryID); err != nil {
			return nil, err
		}

		productModel.CategoryID = coalesce(req.CategoryID, productModel.CategoryID)
	}

	if req.IsActive != nil {
		productModel.IsActive = coalesce(*req.IsActive, productModel.IsActive)
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

func (s *Services) DeleteProduct(id string) error {
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
	if err := s.db.Where("id = ?", productModel.ID).Delete(&productModel).Error; err != nil {
		s.internalLogger("product:delete_product").Warn().Str("product_id", id).Err(err).Msg("error deleting products")
		return errors.New("error deleting product")
	}

	return nil
}

func (s *Services) convertToProductResponse(product *models.Product) *dto.ProductResponse {

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
