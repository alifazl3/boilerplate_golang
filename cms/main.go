package main

import (
	"boilerplate/internal/config"
	"boilerplate/internal/database"
	"boilerplate/router"
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	// Load Config
	conf, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Error loading config: %v", err)
	}

	// Connect to database
	fmt.Println("Connecting to database...")
	db, dbErr := database.ConnectToDatabase(conf)
	if dbErr != nil {
		log.Fatalf("Failed to connect to database: %v", dbErr)
	}

	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	r := &router.Router{
		DB: db,
	}

	// Start server
	go func() {
		if err := r.Routes(conf); err != nil {
			log.Fatalf("Failed to run server: %v", err)
		}
	}()

	<-sigs
	cancel()

	// Close database connection
	sqlDB, err := db.DB()
	if err != nil {
		log.Println("Failed to get underlying DB:", err)
		return
	}
	if err := sqlDB.Close(); err != nil {
		log.Println("Failed to close database connection:", err)
	}
}
