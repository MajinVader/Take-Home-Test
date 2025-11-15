package main

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"

	"take-home-test/database"
	"take-home-test/models"
	"take-home-test/routes"
)

func main() {
	// Load .env
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}

	// Connect DB
	database.ConnectDB()

	// Auto-migrate semua model
	if err := database.DB.AutoMigrate(&models.User{}, &models.Field{}, &models.Booking{}); err != nil {
		log.Fatal("failed to migrate:", err)
	}

	app := fiber.New()

	// Daftarin semua routes di satu tempat
	routes.SetupRoutes(app)

	log.Println("Server running on :3000")
	if err := app.Listen(":3000"); err != nil {
		log.Fatal(err)
	}
}
