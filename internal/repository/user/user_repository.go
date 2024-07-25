package user_repository

import (
	"boilerplate/internal/model"
	"gorm.io/gorm"
)

type UserRepository interface {
	FindByEmail(email string) (*model.User, error)
	CreateUser(user *model.User) error
	CreateSession(session *model.Session) error
	FindByID(id int) (*model.User, error)
	UpdateUser(user *model.User) error
	DeactivateUserSessions(userID int) error
	DeactivateSession(sessionID string) error
}

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db}
}

func (r *userRepository) FindByEmail(email string) (*model.User, error) {
	var user model.User
	if err := r.db.Where("email = ?", email).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) CreateUser(user *model.User) error {
	return r.db.Create(user).Error
}

func (r *userRepository) CreateSession(session *model.Session) error {
	return r.db.Create(session).Error
}

func (r *userRepository) FindByID(id int) (*model.User, error) {
	var user model.User
	if err := r.db.Where("id = ?", id).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) UpdateUser(user *model.User) error {
	return r.db.Save(user).Error
}

func (r *userRepository) DeactivateUserSessions(userID int) error {
	return r.db.Model(&model.Session{}).Where("user_id = ?", userID).Update("status", "deleted").Error
}

func (r *userRepository) DeactivateSession(sessionID string) error {
	return r.db.Model(&model.Session{}).Where("session_id = ?", sessionID).Update("status", "deleted").Error
}
