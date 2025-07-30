package wire

import (
	"go-42/internal/adaptor"
	"go-42/internal/data/repository"
	"go-42/internal/usecase"
	"go-42/pkg/caches"
	"go-42/pkg/jwt"
	"go-42/pkg/middleware"
	"go-42/pkg/utils"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func Wiring(repo repository.Repository, mLogger middleware.LoggerMiddleware, middlwareAuth middleware.AuthMiddleware, logger *zap.Logger, config utils.Configuration, caches caches.Cacher, jwt jwt.JWT) *gin.Engine {
	router := gin.New()
	router.Use(
		mLogger.LoggingMiddleware(),
	)

	api := router.Group("/api/v1")
	wireUser(api, middlwareAuth, repo, logger, config, caches, jwt)
	return router
}

func wireUser(router *gin.RouterGroup, middlwareAuth middleware.AuthMiddleware, repo repository.Repository, logger *zap.Logger, config utils.Configuration, cache caches.Cacher, jwt jwt.JWT) {
	usecaseUser := usecase.NewUserService(repo, logger, config, cache, jwt)
	adaptorUser := adaptor.NewHandlerUser(usecaseUser, logger)
	router.POST("/register", adaptorUser.Register)
	router.GET("/users", jwt.AuthJWT(), adaptorUser.ListUser)
	router.GET("/profile", middlwareAuth.Auth(), adaptorUser.Profile)
	router.POST("/login", adaptorUser.Login)
}
