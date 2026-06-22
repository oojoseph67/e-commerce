package server

import (
	"errors"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/oojoseph67/ecommerce/internal/models"
	jwtp "github.com/oojoseph67/ecommerce/internal/utils/jwt"
	"github.com/oojoseph67/ecommerce/internal/utils/responses"
)

func (s *Server) corsMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.Header("Access-Control-Allow-Origin", "*")
		ctx.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		ctx.Header("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if ctx.Request.Method == "OPTIONS" {
			ctx.AbortWithStatus(http.StatusNoContent)
		}

		ctx.Next()
	}
}

func (s *Server) authMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		authHeader := ctx.GetHeader("Authorization")

		if authHeader == "" {
			responses.UnauthorizedResponse(ctx, "Authorization header is required", errors.New("authorization header is required"))
			ctx.Abort()
			return
		}

		tokenParts := strings.Split(authHeader, " ")
		if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
			responses.UnauthorizedResponse(ctx, "Invalid authorization header format", errors.New("invalid authorization header format"))
			ctx.Abort()
			return
		}

		jwt := tokenParts[1]
		claims, err := jwtp.ValidateToken(jwt, s.config.JWT.Secret)
		if err != nil {
			responses.UnauthorizedResponse(ctx, "Invalid token", errors.New("invalid token"))
			ctx.Abort()
			return
		}

		ctx.Set("user_id", claims.UserID)
		ctx.Set("user_email", claims.Email)
		ctx.Set("user_role", claims.Role)

		ctx.Next()
	}
}

func (s *Server) adminMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		role, exists := ctx.Get("user_role")
		if !exists {
			responses.ForbiddenResponse(ctx, "Forbidden", errors.New("forbidden"))
			ctx.Abort()
			return
		}

		if role != string(models.UserRoleAdmin) {
			responses.ForbiddenResponse(ctx, "Forbidden", errors.New("Forbidden"))
			ctx.Abort()
			return
		}

		ctx.Next()

	}
}
