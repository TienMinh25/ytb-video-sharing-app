package middleware

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
	"ytb-video-sharing-app-be/internal/dto"
	"ytb-video-sharing-app-be/utils"
)

type JwtAuthenticationMiddleware struct {
	KeyManager *utils.KeyManager
}

func NewJWTAuthenticationMiddleware(keyManager *utils.KeyManager) *JwtAuthenticationMiddleware {
	return &JwtAuthenticationMiddleware{
		KeyManager: keyManager,
	}
}

func JWTAuthMiddleware(params *JwtAuthenticationMiddleware) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		authHeader := ctx.GetHeader("Authorization")
		if authHeader == "" {
			utils.ErrorResponse(ctx, http.StatusUnauthorized, &dto.ErrorResponse{Message: "Missing authorization header", Code: http.StatusUnauthorized})
			ctx.Abort()
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			utils.ErrorResponse(ctx, http.StatusUnauthorized, dto.ErrorResponse{Message: "Invalid token format", Code: http.StatusUnauthorized})
			ctx.Abort()
			return
		}

		// Validate token
		claims, errResp := utils.ValidateToken(parts[1], params.KeyManager)

		if errResp != nil {
			utils.ErrorResponse(ctx, errResp.Code, dto.ErrorResponse{Message: errResp.Message, Code: errResp.Code})
			ctx.Abort()
			return
		}

		ctx.Set("claims", claims)

		ctx.Next()
	}
}

func JWTRefreshTokenMiddleware(params *JwtAuthenticationMiddleware) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		authHeader := ctx.GetHeader("X-Authorization")
		if authHeader == "" {
			utils.ErrorResponse(ctx, http.StatusUnauthorized, &dto.ErrorResponse{Message: "Missing X-Authorization header", Code: http.StatusUnauthorized})
			ctx.Abort()
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			utils.ErrorResponse(ctx, http.StatusUnauthorized, dto.ErrorResponse{Message: "Invalid token format", Code: http.StatusUnauthorized})
			ctx.Abort()
			return
		}

		// Validate token
		claims, errResp := utils.ValidateToken(parts[1], params.KeyManager)

		if errResp != nil {
			utils.ErrorResponse(ctx, errResp.Code, dto.ErrorResponse{Message: errResp.Message, Code: errResp.Code})
			ctx.Abort()
			return
		}

		ctx.Set("account_id", claims.AccountID)
		ctx.Set("refresh_token", parts[1])

		ctx.Next()
	}
}
