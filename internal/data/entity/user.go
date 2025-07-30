package entity

import "go-42/pkg/utils"

type User struct {
	Model
	Name     string `gorm:"type:varchar(100);not null" json:"name" binding:"required"`
	Email    string `gorm:"type:varchar(100);uniqueIndex;not null" json:"email" binding:"required,email"`
	Password string `gorm:"type:varchar(100);not null" json:"password" binding:"required,min=6"`
	Photo    string `gorm:"type:text;not null" json:"photo" binding:"required"`
}

func SeedUsers() []User {
	users := []User{
		{
			Name:     "Budi Santoso",
			Email:    "budi@example.com",
			Password: utils.HashPassword("password123"),
			Photo:    "https://example.com/photos/budi.jpg",
		},
		{
			Name:     "Siti Aminah",
			Email:    "siti@example.com",
			Password: utils.HashPassword("rahasia123"),
			Photo:    "https://example.com/photos/siti.jpg",
		},
		{
			Name:     "Andi Saputra",
			Email:    "andi@example.com",
			Password: utils.HashPassword("12345678"),
			Photo:    "https://example.com/photos/andi.jpg",
		},
	}

	return users
}
