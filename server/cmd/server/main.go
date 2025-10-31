package main

import (
	"fmt"
	"log"

	"quizizz.com/wire"
)

//go:generate go run github.com/google/wire/cmd/wire

func main() {
	// Initialize the application
	app, err := wire.InitializeApp()
	if err != nil {
		log.Fatalf("Failed to initialize application: %v", err)
	}

	// Start the server
	fmt.Println("Starting server...")
	if err := app.Run(); err != nil {
		log.Fatalf("Failed to run application: %v", err)
	}
}
