package main

import (
	"ytb-video-sharing-app-be/internal/migrate"

	"github.com/joho/godotenv"
	"go.uber.org/fx"
)

func LoadEnv() error {
	err := godotenv.Load("configs/config.dev.env")

	if err != nil {
		return err
	}

	return nil
}

func MigrateDB() error {
	err := migrate.Migrate("file://db/migrations")

	if err != nil {
		return err
	}

	return nil
}

func main() {
	app := fx.New(
		fx.Invoke(LoadEnv, MigrateDB),
	)

	app.Run()
}
