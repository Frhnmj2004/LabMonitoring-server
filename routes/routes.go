package routes

import (
	"github.com/Frhnmj2004/LabMonitoring-server/controllers"
	"github.com/Frhnmj2004/LabMonitoring-server/middleware"
	"github.com/Frhnmj2004/LabMonitoring-server/websocket"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
)

func SetupRoutes(app *fiber.App) {
	// Public routes
	api := app.Group("/api/v1")
	api.Post("/login", controllers.Login)

	// Protected routes
	protected := api.Use(middleware.AuthMiddleware())
	
	// Admin-only routes
	admin := protected.Use(middleware.AdminOnly())
	admin.Post("/signup", controllers.Signup)
	admin.Get("/history", controllers.GetHistory)
	admin.Get("/alerts", controllers.GetAlerts)

	// Resource monitoring routes (protected but not admin-only)
	protected.Post("/resource", controllers.PostResource)

	// WebSocket route
	app.Get("/ws/resources", websocket.New(func(c *fiber.Ctx) error {
		// IsWebSocketUpgrade returns true if the client
		// requested upgrade to the WebSocket protocol.
		if websocket.IsWebSocketUpgrade(c) {
			c.Locals("allowed", true)
			return c.Next()
		}
		return fiber.ErrUpgradeRequired
	}), websocket.HandleWebSocket)
}
