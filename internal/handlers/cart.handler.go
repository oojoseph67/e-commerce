package handlers

import (
	"errors"

	"github.com/gin-gonic/gin"
	"github.com/oojoseph67/ecommerce/internal/dto"
	"github.com/oojoseph67/ecommerce/internal/middleware"
	"github.com/oojoseph67/ecommerce/internal/utils/responses"
	"github.com/oojoseph67/ecommerce/internal/utils/validators"
)

func (h *Handler) GetCart(ctx *gin.Context) {
	userId := ctx.GetString(middleware.UserIdAuthKey)
	if userId == "" {
		responses.BadRequestResponse(ctx, "user_id not received", errors.New("user_id not received"))
		return
	}

	userCart, err := h.cartService.GetCart(userId)
	if err != nil {
		responses.NotFoundResponse(ctx, "cart not found", err)
	}

	responses.SuccessResponse(ctx, "Cart retrieved successfully", userCart)
}

func (h *Handler) AddItemToCart(ctx *gin.Context) {
	userId := ctx.GetString(middleware.UserIdAuthKey)
	var req *dto.AddToCartRequest

	if userId == "" {
		responses.BadRequestResponse(ctx, "user_id not received", errors.New("user_id not received"))
		return
	}

	if err := ctx.ShouldBindJSON(&req); err != nil {
		responses.BadRequestResponse(ctx, "Invalid request data", validators.FormatValidationError(err))
		return
	}

	addToCart, err := h.cartService.AddToCart(req, userId)
	if err != nil {
		responses.InternalServerResponse(ctx, "Couldnt add item to cart", err)
		return
	}

	responses.SuccessResponse(ctx, "Item added", addToCart)
}

func (h *Handler) UpdateCartItem(ctx *gin.Context) {
	id := ctx.Param("id")
	userId := ctx.GetString(middleware.UserIdAuthKey)
	var req *dto.UpdateCartItemRequest

	if id == "" {
		responses.BadRequestResponse(ctx, "Please provide id param", errors.New("please provide id param"))
	}

	if userId == "" {
		responses.BadRequestResponse(ctx, "user_id not received", errors.New("user_id not received"))
		return
	}

	if err := ctx.ShouldBindJSON(&req); err != nil {
		responses.BadRequestResponse(ctx, "Invalid request data", validators.FormatValidationError(err))
		return
	}

	updateResponse, err := h.cartService.UpdateCartItem(req, id, userId)
	if err != nil {
		responses.InternalServerResponse(ctx, "Unable to update cart item", err)
	}

	responses.SuccessResponse(ctx, "Cart item updated successfully", updateResponse)

}

func (h *Handler) RemoveCartItem(ctx *gin.Context) {
	id := ctx.Param("id")
	userId := ctx.GetString(middleware.UserIdAuthKey)

	if id == "" {
		responses.BadRequestResponse(ctx, "Please provide id param", errors.New("please provide id param"))
	}

	if userId == "" {
		responses.BadRequestResponse(ctx, "user_id not received", errors.New("user_id not received"))
		return
	}

	err := h.cartService.RemoveCartItem(id, userId)
	if err != nil {
		responses.InternalServerResponse(ctx, "couldnt remove item from cart", err)
	}

	responses.SuccessResponse(ctx, "Item removed from cart successfully", nil)

}
