package database

import (
	"boilerplate/internal/config"
	"boilerplate/internal/database/migration"
	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql" // Import the MySQL driver
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func ConnectToDatabase(conf *config.Config) (*gorm.DB, error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		conf.DBUser, conf.DBPassword, conf.DBHost, conf.DBPort, conf.DBName)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}
	// Perform migration
	err = migration.Migrate(db)
	if err != nil {
		log.Fatalf("failed to migrate database: %v", err)
		return nil, fmt.Errorf("failed to migrate database: %w", err)
	}

	return db, nil
}
