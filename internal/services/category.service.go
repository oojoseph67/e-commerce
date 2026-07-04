package services

import (
	"errors"
	"fmt"
	"strings"

	"github.com/oojoseph67/ecommerce/internal/dto"
	"github.com/oojoseph67/ecommerce/internal/models"
	"gorm.io/gorm"
)

func (s *CategoryService) CreateCategory(req *dto.CreateCategoryRequest) (*dto.CategoryResponse, error) {

	var existingCategory *models.Category
	lowercaseName := strings.ToLower(req.Name)

	if err := s.db.Where("name = ?", lowercaseName).First(&existingCategory).Error; err == nil {
		return nil, fmt.Errorf("category with name: '%s' already exists", lowercaseName)
	}

	// creating the model
	categoryModel := models.Category{
		Name:        lowercaseName,
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

func (s *CategoryService) GetCategories() (*dto.GetCategoriesResponse, error) {

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

	response := &dto.GetCategoriesResponse{
		Categories: categories,
		Total:      len(categories),
	}

	return response, nil
}

func (s *CategoryService) UpdateCategory(id string, req *dto.UpdateCategoryRequest) (*dto.CategoryResponse, error) {

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

func (s *CategoryService) DeleteCategory(id string) error {
	// get category
	categoryModel, err := s.getCategory(id)
	if err != nil {
		return err
	}

	// delete category
	// if err := s.db.Where("id = ?", categoryModel.ID).Delete(&categoryModel).Error; err != nil {
	// 	s.internalLogger("category:delete_category").Warn().Str("category_id", id).Err(err).Msg("error deleting category")
	// 	return errors.New("error deleting category")
	// }

	err = s.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&categoryModel).Update("is_active", false).Error; err != nil {
			s.internalLogger("category:delete_category").Warn().Str("category_id", id).Err(err).Msg("error deleting category")
			return err
		}
		// return tx.Unscoped().Delete(&categoryModel).Error // Unscoped = hard delete

		return tx.Delete(&categoryModel).Error
	})

	if err != nil {
		s.internalLogger("category:delete_category").Warn().Str("category_id", id).Err(err).Msg("error deleting category")
		return err
	}

	return nil
}

func (s *CategoryService) getCategory(id string) (*models.Category, error) {
	var categoryModel models.Category
	if err := s.db.Where("id = ?", id).First(&categoryModel).Error; err != nil {
		s.internalLogger("category:get_category").Warn().Str("category_id", id).Err(err).Msg("error finding category")
		return nil, errors.New("category not found")
	}

	return &categoryModel, nil
}
