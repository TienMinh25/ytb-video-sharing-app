package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
	"ytb-video-sharing-app-be/db"
	_ "ytb-video-sharing-app-be/docs"
	"ytb-video-sharing-app-be/internal/handler"
	"ytb-video-sharing-app-be/internal/middleware"
	"ytb-video-sharing-app-be/internal/migrate"
	"ytb-video-sharing-app-be/internal/repository"
	"ytb-video-sharing-app-be/internal/routes"
	"ytb-video-sharing-app-be/internal/service"
	"ytb-video-sharing-app-be/internal/websock"
	"ytb-video-sharing-app-be/utils"

	"github.com/gin-contrib/cors"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"go.uber.org/fx"
)

func LoadEnv() error {
	err := godotenv.Load("configs/config.dev.env")

	if err != nil {
		return err
	}

	fmt.Println("✅ Load .env file successfully!")

	return nil
}

func MigrateDB() error {
	err := migrate.Migrate(os.Getenv("MIGRATION_DIR"))

	if err != nil {
		return err
	}

	fmt.Println("✅ Mirgate database successfully!")
	return nil
}

func NewGinEngine() *gin.Engine {
	router := gin.Default()

	router.Use(cors.New(cors.Config{
		AllowMethods:     []string{http.MethodGet, http.MethodHead, http.MethodPost, http.MethodPut, http.MethodPatch, http.MethodDelete, http.MethodConnect, http.MethodOptions, http.MethodTrace},
		AllowHeaders:     []string{"*"},
		AllowCredentials: true,
		AllowOrigins:     []string{"*"},
		MaxAge:           12 * time.Hour,
	}))

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	return router
}

func StartServer(lifecycle fx.Lifecycle, r *routes.Router) {
	lifecycle.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			address := os.Getenv("SERVER_ADDRESS")

			go func() {
				if err := r.Router.Run(address); err != nil {
					log.Fatal(err)
				}

				log.Println("Server is running on " + address)
			}()

			return nil
		},
		OnStop: func(ctx context.Context) error {
			log.Println("Shut down!")

			return nil
		},
	})
}

func NewRetention() websock.RetentionMap {
	return websock.NewRetentionMap(context.Background(), 1*time.Minute)
}

func StartWebSocketServer(lifecycle fx.Lifecycle, wsMux *http.ServeMux) {
	lifecycle.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			address := os.Getenv("WEBSOCKET_SERVER_ADDRESS")

			go func() {
				fmt.Println("✅ WebSocket Server is running on", address)
				if err := http.ListenAndServe(address, wsMux); err != nil {
					log.Fatal(err)
				}
			}()
			return nil
		},
		OnStop: func(ctx context.Context) error {
			fmt.Println("Shutting down WebSocket server...")

			// if err := q.Close(); err != nil {
			// 	log.Println("Error shutting down queue:", err)
			// } else {
			// 	log.Println("✅ Queue shutdown successfully!")
			// }

			return nil
		},
	})
}

// @title						YouTube Video Sharing API
// @version					1.0
// @description				API cho ứng dụng chia sẻ video YouTube
// @host						localhost:3000
// @BasePath					/api/v1
//
// @securityDefinitions.apikey	BearerAuth
// @in							header
// @name						Authorization
func main() {
	app := fx.New(
		fx.Provide(
			db.NewMySQL,
			repository.NewAccountRepository,
			repository.NewAccountPasswordRepository,
			repository.NewRefreshTokenRepository,
			repository.NewVideoRepository,
			service.NewAccountService,
			service.NewVideoService,
			handler.NewAccountHandler,
			handler.NewVideoHandler,
			NewGinEngine,
			routes.NewRouter,
			utils.LoadKeys,
			middleware.NewJWTAuthenticationMiddleware,
			websock.NewManager,
			NewRetention,
			// third_party.NewQueue,
		),
		fx.Invoke(LoadEnv, MigrateDB, utils.LoadKeys),
		fx.Invoke(StartServer, StartWebSocketServer),
	)

	app.Run()
}
