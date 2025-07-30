package repository

import (
	"go-42/pkg/caches"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

type Repository struct {
	UserRepo UserRepository
}

func NewRepository(db *gorm.DB, log *zap.Logger, cache caches.Cacher) Repository {
	return Repository{
		UserRepo: NewUserRepository(db, log, cache),
	}
}
