package main

import (
	"flag"
	"log"

	"bazar-be/config"
)

func main() {
	freshFlag := flag.Bool("fresh", false, "Reset and truncate all database tables before seeding")
	flag.Parse()

	log.Println("Initializing database connection for seeder...")
	config.ConnectDB()

	if err := config.SeedAll(config.DB, *freshFlag); err != nil {
		log.Fatalf("Database seeding failed: %v", err)
	}

	log.Println("Seeding process finished successfully!")
}
