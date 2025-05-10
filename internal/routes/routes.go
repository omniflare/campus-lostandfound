package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/omniflare/campus-lostandfound/internal/controller"
	"github.com/omniflare/campus-lostandfound/internal/middleware"
)

// SetupRoutes configures all the routes for the application
func SetupRoutes(app *fiber.App) {
	// Health check
	app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status":  "ok",
			"message": "Campus Lost and Found API is running",
		})
	})

	// API routes
	api := app.Group("/api")
	v1 := api.Group("/v1")

	// Auth routes - no authentication required
	auth := v1.Group("/auth")
	auth.Post("/register", controller.RegisterUser)
	auth.Post("/login", controller.LoginUser)

	// User routes - authentication required
	user := v1.Group("/user", middleware.Auth())
	user.Get("/profile", controller.GetUserProfile)
	user.Put("/profile", controller.UpdateUserProfile)
	user.Put("/password", controller.ChangePassword)
	user.Get("/items", controller.GetUserItems)
	user.Get("/messages/unread", controller.GetUnreadMessageCount)
	user.Get("/messages/conversations", controller.GetConversations)
	user.Get("/messages/:id", controller.GetMessages)
	user.Post("/messages", controller.SendMessage)
	user.Post("/reports", controller.CreateReport)

	// Item routes
	items := v1.Group("/items")
	items.Get("", controller.GetItems)           // Public - no auth required
	items.Get("/search", controller.SearchItems) // Public - no auth required
	items.Get("/:id", controller.GetItemDetails) // Public - no auth required

	// Protected item routes - authentication required
	itemsAuth := v1.Group("/items", middleware.Auth())
	itemsAuth.Post("/lost", controller.ReportLostItem)
	itemsAuth.Post("/found", controller.ReportFoundItem)
	itemsAuth.Put("/:id/status", controller.UpdateItemStatus)
	itemsAuth.Post("/:id/image", controller.UploadItemImage)

	// Guard routes - guard and admin only
	guard := v1.Group("/guard", middleware.GuardAndAdmin())
	guard.Get("/items", controller.GetItems) // Reusing the controller but with guard middleware

	// Admin routes - admin only
	admin := v1.Group("/admin", middleware.AdminOnly())
	admin.Get("/users", controller.GetUsers)
	admin.Put("/users/:id/role", controller.UpdateUserRole)
	admin.Get("/reports", controller.GetReports)
	admin.Put("/reports/:id/status", controller.UpdateReportStatus)
	admin.Get("/stats", controller.GetStats)
}
