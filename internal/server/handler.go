package server

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/oojoseph67/ecommerce/internal/dto"
	"github.com/oojoseph67/ecommerce/internal/services"
	"github.com/oojoseph67/ecommerce/internal/utils/responses"
)

func (s *Server) healthCheck(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{"status": "ok"})
}

func (s *Server) signup(ctx *gin.Context) {

	var req dto.SignupRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		responses.BadRequestResponse(ctx, "Invalid request data", err)
		return
	}

	authService := services.NewAuthService(s.db, s.config, s.logger)
	response, err := authService.Signup(&req)

	if err != nil {
		responses.BadRequestResponse(ctx, "signup failed", err)
		return
	}

	responses.CreatedResponse(ctx, "signup successful", response)
}

func (s *Server) login(ctx *gin.Context) {

	var req dto.LoginRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		responses.BadRequestResponse(ctx, "Invalid request data", err)
		return
	}

	authService := services.NewAuthService(s.db, s.config, s.logger)
	response, err := authService.Login(&req)

	if err != nil {
		responses.BadRequestResponse(ctx, "login failed", err)
		return
	}

	responses.SuccessResponse(ctx, "login successful", response)
}

func (s *Server) refreshToken(ctx *gin.Context) {

	var req dto.RefreshTokenRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		responses.BadRequestResponse(ctx, "Invalid request data", err)
		return
	}

	authService := services.NewAuthService(s.db, s.config, s.logger)
	response, err := authService.RefreshToken(&req)

	if err != nil {
		responses.BadRequestResponse(ctx, "requesting refresh_token failed", err)
		return
	}

	responses.SuccessResponse(ctx, "refresh_token requested", response)
}

func (s *Server) logout(ctx *gin.Context) {

	var req dto.LogoutRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		responses.BadRequestResponse(ctx, "Invalid request data", err)
		return
	}

	authService := services.NewAuthService(s.db, s.config, s.logger)
	err := authService.Logout(&req)

	if err != nil {
		responses.BadRequestResponse(ctx, "logout failed", err)
		return
	}

	responses.SuccessResponse(ctx, "logout successful", "logout successful")
}
