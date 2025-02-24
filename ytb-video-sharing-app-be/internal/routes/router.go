package routes

import (
	"github.com/gin-gonic/gin"
	"ytb-video-sharing-app-be/internal/dto"
	"ytb-video-sharing-app-be/internal/handler"
	"ytb-video-sharing-app-be/internal/middleware"
)

type Router struct {
	Router *gin.Engine
}

func NewRouter(
	router *gin.Engine,
	accountHandler *handler.AccountHandler,
	middleware *middleware.JwtAuthenticationMiddleware,
) *Router {
	apiV1Group := router.Group("/api/v1")

	registerAccountEndpoint(accountHandler, apiV1Group, middleware)

	return &Router{
		Router: router,
	}
}

func registerAccountEndpoint(accountHandler *handler.AccountHandler, group *gin.RouterGroup, params *middleware.JwtAuthenticationMiddleware) {
	accountGroup := group.Group("/accounts")

	accountGroup.POST("/register", middleware.ValidateRequest[dto.CreateAccountRequest](), accountHandler.Register)
	accountGroup.POST("/login", middleware.ValidateRequest[dto.LoginRequest](), accountHandler.Login)
	// TODO: add middleware handle access token
	accountGroup.POST("/logout/:accountID", middleware.JWTAuthMiddleware(params), accountHandler.Logout)
	accountGroup.POST("/refresh-token/:accountID", accountHandler.RefreshToken)
}
