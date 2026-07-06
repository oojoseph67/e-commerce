package handlers

import (
	_ "net/http"

	"github.com/gin-gonic/gin"
	"github.com/oojoseph67/ecommerce/internal/dto"
	"github.com/oojoseph67/ecommerce/internal/utils/responses"
	"github.com/oojoseph67/ecommerce/internal/utils/validators"
)

// AdminSignup godoc
// @Summary Admin registration
// @Description Register a new admin account with a valid admin code
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body dto.AdminSignupRequest true "Admin signup credentials"
// @Success 201 {object} responses.Response{data=dto.AuthResponse}
// @Failure 400 {object} responses.Response
// @Router /auth/admin/signup [post]
func (h *Handler) AdminSignup(ctx *gin.Context) {
	var req dto.AdminSignupRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		responses.BadRequestResponse(ctx, "Invalid request data", validators.FormatValidationError(err))
		return
	}

	response, err := h.authService.AdminSignup(&req)
	if err != nil {
		responses.BadRequestResponse(ctx, "Signup failed", err)
		return
	}

	responses.CreatedResponse(ctx, "Signup successful", response)
}

// AdminLogin godoc
// @Summary Admin login
// @Description Authenticate an admin account
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body dto.LoginRequest true "Admin login credentials"
// @Success 200 {object} responses.Response{data=dto.AuthResponse}
// @Failure 400 {object} responses.Response
// @Router /auth/admin/login [post]
func (h *Handler) AdminLogin(ctx *gin.Context) {
	var req dto.LoginRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		responses.BadRequestResponse(ctx, "Invalid request data", validators.FormatValidationError(err))
		return
	}

	response, err := h.authService.AdminLogin(&req)
	if err != nil {
		responses.BadRequestResponse(ctx, "Login failed", err)
		return
	}

	responses.SuccessResponse(ctx, "Login successful", response)
}

// Signup godoc
// @Summary User registration
// @Description Register a new user account
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body dto.SignupRequest true "User signup credentials"
// @Success 201 {object} responses.Response{data=dto.AuthResponse}
// @Failure 400 {object} responses.Response
// @Router /auth/signup [post]
func (h *Handler) Signup(ctx *gin.Context) {
	var req dto.SignupRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		responses.BadRequestResponse(ctx, "Invalid request data", validators.FormatValidationError(err))
		return
	}

	response, err := h.authService.Signup(&req)
	if err != nil {
		responses.BadRequestResponse(ctx, "Signup failed", err)
		return
	}

	responses.CreatedResponse(ctx, "Signup successful", response)
}

// Login godoc
// @Summary User login
// @Description Authenticate a user account
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body dto.LoginRequest true "User login credentials"
// @Success 200 {object} responses.Response{data=dto.AuthResponse}
// @Failure 400 {object} responses.Response
// @Router /auth/login [post]
func (h *Handler) Login(ctx *gin.Context) {
	var req dto.LoginRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		responses.BadRequestResponse(ctx, "Invalid request data", validators.FormatValidationError(err))
		return
	}

	response, err := h.authService.Login(&req)
	if err != nil {
		responses.BadRequestResponse(ctx, "Login failed", err)
		return
	}

	responses.SuccessResponse(ctx, "Login successful", response)
}

// RefreshToken godoc
// @Summary Refresh access token
// @Description Get a new access token using a refresh token
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body dto.RefreshTokenRequest true "Refresh token"
// @Success 200 {object} responses.Response{data=dto.AuthResponse}
// @Failure 400 {object} responses.Response
// @Router /auth/refresh [post]
func (h *Handler) RefreshToken(ctx *gin.Context) {
	var req dto.RefreshTokenRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		responses.BadRequestResponse(ctx, "Invalid request data", validators.FormatValidationError(err))
		return
	}

	response, err := h.authService.RefreshToken(&req)
	if err != nil {
		responses.BadRequestResponse(ctx, "Refresh token failed", err)
		return
	}

	responses.SuccessResponse(ctx, "Token refreshed", response)
}

// Logout godoc
// @Summary User logout
// @Description Invalidate the refresh token
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body dto.LogoutRequest true "Refresh token to invalidate"
// @Success 200 {object} responses.Response
// @Failure 400 {object} responses.Response
// @Router /auth/logout [post]
func (h *Handler) Logout(ctx *gin.Context) {
	var req dto.LogoutRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		responses.BadRequestResponse(ctx, "Invalid request data", validators.FormatValidationError(err))
		return
	}

	if err := h.authService.Logout(&req); err != nil {
		responses.BadRequestResponse(ctx, "Logout failed", err)
		return
	}

	responses.SuccessResponse(ctx, "Logout successful", nil)
}
