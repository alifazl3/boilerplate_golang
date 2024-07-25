package migration

import (
	"boilerplate/internal/model"
	"fmt"

	"gorm.io/gorm"
)

func Migrate(db *gorm.DB) error {
	err := Up(db)
	if err != nil {
		return fmt.Errorf("failed to up migration: %w", err)
	}
	return nil
}

func PerformMigrations(db *gorm.DB, models ...interface{}) error {
	for _, m := range models {
		err := db.AutoMigrate(m)
		if err != nil {
			return err
		}
	}
	return nil
}

// Up performs migration up
func Up(db *gorm.DB) error {
	m := []interface{}{
		&model.User{},
		&model.Session{},
	}
	return PerformMigrations(db, m...)
}
