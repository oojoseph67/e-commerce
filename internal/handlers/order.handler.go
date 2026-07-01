package handlers

import (
	"errors"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/oojoseph67/ecommerce/internal/middleware"
	"github.com/oojoseph67/ecommerce/internal/utils/responses"
)

func (h *Handler) CreateOrder(ctx *gin.Context) {
	userId := ctx.GetString(middleware.UserIdAuthKey)
	if userId == "" {
		responses.BadRequestResponse(ctx, "user_id not received", errors.New("user_id not received"))
		return
	}

	order, err := h.orderService.CreateOrder(userId)
	if err != nil {
		responses.InternalServerResponse(ctx, "Couldnt create order", err)
		return
	}

	responses.CreatedResponse(ctx, "Order created successfully", order)

}

func (h *Handler) GetOrder(ctx *gin.Context) {
	id := ctx.Param("id")
	userId := ctx.GetString(middleware.UserIdAuthKey)

	if id == "" {
		responses.BadRequestResponse(ctx, "Please provide id param", errors.New("please provide id param"))
		return
	}

	if userId == "" {
		responses.BadRequestResponse(ctx, "user_id not received", errors.New("user_id not received"))
		return
	}

	order, err := h.orderService.GetOrder(userId, id)
	if err != nil {
		responses.InternalServerResponse(ctx, "Coudnlt get user order", err)
		return
	}

	responses.SuccessResponse(ctx, "Order retrieved successfully", order)
}

func (h *Handler) GetOrders(ctx *gin.Context) {
	userId := ctx.GetString(middleware.UserIdAuthKey)
	if userId == "" {
		responses.BadRequestResponse(ctx, "user_id not received", errors.New("user_id not received"))
		return
	}

	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(ctx.DefaultQuery("limit", "10"))

	orders, meta, err := h.orderService.GetOrders(userId, page, limit)

	if err != nil {
		responses.InternalServerResponse(ctx, "Couldnt get user orders", err)
		return
	}

	responses.PaginatedSuccessResponse(ctx, "User orders received successfully", orders, *meta)
}
