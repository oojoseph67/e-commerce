package handlers

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/oojoseph67/ecommerce/internal/middleware"
	"github.com/oojoseph67/ecommerce/internal/services"
	"github.com/oojoseph67/ecommerce/internal/utils/responses"
)

// type Handler struct {
// 	service *services.Services
// }

// func NewHandler(service *services.Services) *Handler {
// 	return &Handler{service: service}
// }

// so instead of using struct for typing we use interface to make sure the required methods are passed down correctly as well as there parameters and return types

type Handler struct {
	authService     services.AuthServicer
	userService     services.UserServicer
	productService  services.ProductServicer
	categoryService services.CategoryServicer
	uploadService   services.UploadServicer
	cartService     services.CartServicer
	orderService    services.OrderServicer
}

func NewHandler(
	auth services.AuthServicer,
	user services.UserServicer,
	product services.ProductServicer,
	category services.CategoryServicer,
	upload services.UploadServicer,
	cart services.CartServicer,
	order services.OrderServicer,
) *Handler {
	return &Handler{
		authService:     auth,
		userService:     user,
		productService:  product,
		categoryService: category,
		uploadService:   upload,
		cartService:     cart,
		orderService:    order,
	}
}

func getUserId(ctx *gin.Context) (string, bool) {
	userId := ctx.GetString(middleware.UserIdAuthKey)
	if _, err := uuid.Parse(userId); err != nil {
		responses.BadRequestResponse(ctx, "user_id not received or invalid", err)
		return "", false
	}
	return userId, true
}

func getReqParam(ctx *gin.Context, param string) (string, bool) {
	val := ctx.Param(param)

	if val == "" || val == ":"+param {
		message := fmt.Sprintf("%s not received or invalid", param)
		responses.BadRequestResponse(ctx, message, errors.New(message))
		return "", false
	}

	if _, err := uuid.Parse(val); err != nil {
		message := fmt.Sprintf("%s not received or invalid", param)
		responses.BadRequestResponse(ctx, message, errors.New(message))
		return "", false
	}

	return val, true
}

func getPaginationValues(ctx *gin.Context) (page, limit int) {
	p, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	li, _ := strconv.Atoi(ctx.DefaultQuery("limit", "10"))

	if p < 1 {
		p = 1
	}
	if li < 1 || li > 100 {
		li = 10
	}

	return p, li
}
