package utils

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Response struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Data    any    `json:"data,omitempty"`
	Error   string `json:"error,omitempty"`
}

type PaginadedResponse struct {
	Response Response
	Meta     PaginationMeta `json:"meta"`
}

type PaginationMeta struct {
	Page       int   `json:"page"`
	Limit      int   `json:"limit"`
	Total      int64 `json:"total"`
	TotalPages int   `json:"total_pages"`
}

func SuccessResponse(ctx *gin.Context, message string, data any) {
	ctx.JSON(http.StatusOK, Response{
		Success: true,
		Message: message,
		Data:    data,
	})
}

func CreatedResponse(ctx *gin.Context, message string, data any) {
	ctx.JSON(http.StatusCreated, Response{
		Success: true,
		Message: message,
		Data:    data,
	})
}

func PaginatedSuccessResponse(ctx *gin.Context, message string, data any, meta PaginationMeta) {
	ctx.JSON(http.StatusOK, PaginadedResponse{
		Response: Response{
			Success: true,
			Message: message,
			Data:    data,
		},
		Meta: meta,
	})
}

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
