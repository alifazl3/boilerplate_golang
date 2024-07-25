package model

import (
	"gorm.io/gorm"
	"time"
)

type User struct {
	ID               int       `gorm:"primaryKey"`
	UID              string    `gorm:"unique"`
	UserName         string    `gorm:"type:varchar(100)"`
	Email            *string   `gorm:"type:varchar(255)"`
	ValidEmail       *string   `gorm:"type:varchar(100)"`
	PhoneNumber      *string   `gorm:"type:varchar(100)"`
	ProfileIMG       *string   `gorm:"type:varchar(100)"`
	Age              *int      `gorm:"type:int(11)"`
	AgeRange         *string   `gorm:"type:enum('0-10', '10-18' , '18-21', '21-50' , '50-100')"`
	ValidPhoneNumber *string   `gorm:"type:varchar(100)"`
	Password         string    `gorm:"type:varchar(255);not null"`
	PasswordSalt     string    `gorm:"type:varchar(255);not null"`
	Country          *string   `gorm:"type:varchar(100)"`
	FirstName        *string   `gorm:"type:varchar(100)"`
	LastName         *string   `gorm:"type:varchar(100)"`
	level            string    `gorm:"type:enum('user', 'owner' , 'waiters' , 'us');default:'user';not null"`
	CreatedAt        time.Time `gorm:"autoCreateTime;not null"`
	UpdatedAt        time.Time
	DeletedAt        gorm.DeletedAt `gorm:"index"`
}
