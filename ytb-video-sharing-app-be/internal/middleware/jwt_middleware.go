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
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, &dto.ErrorResponse{Message: "missing authorization header", Code: http.StatusUnauthorized})
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, dto.ErrorResponse{Message: "invalid token format", Code: http.StatusUnauthorized})
			return
		}

		// Validate token
		claims, errResp := utils.ValidateToken(parts[1], params.KeyManager)

		if errResp != nil {
			c.AbortWithStatusJSON(errResp.Code, errResp)
			return
		}

		// Lưu thông tin user vào context
		c.Set("claims", claims)

		c.Next()
	}
}
