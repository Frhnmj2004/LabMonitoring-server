package routes

import (
	"github.com/Frhnmj2004/LabMonitoring-server/controllers"
	"github.com/Frhnmj2004/LabMonitoring-server/middleware"
	"github.com/Frhnmj2004/LabMonitoring-server/websocket"
	"github.com/gofiber/fiber/v2"
	fiberwebsocket "github.com/gofiber/websocket/v2"
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
	
	// Resource routes
	protected.Post("/resource", controllers.PostResource)
	admin.Get("/resources/history", controllers.GetHistory)
	
	// Alert routes (admin only)
	alertGroup := admin.Group("/alerts")
	alertGroup.Get("/", controllers.GetAlerts)          // Get all alerts
	alertGroup.Get("/active", controllers.GetActiveAlerts)
	alertGroup.Get("/history", controllers.GetAlertHistory)
	alertGroup.Get("/stats", controllers.GetAlertStats)
	alertGroup.Put("/:id/resolve", controllers.ResolveAlert)

	// WebSocket setup
	app.Use("/ws", func(c *fiber.Ctx) error {
		if fiberwebsocket.IsWebSocketUpgrade(c) {
			c.Locals("allowed", true)
			return c.Next()
		}
		return fiber.ErrUpgradeRequired
	})

	// WebSocket endpoint for real-time resource updates
	app.Get("/ws/resources", fiberwebsocket.New(websocket.HandleWebSocket, fiberwebsocket.Config{
		EnableCompression: true,
	}))
}
