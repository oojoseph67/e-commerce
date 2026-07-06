package handlers

import (
	"github.com/gin-gonic/gin"
	_ "github.com/oojoseph67/ecommerce/internal/dto"
	"github.com/oojoseph67/ecommerce/internal/utils/responses"
)

// CreateOrder godoc
// @Summary Create an order
// @Description Create a new order from the user's cart items
// @Tags Orders
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 201 {object} responses.Response{data=dto.OrderResponse}
// @Failure 401 {object} responses.Response
// @Failure 500 {object} responses.Response
// @Router /order/ [post]
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

// GetOrder godoc
// @Summary Get an order
// @Description Retrieve a single order by ID
// @Tags Orders
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Order ID"
// @Success 200 {object} responses.Response{data=dto.OrderResponse}
// @Failure 401 {object} responses.Response
// @Failure 404 {object} responses.Response
// @Router /order/{id} [get]
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

// GetOrders godoc
// @Summary Get all orders
// @Description Retrieve a paginated list of the user's orders
// @Tags Orders
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Items per page" default(10)
// @Success 200 {object} responses.PaginatedResponse
// @Failure 401 {object} responses.Response
// @Failure 500 {object} responses.Response
// @Router /order/ [get]
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
