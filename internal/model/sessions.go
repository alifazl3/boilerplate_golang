package model

import (
	"gorm.io/gorm"
	"time"
)

type Session struct {
	ID              int    `gorm:"primaryKey"`
	UserId          int    `gorm:"foreignKey;not null"`
	SessionId       string `gorm:"unique"`
	UserToken       string `gorm:"type:varchar(255)"`
	TokenExpiration time.Time
	TokenFacebook   *string   `gorm:"type:varchar(255)"`
	TokenGoogle     *string   `gorm:"type:varchar(255)"`
	Status          string    `gorm:"type:enum('active', 'deactivate' , 'deleted');default:'active';not null"`
	CreatedAt       time.Time `gorm:"autoCreateTime;not null"`
	UpdatedAt       time.Time
	DeletedAt       gorm.DeletedAt `gorm:"index"`
}
