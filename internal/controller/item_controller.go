package controller

import (
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/omniflare/campus-lostandfound/internal/database"
	"github.com/omniflare/campus-lostandfound/internal/models"
)

// ReportLostItem creates a new lost item report
func ReportLostItem(c *fiber.Ctx) error {
	// Get user ID from JWT context
	userID := c.Locals("user_id").(int)

	// Parse request body
	var itemReq models.ItemRequest
	if err := c.BodyParser(&itemReq); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request format",
		})
	}

	// Validate required fields
	if itemReq.Title == "" || itemReq.Category == "" || itemReq.Location == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Title, category, and location are required",
		})
	}

	// Set default status to "lost"
	status := "lost"

	// Insert the lost item into the database
	now := time.Now()
	var itemID int
	err := database.DB.QueryRow(`
		INSERT INTO items (title, description, category, status, location, lost_time, report_time, reporter_id, created_at, updated_at) 
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
		RETURNING id
	`, itemReq.Title, itemReq.Description, itemReq.Category, status, itemReq.Location, itemReq.LostTime, now, userID, now, now).Scan(&itemID)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Error creating lost item report",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "Lost item reported successfully",
		"item_id": itemID,
	})
}

// ReportFoundItem creates a new found item report
func ReportFoundItem(c *fiber.Ctx) error {
	// Get user ID from JWT context
	userID := c.Locals("user_id").(int)

	// Parse request body
	var itemReq models.ItemRequest
	if err := c.BodyParser(&itemReq); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request format",
		})
	}

	// Validate required fields
	if itemReq.Title == "" || itemReq.Category == "" || itemReq.Location == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Title, category, and location are required",
		})
	}

	// Set status to "found"
	status := "found"

	// Insert the found item into the database
	now := time.Now()
	var itemID int
	err := database.DB.QueryRow(`
		INSERT INTO items (title, description, category, status, location, report_time, finder_id, created_at, updated_at) 
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
		RETURNING id
	`, itemReq.Title, itemReq.Description, itemReq.Category, status, itemReq.Location, now, userID, now, now).Scan(&itemID)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Error creating found item report",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "Found item reported successfully",
		"item_id": itemID,
	})
}

// GetItemDetails retrieves details of a specific item
func GetItemDetails(c *fiber.Ctx) error {
	// Get item ID from URL parameter
	itemID := c.Params("id")
	if itemID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Item ID is required",
		})
	}

	// Get item from database
	var item models.Item
	err := database.DB.Get(&item, "SELECT * FROM items WHERE id = $1", itemID)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Item not found",
		})
	}

	return c.Status(fiber.StatusOK).JSON(item)
}

// GetItems retrieves a list of items based on filters
func GetItems(c *fiber.Ctx) error {
	// Parse query parameters
	status := c.Query("status", "all")     // Filter by status: lost, found, claimed, returned, or all
	category := c.Query("category", "all") // Filter by category
	page := c.QueryInt("page", 1)
	limit := c.QueryInt("limit", 10)
	offset := (page - 1) * limit

	// Build the query based on filters
	query := "SELECT * FROM items WHERE 1=1"
	countQuery := "SELECT COUNT(*) FROM items WHERE 1=1"
	args := []interface{}{}
	argCount := 1

	if status != "all" {
		query += " AND status = $" + string(rune('0'+argCount))
		countQuery += " AND status = $" + string(rune('0'+argCount))
		args = append(args, status)
		argCount++
	}

	if category != "all" {
		query += " AND category = $" + string(rune('0'+argCount))
		countQuery += " AND category = $" + string(rune('0'+argCount))
		args = append(args, category)
		argCount++
	}

	// Add pagination
	query += " ORDER BY created_at DESC LIMIT $" + string(rune('0'+argCount)) + " OFFSET $" + string(rune('0'+argCount+1))
	args = append(args, limit, offset)

	// Get items from database
	var items []models.Item
	err := database.DB.Select(&items, query, args...)
	if err != nil {
		// Log the actual error for debugging
		fmt.Printf("Database error retrieving items: %v\n", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Error retrieving items: " + err.Error(),
		})
	}

	// Get total count for pagination
	var total int
	err = database.DB.Get(&total, countQuery, args[:argCount-1]...)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Error retrieving item count",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"items": items,
		"meta": fiber.Map{
			"total": total,
			"page":  page,
			"limit": limit,
			"pages": (total + limit - 1) / limit,
		},
	})
}

// SearchItems searches for items by title or description
func SearchItems(c *fiber.Ctx) error {
	// Parse query parameters
	query := c.Query("q")
	if query == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Search query is required",
		})
	}

	status := c.Query("status", "all")
	page := c.QueryInt("page", 1)
	limit := c.QueryInt("limit", 10)
	offset := (page - 1) * limit

	// Build the search query
	searchSQL := "SELECT * FROM items WHERE (title ILIKE $1 OR description ILIKE $1)"
	countSQL := "SELECT COUNT(*) FROM items WHERE (title ILIKE $1 OR description ILIKE $1)"
	args := []interface{}{"%" + query + "%"}
	argCount := 2

	if status != "all" {
		searchSQL += " AND status = $" + string(rune('0'+argCount))
		countSQL += " AND status = $" + string(rune('0'+argCount))
		args = append(args, status)
		argCount++
	}

	// Add pagination
	searchSQL += " ORDER BY created_at DESC LIMIT $" + string(rune('0'+argCount)) + " OFFSET $" + string(rune('0'+argCount+1))
	args = append(args, limit, offset)

	// Get items from database
	var items []models.Item
	err := database.DB.Select(&items, searchSQL, args...)
	if err != nil {
		// Log the actual error for debugging
		fmt.Printf("Database error searching items: %v\n", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Error searching items: " + err.Error(),
		})
	}

	// Get total count for pagination
	var total int
	err = database.DB.Get(&total, countSQL, args[:argCount-1]...)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Error retrieving item count",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"items": items,
		"meta": fiber.Map{
			"total": total,
			"page":  page,
			"limit": limit,
			"pages": (total + limit - 1) / limit,
		},
	})
}

// UpdateItemStatus updates the status of an item
func UpdateItemStatus(c *fiber.Ctx) error {
	// Get user ID and role from JWT context
	userID := c.Locals("user_id").(int)
	role := c.Locals("role").(string)

	// Get item ID from URL parameter
	itemID := c.Params("id")
	if itemID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Item ID is required",
		})
	}

	// Parse request body
	var statusReq struct {
		Status string `json:"status"`
	}
	if err := c.BodyParser(&statusReq); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request format",
		})
	}

	// Validate status
	validStatuses := map[string]bool{"lost": true, "found": true, "claimed": true, "returned": true}
	if !validStatuses[statusReq.Status] {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid status. Must be one of: lost, found, claimed, returned",
		})
	}

	// Get the item to check ownership
	var item models.Item
	err := database.DB.Get(&item, "SELECT * FROM items WHERE id = $1", itemID)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Item not found",
		})
	}

	// Check if user has permission to update this item
	// Guards and admins can update any item, users can only update their own items
	if role != "admin" && role != "guard" {
		isReporter := item.ReporterID != nil && *item.ReporterID == userID
		isFinder := item.FinderID != nil && *item.FinderID == userID

		if !isReporter && !isFinder {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"error": "You do not have permission to update this item",
			})
		}
	}

	// Update the status
	now := time.Now()
	_, err = database.DB.Exec("UPDATE items SET status = $1, updated_at = $2 WHERE id = $3",
		statusReq.Status, now, itemID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Error updating item status",
		})
	}

	// If status is "claimed", update the claimed time
	if statusReq.Status == "claimed" {
		_, err = database.DB.Exec("UPDATE items SET claimed_time = $1 WHERE id = $2", now, itemID)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Error updating claimed time",
			})
		}
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Item status updated successfully",
	})
}

// UploadItemImage handles image upload for an item
func UploadItemImage(c *fiber.Ctx) error {
	// Get item ID from URL parameter
	itemID := c.Params("id")
	if itemID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Item ID is required",
		})
	}

	// Check if item exists
	var exists bool
	err := database.DB.Get(&exists, "SELECT EXISTS(SELECT 1 FROM items WHERE id = $1)", itemID)
	if err != nil || !exists {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Item not found",
		})
	}

	// Get the file from form data
	file, err := c.FormFile("image")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "No image file provided",
		})
	}

	// Generate a unique filename
	filename := time.Now().Format("20060102150405") + "_" + file.Filename

	// Save the file to a directory
	err = c.SaveFile(file, "./uploads/"+filename)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Error saving image file",
		})
	}

	// Create image URL
	imageURL := "/uploads/" + filename

	// Extract metadata from the image (this would require additional libraries)
	// For now, we'll just save the image URL
	_, err = database.DB.Exec(`
		INSERT INTO images (item_id, image_url, timestamp, created_at) 
		VALUES ($1, $2, $3, $4)
	`, itemID, imageURL, time.Now(), time.Now())
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Error saving image information",
		})
	}

	// Update item with image URL
	_, err = database.DB.Exec("UPDATE items SET image_url = $1 WHERE id = $2", imageURL, itemID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Error updating item with image URL",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message":   "Image uploaded successfully",
		"image_url": imageURL,
	})
}

// GetUserItems gets items reported by the current user
func GetUserItems(c *fiber.Ctx) error {
	// Get user ID from JWT context
	userID := c.Locals("user_id").(int)

	// Parse query parameters
	status := c.Query("status", "all")
	page := c.QueryInt("page", 1)
	limit := c.QueryInt("limit", 10)
	offset := (page - 1) * limit

	// Build the query that correctly handles NULL reporter_id and finder_id values
	query := "SELECT * FROM items WHERE ((reporter_id IS NOT NULL AND reporter_id = $1) OR (finder_id IS NOT NULL AND finder_id = $1))"
	countQuery := "SELECT COUNT(*) FROM items WHERE ((reporter_id IS NOT NULL AND reporter_id = $1) OR (finder_id IS NOT NULL AND finder_id = $1))"
	args := []interface{}{userID}
	argCount := 2

	if status != "all" {
		query += " AND status = $" + string(rune('0'+argCount))
		countQuery += " AND status = $" + string(rune('0'+argCount))
		args = append(args, status)
		argCount++
	}

	// Add pagination
	query += " ORDER BY created_at DESC LIMIT $" + string(rune('0'+argCount)) + " OFFSET $" + string(rune('0'+argCount+1))
	args = append(args, limit, offset)
	// Get items from database
	var items []models.Item
	err := database.DB.Select(&items, query, args...)
	if err != nil {
		// Log the actual error for debugging
		fmt.Printf("Database error retrieving items: %v\n", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Error retrieving items: " + err.Error(),
		})
	}

	// Get total count for pagination
	var total int
	err = database.DB.Get(&total, countQuery, args[:argCount-1]...)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Error retrieving item count",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"items": items,
		"meta": fiber.Map{
			"total": total,
			"page":  page,
			"limit": limit,
			"pages": (total + limit - 1) / limit,
		},
	})
}
