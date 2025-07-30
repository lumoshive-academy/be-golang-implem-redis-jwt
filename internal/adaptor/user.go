package adaptor

import (
	"encoding/json"
	"go-42/internal/data/entity"
	"go-42/internal/dto"
	"go-42/internal/usecase"
	"go-42/pkg/response"
	"go-42/pkg/utils"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type HandlerUser struct {
	User   usecase.UserService
	Logger *zap.Logger
}

func NewHandlerUser(user usecase.UserService, logger *zap.Logger) HandlerUser {
	return HandlerUser{
		User:   user,
		Logger: logger,
	}
}

func (h *HandlerUser) Register(ctx *gin.Context) {
	user := entity.User{}

	// Decode JSON
	if err := json.NewDecoder(ctx.Request.Body).Decode(&user); err != nil {
		response.ResponseBadRequest(ctx, http.StatusBadRequest, "Invalid request payload")
		return
	}

	// Validate data
	if validationErrors, err := utils.ValidateData(user); err != nil {
		h.Logger.Error("validate error", zap.Error(err))
		response.ResponseBadRequest2(ctx, http.StatusBadRequest, validationErrors)
		return
	}

	// save to DB
	err := h.User.Create(ctx.Request.Context(), &user)
	if err != nil {
		response.ResponseBadRequest(ctx, http.StatusBadRequest, "Register failed")
		return
	}

	response.ResponseSuccess(ctx, http.StatusCreated, "success register", nil)
}

func (h *HandlerUser) ListUser(ctx *gin.Context) {
	users, err := h.User.List(ctx.Request.Context())
	if err != nil {
		response.ResponseBadRequest(ctx, http.StatusBadRequest, "not found")
		return
	}
	response.ResponseSuccess(ctx, http.StatusOK, "success", users)
}

func (h *HandlerUser) Login(ctx *gin.Context) {
	var User dto.LoginRequest
	if err := ctx.ShouldBindJSON(&User); err != nil {
		response.ResponseBadRequest(ctx, http.StatusBadRequest, "invalid body")
		return
	}

	data, err := h.User.Login(ctx.Request.Context(), User)
	if err != nil {
		response.ResponseBadRequest(ctx, http.StatusBadRequest, err.Error())
		return
	}
	response.ResponseSuccess(ctx, http.StatusOK, "success login", data)
}

func (h *HandlerUser) Profile(ctx *gin.Context) {

}
