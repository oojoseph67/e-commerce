package handlers

import (
	"errors"
	_ "net/http"

	"github.com/gin-gonic/gin"
	"github.com/oojoseph67/ecommerce/internal/dto"
	"github.com/oojoseph67/ecommerce/internal/middleware"
	"github.com/oojoseph67/ecommerce/internal/utils/responses"
	"github.com/oojoseph67/ecommerce/internal/utils/validators"
)

func (h *Handler) Me(ctx *gin.Context) {
	service := h.service.UserService
	userId := ctx.GetString(middleware.UserIdAuthKey)
	if userId == "" {
		responses.BadRequestResponse(ctx, "user_id not received", errors.New("user_id not received"))
		return
	}

	profile, err := service.GetUserProfile(userId)
	if err != nil {
		responses.NotFoundResponse(ctx, "User not found", err)
		return
	}

	responses.SuccessResponse(ctx, "Profile retrieved successfully", profile)
}

func (h *Handler) UpdateProfile(ctx *gin.Context) {
	service := h.service.UserService
	userId := ctx.GetString(middleware.UserIdAuthKey)
	var req *dto.UpdateProfileRequest

	if userId == "" {
		responses.BadRequestResponse(ctx, "user_id not received", errors.New("user_id not received"))
		return
	}

	if err := ctx.ShouldBindJSON(&req); err != nil {
		responses.BadRequestResponse(ctx, "Invalid request data", validators.FormatValidationError(err))
		return
	}

	profile, err := service.UpdateProfile(userId, req)
	if err != nil {
		responses.InternalServerResponse(ctx, "Couldnt update profile", err)
		return
	}

	responses.SuccessResponse(ctx, "Profile updated successfully", profile)
}
