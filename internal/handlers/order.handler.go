package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/oojoseph67/ecommerce/internal/utils/responses"
)

func (h *Handler) CreateOrder(ctx *gin.Context) {
	userId, ok := getUserId(ctx)
	if !ok {
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
	id, ok := getReqParam(ctx, "id")
	if !ok {
		return
	}

	userId, ok := getUserId(ctx)
	if !ok {
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
	userId, ok := getUserId(ctx)
	if !ok {
		return
	}

	page, limit := getPaginationValues(ctx)

	orders, meta, err := h.orderService.GetOrders(userId, page, limit)
	if err != nil {
		responses.InternalServerResponse(ctx, "Couldnt get user orders", err)
		return
	}

	responses.PaginatedSuccessResponse(ctx, "User orders received successfully", orders, *meta)
}
