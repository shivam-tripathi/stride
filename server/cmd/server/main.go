package main

import (
	"context"
	"fmt"
	"log"

	"quizizz.com/internal/config"
	"quizizz.com/internal/resources"
	"quizizz.com/wire"
)

//go:generate go run github.com/google/wire/cmd/wire

func main() {
	// Initialize configuration
	cfg := config.NewConfig()

	// Create resources (not yet connected)
	db := resources.NewDB(cfg)
	redis := resources.NewRedis(cfg)
	res := &resources.Resources{
		DB:    db,
		Redis: redis,
	}

	// Initialize resources BEFORE creating the app
	// This ensures resources are connected when repositories are created
	fmt.Println("Initializing resources...")
	ctx := context.Background()
	if err := resources.InitResources(ctx, res); err != nil {
		log.Fatalf("Failed to initialize resources: %v", err)
	}

	// Now initialize the application with connected resources
	app, err := wire.InitializeAppWithResources(cfg, res)
	if err != nil {
		log.Fatalf("Failed to initialize application: %v", err)
	}

	// Start the server
	fmt.Println("Starting server...")
	if err := app.Run(); err != nil {
		log.Fatalf("Failed to run application: %v", err)
	}
}
