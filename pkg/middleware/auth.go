package middleware

import (
	"fmt"
	"go-42/pkg/caches"
	"go-42/pkg/response"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type AuthMiddleware struct {
	Logger *zap.Logger
	Cache  caches.Cacher
}

func NewAuthMiddleware(logger *zap.Logger, cache caches.Cacher) AuthMiddleware {
	return AuthMiddleware{logger, cache}
}

func (m *AuthMiddleware) Auth() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		token := ctx.Request.Header.Get("token")
		data, err := m.Cache.Get("session_" + token)
		if err != nil {
			response.ResponseBadRequest(ctx, http.StatusUnauthorized, "token invalid")
			return
		}
		fmt.Println(data)
		// ctx.Set("userid", 123)
		ctx.Next()
	}
}
