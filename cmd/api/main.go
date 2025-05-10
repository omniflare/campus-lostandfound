package main

import (
	"encoding/json"
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/omniflare/campus-lostandfound/internal/database"
	"github.com/omniflare/campus-lostandfound/internal/routes"
)

func main() {
	// Create uploads directory if it doesn't exist
	err := os.MkdirAll("./uploads", 0755)
	if err != nil {
		log.Fatalf("Error creating uploads directory: %v", err)
	}

	// Connect to database
	database.ConnectDB()
	database.InitDB()

	// Create a new Fiber app
	app := fiber.New(fiber.Config{
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			// Default 500 status code
			code := fiber.StatusInternalServerError

			// Check if it's a fiber error
			if e, ok := err.(*fiber.Error); ok {
				code = e.Code
			}

			// Return JSON response
			return c.Status(code).JSON(fiber.Map{
				"error": err.Error(),
			})
		},
		// Set more lenient JSON parsing
		JSONEncoder: json.Marshal,
		JSONDecoder: json.Unmarshal,
		// Increase request body limit
		BodyLimit: 10 * 1024 * 1024, // 10MB
	})

	// Use middlewares
	app.Use(recover.New())
	app.Use(logger.New())
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*", // Allow all origins
		AllowHeaders: "Origin, Content-Type, Accept, Authorization",
		AllowMethods: "GET, POST, PUT, DELETE",
	}))

	// Serve static files from uploads directory
	app.Static("/uploads", "./uploads")

	// Setup routes
	routes.SetupRoutes(app)

	// Start server
	port := getEnv("PORT", "3000")
	log.Printf("Server starting on port %s", port)
	log.Fatal(app.Listen(":" + port))
}

// getEnv gets the environment variable or returns a default value
func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}
