package services

import (
	"github.com/oojoseph67/ecommerce/internal/dto"
	"github.com/oojoseph67/ecommerce/internal/models"
)

func (s *Services) GetUserProfile(userId string) (*dto.UserResponse, error) {
	var userModel models.User

	if err := s.db.Where("id = ?", userId).First(&userModel).Error; err != nil {
		s.internalLogger("user").Warn().Str("user_id", userId).Err(err).Msg("user not found")
		return nil, err
	}

	user := &dto.UserResponse{
		ID:        userModel.ID,
		Email:     userModel.Email,
		FirstName: userModel.FirstName,
		LastName:  userModel.LastName,
		Phone:     userModel.Phone,
		Role:      string(userModel.Role),
		IsActive:  userModel.IsActive,
	}

	return user, nil
}

func (s *Services) UpdateProfile(userId string, req *dto.UpdateProfileRequest) (*dto.UserResponse, error) {
	var userModel models.User

	// find user
	if err := s.db.Where("id = ?", userId).First(&userModel).Error; err != nil {
		s.internalLogger("user").Warn().Str("user_id", userId).Err(err).Msg("user not found")
		return nil, err
	}

	// update values
	userModel.FirstName = coalesce(req.FirstName, userModel.FirstName)
	userModel.LastName = coalesce(req.LastName, userModel.LastName)
	userModel.Phone = coalesce(req.Phone, userModel.Phone)

	// save users
	err := s.db.Model(userModel).Updates(map[string]interface{}{
		"first_name": userModel.FirstName,
		"last_name":  userModel.LastName,
		"phone":      userModel.Phone,
	}).Error

	if err != nil {
		s.internalLogger("user").Warn().Str("user_id", userId).Err(err).Msg("error updating user profile")
		return nil, err
	}

	user, err := s.GetUserProfile(userModel.ID)
	if err != nil {
		return nil, err
	}

	return user, nil
}
