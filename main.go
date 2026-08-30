package main

import (
	"log"
	"os"

	"bazar-be/config"
	"bazar-be/routes"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/cors"
	"github.com/gofiber/fiber/v3/middleware/logger"
	"github.com/gofiber/fiber/v3/middleware/recover"
)

func main() {
	// Connect to Database
	config.ConnectDB()

	// Check CLI flags for seeder execution
	args := os.Args[1:]
	for _, arg := range args {
		if arg == "--seed" || arg == "-seed" {
			log.Println("CLI Flag detected: Seeding database...")
			if err := config.SeedAll(config.DB, false); err != nil {
				log.Fatalf("Seeding failed: %v", err)
			}
			break
		} else if arg == "--seed-fresh" || arg == "-seed-fresh" || arg == "--fresh" {
			log.Println("CLI Flag detected: Fresh seeding database...")
			if err := config.SeedAll(config.DB, true); err != nil {
				log.Fatalf("Fresh seeding failed: %v", err)
			}
			break
		}
	}

	// Initialize Fiber App
	app := fiber.New(fiber.Config{
		AppName: "Bazar Backend API v1.0",
	})

	// Middlewares
	app.Use(recover.New())
	app.Use(logger.New())
	app.Use(cors.New(cors.Config{
		AllowOrigins: []string{"*"},
		AllowHeaders: []string{"Origin", "Content-Type", "Accept", "Authorization"},
		AllowMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
	}))

	// Setup Routes
	routes.SetupRoutes(app)

	// Start Server
	port := os.Getenv("APP_PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Server starting on port %s...\n", port)
	if err := app.Listen(":" + port); err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
}
