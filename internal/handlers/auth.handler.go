package handlers

import (
	_ "net/http"

	"github.com/gin-gonic/gin"
	"github.com/oojoseph67/ecommerce/internal/dto"
	"github.com/oojoseph67/ecommerce/internal/utils/responses"
	"github.com/oojoseph67/ecommerce/internal/utils/validators"
)

func (h *Handler) AdminSignup(ctx *gin.Context) {
	service := h.service.AuthService
	var req dto.AdminSignupRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		responses.BadRequestResponse(ctx, "Invalid request data", validators.FormatValidationError(err))
		return
	}

	response, err := service.AdminSignup(&req)
	if err != nil {
		responses.BadRequestResponse(ctx, "Signup failed", err)
		return
	}

	responses.CreatedResponse(ctx, "Signup successful", response)
}

func (h *Handler) AdminLogin(ctx *gin.Context) {
	service := h.service.AuthService
	var req dto.LoginRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		responses.BadRequestResponse(ctx, "Invalid request data", validators.FormatValidationError(err))
		return
	}

	response, err := service.AdminLogin(&req)
	if err != nil {
		responses.BadRequestResponse(ctx, "Login failed", err)
		return
	}

	responses.SuccessResponse(ctx, "Login successful", response)
}

func (h *Handler) Signup(ctx *gin.Context) {
	service := h.service.AuthService
	var req dto.SignupRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		responses.BadRequestResponse(ctx, "Invalid request data", validators.FormatValidationError(err))
		return
	}

	response, err := service.Signup(&req)
	if err != nil {
		responses.BadRequestResponse(ctx, "Signup failed", err)
		return
	}

	responses.CreatedResponse(ctx, "Signup successful", response)
}

func (h *Handler) Login(ctx *gin.Context) {
	service := h.service.AuthService
	var req dto.LoginRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		responses.BadRequestResponse(ctx, "Invalid request data", validators.FormatValidationError(err))
		return
	}

	response, err := service.Login(&req)
	if err != nil {
		responses.BadRequestResponse(ctx, "Login failed", err)
		return
	}

	responses.SuccessResponse(ctx, "Login successful", response)
}

func (h *Handler) RefreshToken(ctx *gin.Context) {
	service := h.service.AuthService
	var req dto.RefreshTokenRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		responses.BadRequestResponse(ctx, "Invalid request data", validators.FormatValidationError(err))
		return
	}

	response, err := service.RefreshToken(&req)
	if err != nil {
		responses.BadRequestResponse(ctx, "Refresh token failed", err)
		return
	}

	responses.SuccessResponse(ctx, "Token refreshed", response)
}

func (h *Handler) Logout(ctx *gin.Context) {
	service := h.service.AuthService
	var req dto.LogoutRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		responses.BadRequestResponse(ctx, "Invalid request data", validators.FormatValidationError(err))
		return
	}

	if err := service.Logout(&req); err != nil {
		responses.BadRequestResponse(ctx, "Logout failed", err)
		return
	}

	responses.SuccessResponse(ctx, "Logout successful", nil)
}
