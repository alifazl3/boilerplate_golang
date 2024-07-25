package repository

import (
	auth_repository "boilerplate/internal/repository/auth"
	user_repository "boilerplate/internal/repository/user"
	"gorm.io/gorm"
)

type AuthRepository = auth_repository.AuthRepository
type UserRepository = user_repository.UserRepository

func NewAuthRepository(db *gorm.DB) AuthRepository {
	return auth_repository.NewAuthRepository(db)
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return user_repository.NewUserRepository(db)
}
