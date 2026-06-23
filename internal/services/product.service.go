package services

import (
	"errors"

	"github.com/oojoseph67/ecommerce/internal/dto"
	"github.com/oojoseph67/ecommerce/internal/models"
)

func (s *Services) CreateCategory(req *dto.CreateCategoryRequest) (*dto.CategoryResponse, error) {

	// creating the model
	categoryModel := models.Category{
		Name:        req.Name,
		Description: req.Description,
	}

	// db create
	if err := s.db.Create(&categoryModel).Error; err != nil {
		s.internalLogger("category").Warn().Err(err).Msg("error creating new category")
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
		s.internalLogger("category").Warn().Err(err).Msg("error getting active categories")
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
	var categoryModel *models.Category

	// if err := s.db.Where("id = ?", id).First(&categoryModel).Error; err != nil {
	// 	s.internalLogger("category").Warn().Str("category_id", id).Err(err).Msg("error finding category")
	// 	return nil, errors.New("category not found")
	// }

	categoryModel, err := s.getCategory(id, categoryModel)
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
		s.internalLogger("category").Warn().Str("category_id", id).Err(err).Msg("error updating category")
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
	var categoryModel *models.Category

	categoryModel, err := s.getCategory(id, categoryModel)
	if err != nil {
		return err
	}

	// delete category
	if err := s.db.Where("id = ?", categoryModel.ID).Delete(&categoryModel).Error; err != nil {
		s.internalLogger("category").Warn().Str("category_id", id).Err(err).Msg("error deleting category")
		return errors.New("error deleting category")
	}

	return nil
}

func (s *Services) getCategory(id string, categoryModel *models.Category) (*models.Category, error) {
	if err := s.db.Where("id = ?", id).First(&categoryModel).Error; err != nil {
		s.internalLogger("category").Warn().Str("category_id", id).Err(err).Msg("error finding category")
		return nil, errors.New("category not found")
	}

	return categoryModel, nil
}
