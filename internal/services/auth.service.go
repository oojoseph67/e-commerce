package services

import (
	"errors"
	"strings"
	"time"

	"github.com/oojoseph67/ecommerce/internal/dto"
	"github.com/oojoseph67/ecommerce/internal/models"
	jwtp "github.com/oojoseph67/ecommerce/internal/utils/jwt"
	"github.com/oojoseph67/ecommerce/internal/utils/password"
)

// Admin Signup
func (s *Services) AdminSignup(req *dto.AdminSignupRequest) (*dto.AuthResponse, error) {
	adminConfig := s.config.Admin

	lowercaseEmail := strings.ToLower(req.Email)

	// check email is a company email
	emailParts := strings.Split(lowercaseEmail, "@")

	if emailParts[1] != adminConfig.DomainName {
		s.internalLogger("auth:admin_signup").Warn().Str("email_domain_name", emailParts[1]).Msg("invalid email domain")
		return nil, errors.New("invalid email")
	}

	if strings.ToLower(req.Code) != adminConfig.AdminCode {
		s.internalLogger("auth:admin_signup").Warn().Str("admin_code", adminConfig.AdminCode).Msg("invalid admin code")
		return nil, errors.New("invalid admin code")
	}

	// check if user already exists
	var existingUserByEmail models.User
	if err := s.db.Where("email = ?", lowercaseEmail).First(&existingUserByEmail).Error; err == nil {
		s.internalLogger("auth:admin_signup").Warn().Str("email", lowercaseEmail).Msg("signup attempted for existing user")
		return nil, errors.New("user already exists with this email")
	}

	// hash password
	hashedPassword, err := password.HashPassword(req.Password)
	if err != nil {
		s.internalLogger("auth:admin_signup").Error().Err(err).Msg("failed to hash password")
		return nil, errors.New("failed to process password")
	}

	// create user
	user := models.User{
		Email:     lowercaseEmail,
		Password:  hashedPassword,
		FirstName: "admin",
		LastName:  "admin",
		Phone:     "12345678900",
		Role:      models.UserRoleAdmin,
	}

	if err := s.db.Create(&user).Error; err != nil {
		s.internalLogger("auth:admin_signup").Error().Err(err).Str("email", lowercaseEmail).Msg("failed to create user")
		return nil, errors.New("failed to create admin")
	}

	s.internalLogger("auth:admin_signup").Info().Str("user_id", user.ID).Str("email", lowercaseEmail).Msg("user created")

	// create cart
	cart := models.Cart{
		UserID: user.ID,
	}
	if err := s.db.Create(&cart).Error; err != nil {
		s.internalLogger("auth:admin_signup").Warn().Err(err).Str("user_id", user.ID).Msg("failed to create cart for new user")
	}

	authResponse, err := s.generateAuthResponse(&user)
	if err != nil {
		s.internalLogger("auth:admin_signup").Error().Err(err).Str("user_id", user.ID).Msg("failed to generate auth response after signup")
		return nil, err
	}

	return authResponse, nil
}

// Signup registers a new user.
func (s *Services) Signup(req *dto.SignupRequest) (*dto.AuthResponse, error) {
	lowercaseEmail := strings.ToLower(req.Email)

	// check if user already exists
	var existingUserByEmail models.User
	if err := s.db.Where("email = ?", lowercaseEmail).First(&existingUserByEmail).Error; err == nil {
		s.internalLogger("auth:signup").Warn().Str("email", lowercaseEmail).Msg("signup attempted for existing user")
		return nil, errors.New("user already exists with this email")
	}

	// hash password
	hashedPassword, err := password.HashPassword(req.Password)
	if err != nil {
		s.internalLogger("auth:signup").Error().Err(err).Msg("failed to hash password")
		return nil, errors.New("failed to process password")
	}

	// create user
	user := models.User{
		Email:     lowercaseEmail,
		Password:  hashedPassword,
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Phone:     req.Phone,
		Role:      models.UserRoleCustomer,
	}

	if err := s.db.Create(&user).Error; err != nil {
		s.internalLogger("auth:signup").Error().Err(err).Str("email", lowercaseEmail).Msg("failed to create user")
		return nil, errors.New("failed to create user")
	}

	s.internalLogger("auth:signup").Info().Str("user_id", user.ID).Str("email", lowercaseEmail).Msg("user created")

	// create cart
	cart := models.Cart{
		UserID: user.ID,
	}
	if err := s.db.Create(&cart).Error; err != nil {
		s.internalLogger("auth:signup").Warn().Err(err).Str("user_id", user.ID).Msg("failed to create cart for new user")
	}

	authResponse, err := s.generateAuthResponse(&user)
	if err != nil {
		s.internalLogger("auth:signup").Error().Err(err).Str("user_id", user.ID).Msg("failed to generate auth response after signup")
		return nil, err
	}

	return authResponse, nil
}

// Login authenticates a user.
func (s *Services) Login(req *dto.LoginRequest) (*dto.AuthResponse, error) {
	lowercaseEmail := strings.ToLower(req.Email)

	var user models.User
	if err := s.db.Where("email = ? AND is_active = ?", lowercaseEmail, true).First(&user).Error; err != nil {
		s.internalLogger("auth:login").Warn().Str("email", lowercaseEmail).Msg("login failed: user not found or inactive")
		return nil, errors.New("incorrect email or password")
	}

	if !password.ComparePassword(req.Password, user.Password) {
		s.internalLogger("auth:login").Warn().Str("user_id", user.ID).Str("email", lowercaseEmail).Msg("login failed: incorrect password")
		return nil, errors.New("incorrect email or password")
	}

	s.internalLogger("auth:login").Info().Str("user_id", user.ID).Str("email", lowercaseEmail).Msg("user logged in")

	authResponse, err := s.generateAuthResponse(&user)
	if err != nil {
		s.internalLogger("auth:login").Error().Err(err).Str("user_id", user.ID).Msg("failed to generate auth response after login")
		return nil, err
	}

	return authResponse, nil
}

// Admin Login
func (s *Services) AdminLogin(req *dto.LoginRequest) (*dto.AuthResponse, error) {
	lowercaseEmail := strings.ToLower(req.Email)

	var user models.User
	if err := s.db.Where("email = ? AND is_active = ?", lowercaseEmail, true).First(&user).Error; err != nil {
		s.internalLogger("auth:admin_login").Warn().Str("email", lowercaseEmail).Msg("login failed: user not found or inactive")
		return nil, errors.New("incorrect email or password")
	}

	if user.Role != models.UserRoleAdmin {
		return nil, errors.New("admin route")
	}

	if !password.ComparePassword(req.Password, user.Password) {
		s.internalLogger("auth:admin_login").Warn().Str("user_id", user.ID).Str("email", lowercaseEmail).Msg("login failed: incorrect password")
		return nil, errors.New("incorrect email or password")
	}

	s.internalLogger("auth:admin_login").Info().Str("user_id", user.ID).Str("email", lowercaseEmail).Msg("user logged in")

	authResponse, err := s.generateAuthResponse(&user)
	if err != nil {
		s.internalLogger("auth:admin_login").Error().Err(err).Str("user_id", user.ID).Msg("failed to generate auth response after login")
		return nil, err
	}

	return authResponse, nil
}

// RefreshToken exchanges a valid refresh token for a new auth pair.
func (s *Services) RefreshToken(req *dto.RefreshTokenRequest) (*dto.AuthResponse, error) {
	claims, err := jwtp.ValidateToken(req.RefreshToken, s.config.JWT.Secret)
	if err != nil {
		s.internalLogger("auth:refresh_token").Warn().Msg("refresh token failed: invalid token")
		return nil, errors.New("invalid refresh token")
	}

	var refreshToken models.RefreshToken
	now := time.Now()
	if err := s.db.Where("token = ? AND expires_at > ?", req.RefreshToken, now).First(&refreshToken).Error; err != nil {
		s.internalLogger("auth:refresh_token").Warn().Str("user_id", claims.UserID).Msg("refresh token failed: not found or expired")
		return nil, errors.New("refresh token not found or expired")
	}

	var user models.User
	if err := s.db.Where("id = ?", claims.UserID).First(&user).Error; err != nil {
		s.internalLogger("auth:refresh_token").Error().Err(err).Str("user_id", claims.UserID).Msg("refresh token failed: user not found")
		return nil, errors.New("user not found")
	}

	// delete(soft) old refresh token
	if err := s.db.Delete(&refreshToken).Error; err != nil {
		s.internalLogger("auth:refresh_token").Warn().Err(err).Str("user_id", user.ID).Msg("failed to delete old refresh token")
	}

	authResponse, err := s.generateAuthResponse(&user)
	if err != nil {
		s.internalLogger("auth:refresh_token").Error().Err(err).Str("user_id", user.ID).Msg("failed to generate auth response during refresh")
		return nil, err
	}

	s.internalLogger("auth:refresh_token").Info().Str("user_id", user.ID).Msg("token refreshed")
	return authResponse, nil
}

// Logout invalidates a refresh token.
func (s *Services) Logout(req *dto.LogoutRequest) error {
	var refreshToken models.RefreshToken
	result := s.db.Where("token = ?", req.RefreshToken).Delete(&refreshToken)
	if result.Error != nil {
		s.internalLogger("auth:logout").Error().Err(result.Error).Msg("failed to delete refresh token on logout")
		return result.Error
	}

	if result.RowsAffected == 0 {
		s.internalLogger("auth:logout").Warn().Msg("logout attempted with unknown refresh token")
	} else {
		s.internalLogger("auth:logout").Info().Msg("user logged out")
	}

	return nil
}

func (s *Services) generateAuthResponse(user *models.User) (*dto.AuthResponse, error) {
	accessToken, refreshToken, err := jwtp.GenerateTokenPair(
		&s.config.JWT, user.ID, user.Email, string(user.Role),
	)
	if err != nil {
		s.internalLogger("auth:generate_auth_response").Error().Err(err).Str("user_id", user.ID).Msg("failed to generate token pair")
		return nil, errors.New("failed to generate tokens")
	}

	refreshTokenModel := models.RefreshToken{
		UserID:    user.ID,
		Token:     refreshToken,
		ExpiresAt: time.Now().Add(s.config.JWT.RefreshTokenExpires),
	}
	if err := s.db.Create(&refreshTokenModel).Error; err != nil {
		s.internalLogger("auth:generate_auth_response").Error().Err(err).Str("user_id", user.ID).Msg("failed to save refresh token")
	}

	userModel := dto.UserResponse{
		ID:        user.ID,
		Email:     user.Email,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Phone:     user.Phone,
		Role:      string(user.Role),
		IsActive:  user.IsActive,
	}

	return &dto.AuthResponse{
		User:         userModel,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}
