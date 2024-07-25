package auth_repository

import (
	"boilerplate/internal/model"
	"time"

	"gorm.io/gorm"
)

type AuthRepository interface {
	GetUserTokenBySessionId(sessionID string) (*model.Session, error)
}

type authRepository struct {
	db *gorm.DB
}

func NewAuthRepository(db *gorm.DB) AuthRepository {
	return &authRepository{db}
}

func (r *authRepository) GetUserTokenBySessionId(sessionID string) (*model.Session, error) {
	var session model.Session
	if err := r.db.Where("session_id = ? AND status = ? AND token_expiration > ?", sessionID, "active", time.Now()).First(&session).Error; err != nil {
		return nil, err
	}
	return &session, nil
}
