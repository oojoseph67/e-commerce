package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/oojoseph67/ecommerce/internal/dto"
	"github.com/oojoseph67/ecommerce/internal/utils/responses"
	"github.com/oojoseph67/ecommerce/internal/utils/validators"
)

// Me godoc
// @Summary Get current user profile
// @Description Retrieve the authenticated user's profile information
// @Tags User
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} responses.Response{data=dto.UserResponse}
// @Failure 401 {object} responses.Response
// @Failure 404 {object} responses.Response
// @Router /user/me [get]
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

// UpdateProfile godoc
// @Summary Update user profile
// @Description Update the authenticated user's profile information
// @Tags User
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body dto.UpdateProfileRequest true "Profile update data"
// @Success 200 {object} responses.Response{data=dto.UserResponse}
// @Failure 400 {object} responses.Response
// @Failure 401 {object} responses.Response
// @Router /user/update [patch]
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
