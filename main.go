package main

import (
	"log"

	"inventory-juanfe/config"
	"inventory-juanfe/routers"
	"inventory-juanfe/utils"
)

func main() {

	// Load database configuration and connect
	dbConfig := config.LoadDBConfig()
	db := config.ConnectDB(dbConfig)
	defer db.Close()

	// Create Fiber app
	app, err := utils.InitFiber()
	if err != nil {
		log.Fatalf("Error initializing Fiber: %v", err)
	}

	// Setup routes
	routers.SetupRoutes(app, db)

	// Start server
	if err := app.Listen(":8080"); err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
}
