package middleware

import (
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/omniflare/campus-lostandfound/internal/utils/jwt"
)

// Auth middleware to validate JWT tokens
func Auth() fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Get the Authorization header
		authHeader := c.Get("Authorization")

		// Check if the header is empty or doesn't contain "Bearer "
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Unauthorized: Missing or invalid authorization token",
			})
		}

		// Extract the token
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

		// Validate the token
		claims, err := jwt.Validate(tokenString)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Unauthorized: " + err.Error(),
			})
		}

		// Store user information in context for later use
		c.Locals("user_id", claims.UserID)
		c.Locals("username", claims.Username)
		c.Locals("role", claims.Role)

		// Continue to the next middleware/handler
		return c.Next()
	}
}

// AdminOnly middleware to restrict endpoints to admin users only
func AdminOnly() fiber.Handler {
	return func(c *fiber.Ctx) error {
		// First check if the user is authenticated
		role := c.Locals("role")
		if role == nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Unauthorized: Authentication required",
			})
		}

		// Check if the user has admin role
		if role != "admin" {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"error": "Forbidden: Admin access required",
			})
		}

		// Continue to the next middleware/handler
		return c.Next()
	}
}

// GuardAndAdmin middleware to restrict endpoints to guard and admin users
func GuardAndAdmin() fiber.Handler {
	return func(c *fiber.Ctx) error {
		// First check if the user is authenticated
		role := c.Locals("role")
		if role == nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Unauthorized: Authentication required",
			})
		}

		// Check if the user has guard or admin role
		if role != "guard" && role != "admin" {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"error": "Forbidden: Guard or admin access required",
			})
		}

		// Continue to the next middleware/handler
		return c.Next()
	}
}
