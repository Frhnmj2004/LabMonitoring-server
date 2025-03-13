package main

import (
	"log"
	"os"

	"github.com/Frhnmj2004/LabMonitoring-server/config"
	"github.com/Frhnmj2004/LabMonitoring-server/routes"
	"github.com/Frhnmj2004/LabMonitoring-server/utils"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/joho/godotenv"
)

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}

	// Initialize database connection
	config.InitDB()

	// Create Fiber app
	app := fiber.New(fiber.Config{
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			utils.LogError("Unhandled error: %v", err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Internal server error",
			})
		},
	})

	// Middleware
	app.Use(cors.New())
	app.Use(logger.New())
	app.Use(recover.New())

	// Setup routes
	routes.SetupRoutes(app)

	// Get port from environment variable
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	utils.LogInfo("Server starting on port %s", port)
	if err := app.Listen(":" + port); err != nil {
		utils.LogError("Server failed to start: %v", err)
		os.Exit(1)
	}
}
