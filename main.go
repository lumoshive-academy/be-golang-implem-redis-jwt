package main

import (
	"go-42/cmd"
	"go-42/internal/data"
	"go-42/internal/data/repository"
	"go-42/internal/wire"
	"go-42/pkg/caches"
	"go-42/pkg/database"
	"go-42/pkg/jwt"
	"go-42/pkg/middleware"
	"go-42/pkg/utils"
	"log"
	"os"

	"go.uber.org/zap"
)

func main() {
	// read config
	config, err := utils.ReadConfiguration()
	if err != nil {
		log.Fatal(err)
	}

	// init logger
	logger, err := utils.InitLogger(config.PathLogger, config)
	if err != nil {
		log.Fatal("can't init logger %w", zap.Error(err))
	}

	//Init db
	db, err := database.InitDB(config)
	if err != nil {
		logger.Fatal("can't connect to database ", zap.Error(err))
	}

	//ini redis (cache)
	rdb := caches.NewCacher(config, 60*60)

	// migration
	if err := data.AutoMigrate(db); err != nil {
		logger.Fatal("failed to run migrations", zap.Error(err))
	}

	// seeder
	if err := data.SeedAll(db); err != nil {
		logger.Fatal("failed to seed initial data", zap.Error(err))
	}

	pemPrivate, err := os.ReadFile("private.pem")
	if err != nil {
		logger.Fatal("Failed to read private key:", zap.Error(err))
	}

	pemPublic, err := os.ReadFile("public.pem")
	if err != nil {
		logger.Fatal("Failed to read public key:", zap.Error(err))
	}

	// jwt
	jwt := jwt.NewJWT(pemPrivate, pemPublic, logger)

	repo := repository.NewRepository(db, logger, rdb)
	mLogger := middleware.NewLoggerMiddleware(logger)
	mAuth := middleware.NewAuthMiddleware(logger, rdb)
	router := wire.Wiring(repo, mLogger, mAuth, logger, config, rdb, jwt)

	cmd.ApiServer(config, logger, router)
}
