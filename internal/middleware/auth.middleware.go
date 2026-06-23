package middleware

import (
	"errors"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/oojoseph67/ecommerce/internal/models"
	jwtp "github.com/oojoseph67/ecommerce/internal/utils/jwt"
	"github.com/oojoseph67/ecommerce/internal/utils/responses"
)

const (
	UserIdAuthKey    string = "user_id"
	UserEmailAuthKey string = "user_email"
	UserRoleAuthKey  string = "user_role"
)

// Auth validates JWT tokens and sets user context.
func Auth(secret string) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		authHeader := ctx.GetHeader("Authorization")
		if authHeader == "" {
			responses.UnauthorizedResponse(ctx, "Authorization header is required", errors.New("missing authorization header"))
			ctx.Abort()
			return
		}

		tokenParts := strings.Split(authHeader, " ")
		if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
			responses.UnauthorizedResponse(ctx, "Invalid authorization header format", errors.New("invalid authorization header format"))
			ctx.Abort()
			return
		}

		claims, err := jwtp.ValidateToken(tokenParts[1], secret)
		if err != nil {
			responses.UnauthorizedResponse(ctx, "Invalid token", errors.New("invalid token"))
			ctx.Abort()
			return
		}

		ctx.Set(UserIdAuthKey, claims.UserID)
		ctx.Set(UserEmailAuthKey, claims.Email)
		ctx.Set(UserRoleAuthKey, claims.Role)

		ctx.Next()
	}
}

// Admin restricts access to admin users only.
func Admin() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		role, exists := ctx.Get("user_role")
		if !exists {
			responses.ForbiddenResponse(ctx, "Forbidden", errors.New("missing user role"))
			ctx.Abort()
			return
		}

		if role != string(models.UserRoleAdmin) {
			responses.ForbiddenResponse(ctx, "Forbidden", errors.New("admin access required"))
			ctx.Abort()
			return
		}

		ctx.Next()
	}
}
