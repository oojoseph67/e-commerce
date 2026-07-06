package handlers

import (
	"errors"

	"github.com/gin-gonic/gin"
	"github.com/oojoseph67/ecommerce/internal/dto"
	"github.com/oojoseph67/ecommerce/internal/utils/responses"
	"github.com/oojoseph67/ecommerce/internal/utils/validators"
)

// GetCart godoc
// @Summary Get user cart
// @Description Retrieve the authenticated user's shopping cart
// @Tags Cart
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} responses.Response{data=dto.CartResponse}
// @Failure 401 {object} responses.Response
// @Failure 404 {object} responses.Response
// @Router /cart/ [get]
func (h *Handler) GetCart(ctx *gin.Context) {
	userId, ok := getUserId(ctx)
	if !ok {
		return
	}

	userCart, err := h.cartService.GetCart(userId)
	if err != nil {
		responses.NotFoundResponse(ctx, "cart not found", err)
		return
	}

	responses.SuccessResponse(ctx, "Cart retrieved successfully", userCart)
}

// AddItemToCart godoc
// @Summary Add item to cart
// @Description Add a product to the authenticated user's shopping cart
// @Tags Cart
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body dto.AddToCartRequest true "Product to add to cart"
// @Success 200 {object} responses.Response{data=dto.CartResponse}
// @Failure 400 {object} responses.Response
// @Failure 401 {object} responses.Response
// @Router /cart/items [post]
func (h *Handler) AddItemToCart(ctx *gin.Context) {
	userId, ok := getUserId(ctx)
	if !ok {
		return
	}

	var req *dto.AddToCartRequest
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

// UpdateCartItem godoc
// @Summary Update cart item quantity
// @Description Update the quantity of an item in the shopping cart
// @Tags Cart
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Cart Item ID"
// @Param request body dto.UpdateCartItemRequest true "Updated quantity"
// @Success 200 {object} responses.Response{data=dto.CartItemResponse}
// @Failure 400 {object} responses.Response
// @Failure 401 {object} responses.Response
// @Router /cart/items/{id} [patch]
func (h *Handler) UpdateCartItem(ctx *gin.Context) {
	id := ctx.Param("id")
	if id == "" {
		responses.BadRequestResponse(ctx, "Please provide id param", errors.New("please provide id param"))
		return
	}

	userId, ok := getUserId(ctx)
	if !ok {
		return
	}

	var req *dto.UpdateCartItemRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		responses.BadRequestResponse(ctx, "Invalid request data", validators.FormatValidationError(err))
		return
	}

	updateResponse, err := h.cartService.UpdateCartItem(req, id, userId)
	if err != nil {
		responses.InternalServerResponse(ctx, "Unable to update cart item", err)
		return
	}

	responses.SuccessResponse(ctx, "Cart item updated successfully", updateResponse)
}

// RemoveCartItem godoc
// @Summary Remove item from cart
// @Description Remove an item from the shopping cart
// @Tags Cart
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Cart Item ID"
// @Success 200 {object} responses.Response
// @Failure 400 {object} responses.Response
// @Failure 401 {object} responses.Response
// @Router /cart/items/{id} [delete]
func (h *Handler) RemoveCartItem(ctx *gin.Context) {
	id := ctx.Param("id")
	if id == "" {
		responses.BadRequestResponse(ctx, "Please provide id param", errors.New("please provide id param"))
		return
	}

	userId, ok := getUserId(ctx)
	if !ok {
		return
	}

	err := h.cartService.RemoveCartItem(id, userId)
	if err != nil {
		responses.InternalServerResponse(ctx, "couldnt remove item from cart", err)
		return
	}

	responses.SuccessResponse(ctx, "Item removed from cart successfully", nil)
}
