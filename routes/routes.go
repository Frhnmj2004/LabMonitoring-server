package routes

import (
	"github.com/Frhnmj2004/LabMonitoring-server/controllers"
	//"github.com/Frhnmj2004/LabMonitoring-server/middleware"
	"github.com/Frhnmj2004/LabMonitoring-server/websocket"
	"github.com/gofiber/fiber/v2"
	fiberwebsocket "github.com/gofiber/websocket/v2"
)

func SetupRoutes(app *fiber.App) {
	// Public routes
	api := app.Group("/api/v1")
	api.Post("/login", controllers.Login)

	// System signup routes (public)
	api.Post("/system-signup", controllers.SystemSignup)
	api.Get("/download-collector", controllers.DownloadCollector)

	// Protected routes
	//protected := api.Use(middleware.AuthMiddleware())

	// Admin-only routes
	//admin := protected.Use(middleware.AdminOnly())
	api.Post("/signup", controllers.Signup)

	// Resource routes
	api.Post("/resource", controllers.PostResource)
	api.Get("/resources/history", controllers.GetHistory)

	// Computer management routes (admin only)
	api.Get("/computers", controllers.GetAllComputers)

	// Alert routes (admin only)
	alertGroup := api.Group("/alerts")
	alertGroup.Get("/", controllers.GetAlerts)
	alertGroup.Get("/active", controllers.GetActiveAlerts)
	alertGroup.Get("/history", controllers.GetAlertHistory)
	alertGroup.Get("/stats", controllers.GetAlertStats)
	alertGroup.Put("/:id/resolve", controllers.ResolveAlert)

	// Internet usage routes
	api.Post("/internet-usage", controllers.PostInternetUsage)
	api.Get("/internet-usage", controllers.GetInternetUsage)

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
