package controller

import (
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/omniflare/campus-lostandfound/internal/database"
	"github.com/omniflare/campus-lostandfound/internal/models"
	"github.com/omniflare/campus-lostandfound/internal/utils/jwt"
	"golang.org/x/crypto/bcrypt"
)

// RegisterUser handles user registration
func RegisterUser(c *fiber.Ctx) error {
	// Log the content type to debug
	fmt.Println("Content-Type:", c.Get("Content-Type"))

	// Parse request body
	var register models.Register
	if err := c.BodyParser(&register); err != nil {
		fmt.Println("Error parsing request body:", err)
		// Return more detailed error message
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request format: " + err.Error(),
		})
	}

	// Validate required fields
	if register.Username == "" || register.Email == "" || register.Password == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Username, email and password are required",
		})
	}

	// Check if username already exists
	var count int
	err := database.DB.Get(&count, "SELECT COUNT(*) FROM users WHERE username = $1", register.Username)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Database error",
		})
	}
	if count > 0 {
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{
			"error": "Username already exists",
		})
	}

	// Check if email already exists
	err = database.DB.Get(&count, "SELECT COUNT(*) FROM users WHERE email = $1", register.Email)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Database error",
		})
	}
	if count > 0 {
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{
			"error": "Email already exists",
		})
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(register.Password), bcrypt.DefaultCost)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Error hashing password",
		})
	}

	// Create new user
	now := time.Now()
	var userID int
	err = database.DB.QueryRow(`
		INSERT INTO users (username, email, password_hash, role, first_name, last_name, phone, created_at, updated_at) 
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
		RETURNING id
	`, register.Username, register.Email, string(hashedPassword), "student", register.FirstName, register.LastName, register.Phone, now, now).Scan(&userID)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Error creating user",
		})
	}

	// Return success message
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "User registered successfully",
		"user_id": userID,
	})
}

// LoginUser handles user login
func LoginUser(c *fiber.Ctx) error {
	// Parse request body
	var login models.Login
	if err := c.BodyParser(&login); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request format",
		})
	}

	// Validate required fields
	if login.Username == "" || login.Password == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Username and password are required",
		})
	}

	// Get user from database
	var user models.User
	err := database.DB.Get(&user, "SELECT * FROM users WHERE username = $1", login.Username)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Invalid username or password",
		})
	}

	// Verify password
	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(login.Password))
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Invalid username or password",
		})
	}

	// Generate JWT token
	token, err := jwt.Generate(user)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Error generating token",
		})
	}

	// Return token
	return c.Status(fiber.StatusOK).JSON(models.TokenResponse{
		Token: token,
	})
}

// GetUserProfile gets the current user's profile
func GetUserProfile(c *fiber.Ctx) error {
	// Get user ID from JWT context
	userID := c.Locals("user_id").(int)

	// Get user from database
	var user models.User
	err := database.DB.Get(&user, "SELECT id, username, email, role, first_name, last_name, phone, created_at, updated_at FROM users WHERE id = $1", userID)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "User not found",
		})
	}

	return c.Status(fiber.StatusOK).JSON(user)
}

// UpdateUserProfile updates the current user's profile
func UpdateUserProfile(c *fiber.Ctx) error {
	// Get user ID from JWT context
	userID := c.Locals("user_id").(int)

	// Parse request body
	var updateData struct {
		FirstName string `json:"first_name"`
		LastName  string `json:"last_name"`
		Phone     string `json:"phone"`
		Email     string `json:"email"`
	}

	if err := c.BodyParser(&updateData); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request format",
		})
	}

	// Update user in database
	_, err := database.DB.Exec(`
		UPDATE users 
		SET first_name = $1, last_name = $2, phone = $3, email = $4, updated_at = $5 
		WHERE id = $6
	`, updateData.FirstName, updateData.LastName, updateData.Phone, updateData.Email, time.Now(), userID)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Error updating user profile",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Profile updated successfully",
	})
}

// ChangePassword changes the user's password
func ChangePassword(c *fiber.Ctx) error {
	// Get user ID from JWT context
	userID := c.Locals("user_id").(int)

	// Parse request body
	var passwordData struct {
		CurrentPassword string `json:"current_password"`
		NewPassword     string `json:"new_password"`
	}

	if err := c.BodyParser(&passwordData); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request format",
		})
	}

	// Get current password hash from database
	var passwordHash string
	err := database.DB.Get(&passwordHash, "SELECT password_hash FROM users WHERE id = $1", userID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Database error",
		})
	}

	// Verify current password
	err = bcrypt.CompareHashAndPassword([]byte(passwordHash), []byte(passwordData.CurrentPassword))
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Current password is incorrect",
		})
	}

	// Hash new password
	newHashedPassword, err := bcrypt.GenerateFromPassword([]byte(passwordData.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Error hashing password",
		})
	}

	// Update password in database
	_, err = database.DB.Exec("UPDATE users SET password_hash = $1, updated_at = $2 WHERE id = $3",
		string(newHashedPassword), time.Now(), userID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Error updating password",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Password changed successfully",
	})
}
