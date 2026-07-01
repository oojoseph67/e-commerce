package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/oojoseph67/ecommerce/internal/dto"
	"github.com/oojoseph67/ecommerce/internal/utils/responses"
	"github.com/oojoseph67/ecommerce/internal/utils/validators"
)

func (h *Handler) Me(ctx *gin.Context) {
	userId, ok := getUserId(ctx)
	if !ok {
		return
	}

	profile, err := h.userService.GetUserProfile(userId)
	if err != nil {
		responses.NotFoundResponse(ctx, "User not found", err)
		return
	}

	responses.SuccessResponse(ctx, "Profile retrieved successfully", profile)
}

func (h *Handler) UpdateProfile(ctx *gin.Context) {
	userId, ok := getUserId(ctx)
	if !ok {
		return
	}

	var req *dto.UpdateProfileRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		responses.BadRequestResponse(ctx, "Invalid request data", validators.FormatValidationError(err))
		return
	}

	profile, err := h.userService.UpdateProfile(userId, req)
	if err != nil {
		responses.InternalServerResponse(ctx, "Couldnt update profile", err)
		return
	}

	responses.SuccessResponse(ctx, "Profile updated successfully", profile)
}
