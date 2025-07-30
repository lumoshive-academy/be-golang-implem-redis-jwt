package usecase

import (
	"context"
	"errors"
	"fmt"
	"go-42/internal/data/entity"
	"go-42/internal/data/repository"
	"go-42/internal/dto"
	"go-42/pkg/caches"
	"go-42/pkg/jwt"
	"go-42/pkg/utils"

	"go.uber.org/zap"
)

type UserService interface {
	Create(ctx context.Context, user *entity.User) error
	List(ctx context.Context) (*[]entity.User, error)
	Login(ctx context.Context, user dto.LoginRequest) (*dto.ResponseUser, error)
}

type userService struct {
	Repo   repository.Repository
	Logger *zap.Logger
	Config utils.Configuration
	Cache  caches.Cacher
	Jwt    jwt.JWT
}

func NewUserService(repo repository.Repository, logger *zap.Logger, config utils.Configuration, cache caches.Cacher, jwt jwt.JWT) UserService {
	return &userService{
		Repo:   repo,
		Logger: logger,
		Config: config,
		Cache:  cache,
		Jwt:    jwt,
	}
}

func (s *userService) Create(ctx context.Context, user *entity.User) error {
	// check email
	err := s.Repo.UserRepo.GetUserByEmail(ctx, user)
	if err != nil {
		s.Logger.Error("failed to check existing user by email:", zap.Error(err))
		return err
	}

	// Create user to DB
	err = s.Repo.UserRepo.Create(user)
	if err != nil {
		s.Logger.Error("failed to create user:", zap.Error(err))
		return err
	}

	return nil
}

func (s *userService) List(ctx context.Context) (*[]entity.User, error) {
	// get user id
	userID := ctx.Value("userID")
	fmt.Println("user id ", userID)
	return s.Repo.UserRepo.List(ctx)
}

func (s *userService) Login(ctx context.Context, user dto.LoginRequest) (*dto.ResponseUser, error) {
	// check user
	// token := utils.GenerateUUIDToken()

	if user.Email == "lumos@mail.com" {
		// secret := make(map[string]interface{})
		// secret["id"] = 1
		// secret["name"] = user.Name
		// secret["email"] = user.Email

		secret := map[string]interface{}{
			"id":    1,
			"name":  "lumoshiveAcademy",
			"email": user.Email,
		}

		token, err := s.Jwt.CreateToken(user.Email, "192.168.1.45", "1")

		// dataJson, _ := json.Marshal(secret)
		// // err := s.Cache.Set("session_"+token, string(dataJson))
		if err != nil {
			// print error
			return nil, err
		}

		dataUser := dto.ResponseUser{
			Name:  secret["name"].(string),
			Email: user.Email,
			Photo: "https:/lumoshive.academy/avatar.png",
			Token: token,
		}

		return &dataUser, nil
	}

	return nil, errors.New("not found")
}
