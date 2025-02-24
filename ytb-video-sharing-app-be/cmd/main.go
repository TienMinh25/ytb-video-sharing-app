package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"ytb-video-sharing-app-be/db"
	_ "ytb-video-sharing-app-be/docs"
	"ytb-video-sharing-app-be/internal/handler"
	"ytb-video-sharing-app-be/internal/middleware"
	"ytb-video-sharing-app-be/internal/migrate"
	"ytb-video-sharing-app-be/internal/repository"
	"ytb-video-sharing-app-be/internal/routes"
	"ytb-video-sharing-app-be/internal/service"
	"ytb-video-sharing-app-be/utils"

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
		),
		fx.Invoke(LoadEnv, MigrateDB, utils.LoadKeys),
		fx.Invoke(StartServer),
	)

	app.Run()
}
