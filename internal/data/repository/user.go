package repository

import (
	"context"
	"errors"
	"go-42/internal/data/entity"
	"go-42/pkg/caches"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

type UserRepository interface {
	Create(user *entity.User) error
	GetUserByEmail(ctx context.Context, user *entity.User) error
	List(ctx context.Context) (*[]entity.User, error)
}

type userRepositoryImpl struct {
	DB     *gorm.DB
	Log    *zap.Logger
	Cacher caches.Cacher
}

func NewUserRepository(db *gorm.DB, log *zap.Logger, caches caches.Cacher) UserRepository {
	return &userRepositoryImpl{
		DB:     db,
		Log:    log,
		Cacher: caches,
	}
}

func (r *userRepositoryImpl) Create(user *entity.User) error {
	// query := `
	// 	INSERT INTO users (name, email, password, photo, created_at, updated_at)
	// 	VALUES ($1, $2, $3, $4, NOW(), NOW())
	// 	RETURNING id, created_at, updated_at
	// `
	// result := r.DB.Raw(
	// 	query,
	// 	user.Name,
	// 	user.Email,
	// 	user.Password,
	// 	user.Photo,
	// ).Scan(&user)
	// if result.Error != nil {
	// 	return errors.New("error insert data")
	// }

	// return nil

	result := r.DB.Create(user)
	if result.Error != nil {
		return errors.New("error insert data")
	}

	return nil

}

func (r *userRepositoryImpl) GetUserByEmail(ctx context.Context, user *entity.User) error {
	result := r.DB.WithContext(ctx).Where("email = ?", user.Email).First(user)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil
		}
		return result.Error
	}
	return nil
}

func (r *userRepositoryImpl) List(ctx context.Context) (*[]entity.User, error) {
	var users []entity.User
	if err := r.DB.WithContext(ctx).Limit(2).Find(&users).Error; err != nil {
		return nil, err
	}
	return &users, nil

}
