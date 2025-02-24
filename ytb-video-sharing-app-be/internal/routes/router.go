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
	videoHandler *handler.VideoHandler,
	middleware *middleware.JwtAuthenticationMiddleware,
) *Router {
	apiV1Group := router.Group("/api/v1")

	registerAccountEndpoint(accountHandler, apiV1Group, middleware)
	registerVideoEndpoint(videoHandler, apiV1Group, middleware)

	return &Router{
		Router: router,
	}
}

func registerAccountEndpoint(accountHandler *handler.AccountHandler, group *gin.RouterGroup, params *middleware.JwtAuthenticationMiddleware) {
	accountGroup := group.Group("/accounts")

	accountGroup.POST("/register", middleware.ValidateRequest[dto.CreateAccountRequest](), accountHandler.Register)
	accountGroup.POST("/login", middleware.ValidateRequest[dto.LoginRequest](), accountHandler.Login)
	accountGroup.POST("/logout/:accountID", middleware.JWTAuthMiddleware(params), accountHandler.Logout)
	accountGroup.POST("/refresh-token", middleware.JWTRefreshTokenMiddleware(params), accountHandler.RefreshToken)
}

func registerVideoEndpoint(videoHandler *handler.VideoHandler, group *gin.RouterGroup, params *middleware.JwtAuthenticationMiddleware) {
	videoGroup := group.Group("/videos")

	videoGroup.POST("", middleware.JWTAuthMiddleware(params), middleware.ValidateRequest[dto.ShareVideoRequest](), videoHandler.ShareVideo)
	videoGroup.GET("", videoHandler.GetListVideos)
}
