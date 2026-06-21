package utils

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func ErrorResponse(ctx *gin.Context, statusCode int, message string, err error) {
	response := Response{
		Success: false,
		Message: message,
	}

	if err != nil {
		response.Error = err.Error()
	}

	ctx.JSON(statusCode, response)
}

func BadRequestResponse(ctx *gin.Context, message string, err error) {
	ErrorResponse(ctx, http.StatusBadRequest, message, err)
}

func UnauthorizedResponse(ctx *gin.Context, message string, err error) {
	ErrorResponse(ctx, http.StatusUnauthorized, message, err)
}

func ForbiddenResponse(ctx *gin.Context, message string, err error) {
	ErrorResponse(ctx, http.StatusForbidden, message, err)
}

func NotFoundResponse(ctx *gin.Context, message string, err error) {
	ErrorResponse(ctx, http.StatusNotFound, message, err)
}

func InternalServerResponse(ctx *gin.Context, message string, err error) {
	ErrorResponse(ctx, http.StatusInternalServerError, message, err)
}
