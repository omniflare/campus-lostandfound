package controller

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/omniflare/campus-lostandfound/internal/database"
	"github.com/omniflare/campus-lostandfound/internal/models"
)

// GetUsers gets a list of users (admin only)
func GetUsers(c *fiber.Ctx) error {
	// Parse query parameters
	role := c.Query("role", "all")
	search := c.Query("search", "")
	page := c.QueryInt("page", 1)
	limit := c.QueryInt("limit", 20)
	offset := (page - 1) * limit

	// Build the query
	query := "SELECT id, username, email, role, first_name, last_name, phone, created_at, updated_at FROM users WHERE 1=1"
	countQuery := "SELECT COUNT(*) FROM users WHERE 1=1"
	args := []interface{}{}
	argCount := 1

	// Add role filter if provided
	if role != "all" {
		query += " AND role = $" + string(rune('0'+argCount))
		countQuery += " AND role = $" + string(rune('0'+argCount))
		args = append(args, role)
		argCount++
	}

	// Add search filter if provided
	if search != "" {
		query += " AND (username ILIKE $" + string(rune('0'+argCount)) + " OR email ILIKE $" + string(rune('0'+argCount))
		query += " OR first_name ILIKE $" + string(rune('0'+argCount)) + " OR last_name ILIKE $" + string(rune('0'+argCount)) + ")"
		countQuery += " AND (username ILIKE $" + string(rune('0'+argCount)) + " OR email ILIKE $" + string(rune('0'+argCount))
		countQuery += " OR first_name ILIKE $" + string(rune('0'+argCount)) + " OR last_name ILIKE $" + string(rune('0'+argCount)) + ")"
		args = append(args, "%"+search+"%")
		argCount++
	}

	// Add pagination
	query += " ORDER BY created_at DESC LIMIT $" + string(rune('0'+argCount)) + " OFFSET $" + string(rune('0'+argCount+1))
	args = append(args, limit, offset)

	// Get users from database
	var users []models.User
	err := database.DB.Select(&users, query, args...)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Error retrieving users",
		})
	}

	// Get total count for pagination
	var total int
	err = database.DB.Get(&total, countQuery, args[:argCount-1]...)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Error retrieving user count",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"users": users,
		"meta": fiber.Map{
			"total": total,
			"page":  page,
			"limit": limit,
			"pages": (total + limit - 1) / limit,
		},
	})
}

// UpdateUserRole updates a user's role (admin only)
func UpdateUserRole(c *fiber.Ctx) error {
	// Get user ID from URL parameter
	userID := c.Params("id")
	if userID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "User ID is required",
		})
	}

	// Get admin ID to ensure they're not changing their own role
	adminID := c.Locals("user_id").(int)
	if userID == string(rune('0'+adminID)) { // Convert userID string to match adminID for comparison
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error": "Cannot change your own role",
		})
	}

	// Parse request body
	var roleUpdate struct {
		Role string `json:"role"`
	}
	if err := c.BodyParser(&roleUpdate); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request format",
		})
	}

	// Validate role
	validRoles := map[string]bool{"student": true, "guard": true, "admin": true}
	if !validRoles[roleUpdate.Role] {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid role. Must be one of: student, guard, admin",
		})
	}

	// Update user role in database
	_, err := database.DB.Exec("UPDATE users SET role = $1, updated_at = $2 WHERE id = $3",
		roleUpdate.Role, time.Now(), userID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Error updating user role",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "User role updated successfully",
	})
}

// GetReports gets a list of reports (admin only)
func GetReports(c *fiber.Ctx) error {
	// Parse query parameters
	status := c.Query("status", "all")
	page := c.QueryInt("page", 1)
	limit := c.QueryInt("limit", 20)
	offset := (page - 1) * limit

	// Build the query
	query := `
		SELECT r.*, 
		   reporter.username as reporter_username,
		   reported.username as reported_username
		FROM reports r
		JOIN users reporter ON r.reporter_id = reporter.id
		JOIN users reported ON r.reported_id = reported.id
		WHERE 1=1
	`
	countQuery := "SELECT COUNT(*) FROM reports WHERE 1=1"
	args := []interface{}{}
	argCount := 1

	// Add status filter if provided
	if status != "all" {
		query += " AND r.status = $" + string(rune('0'+argCount))
		countQuery += " AND status = $" + string(rune('0'+argCount))
		args = append(args, status)
		argCount++
	}

	// Add pagination
	query += " ORDER BY r.created_at DESC LIMIT $" + string(rune('0'+argCount)) + " OFFSET $" + string(rune('0'+argCount+1))
	args = append(args, limit, offset)

	// Get reports from database
	var reports []struct {
		models.Report
		ReporterUsername string `db:"reporter_username" json:"reporter_username"`
		ReportedUsername string `db:"reported_username" json:"reported_username"`
	}
	err := database.DB.Select(&reports, query, args...)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Error retrieving reports",
		})
	}

	// Get total count for pagination
	var total int
	err = database.DB.Get(&total, countQuery, args[:argCount-1]...)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Error retrieving report count",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"reports": reports,
		"meta": fiber.Map{
			"total": total,
			"page":  page,
			"limit": limit,
			"pages": (total + limit - 1) / limit,
		},
	})
}

// UpdateReportStatus updates a report's status (admin only)
func UpdateReportStatus(c *fiber.Ctx) error {
	// Get report ID from URL parameter
	reportID := c.Params("id")
	if reportID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Report ID is required",
		})
	}

	// Parse request body
	var statusUpdate struct {
		Status  string `json:"status"`
		Comment string `json:"comment"`
	}
	if err := c.BodyParser(&statusUpdate); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request format",
		})
	}

	// Validate status
	validStatuses := map[string]bool{"pending": true, "resolved": true, "dismissed": true}
	if !validStatuses[statusUpdate.Status] {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid status. Must be one of: pending, resolved, dismissed",
		})
	}

	// Update report status in database
	_, err := database.DB.Exec(`
		UPDATE reports 
		SET status = $1, updated_at = $2, admin_comment = $3 
		WHERE id = $4
	`, statusUpdate.Status, time.Now(), statusUpdate.Comment, reportID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Error updating report status",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Report status updated successfully",
	})
}

// GetStats gets statistics for the admin dashboard
func GetStats(c *fiber.Ctx) error {
	// Get statistics from database
	stats := struct {
		TotalUsers     int `json:"total_users"`
		StudentCount   int `json:"student_count"`
		GuardCount     int `json:"guard_count"`
		AdminCount     int `json:"admin_count"`
		TotalItems     int `json:"total_items"`
		LostItems      int `json:"lost_items"`
		FoundItems     int `json:"found_items"`
		ClaimedItems   int `json:"claimed_items"`
		ReturnedItems  int `json:"returned_items"`
		PendingReports int `json:"pending_reports"`
	}{}

	// Get user counts
	err := database.DB.Get(&stats.TotalUsers, "SELECT COUNT(*) FROM users")
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Error retrieving user statistics",
		})
	}

	err = database.DB.Get(&stats.StudentCount, "SELECT COUNT(*) FROM users WHERE role = 'student'")
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Error retrieving student count",
		})
	}

	err = database.DB.Get(&stats.GuardCount, "SELECT COUNT(*) FROM users WHERE role = 'guard'")
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Error retrieving guard count",
		})
	}

	err = database.DB.Get(&stats.AdminCount, "SELECT COUNT(*) FROM users WHERE role = 'admin'")
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Error retrieving admin count",
		})
	}

	// Get item counts
	err = database.DB.Get(&stats.TotalItems, "SELECT COUNT(*) FROM items")
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Error retrieving item statistics",
		})
	}

	err = database.DB.Get(&stats.LostItems, "SELECT COUNT(*) FROM items WHERE status = 'lost'")
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Error retrieving lost items count",
		})
	}

	err = database.DB.Get(&stats.FoundItems, "SELECT COUNT(*) FROM items WHERE status = 'found'")
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Error retrieving found items count",
		})
	}

	err = database.DB.Get(&stats.ClaimedItems, "SELECT COUNT(*) FROM items WHERE status = 'claimed'")
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Error retrieving claimed items count",
		})
	}

	err = database.DB.Get(&stats.ReturnedItems, "SELECT COUNT(*) FROM items WHERE status = 'returned'")
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Error retrieving returned items count",
		})
	}

	// Get report counts
	err = database.DB.Get(&stats.PendingReports, "SELECT COUNT(*) FROM reports WHERE status = 'pending'")
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Error retrieving pending reports count",
		})
	}

	return c.Status(fiber.StatusOK).JSON(stats)
}

// CreateReport creates a new report
func CreateReport(c *fiber.Ctx) error {
	// Get user ID from JWT context
	userID := c.Locals("user_id").(int)

	// Parse request body
	var reportReq models.ReportRequest
	if err := c.BodyParser(&reportReq); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request format",
		})
	}

	// Validate required fields
	if reportReq.ReportedID == 0 || reportReq.Reason == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Reported user ID and reason are required",
		})
	}

	// Validate that reported user exists
	var exists bool
	err := database.DB.Get(&exists, "SELECT EXISTS(SELECT 1 FROM users WHERE id = $1)", reportReq.ReportedID)
	if err != nil || !exists {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Reported user not found",
		})
	}

	// Validate that item exists if provided
	if reportReq.ItemID != nil {
		err := database.DB.Get(&exists, "SELECT EXISTS(SELECT 1 FROM items WHERE id = $1)", *reportReq.ItemID)
		if err != nil || !exists {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Item not found",
			})
		}
	}

	// Insert the report into the database
	now := time.Now()
	var reportID int
	err = database.DB.QueryRow(`
		INSERT INTO reports (reporter_id, reported_id, item_id, reason, status, created_at, updated_at) 
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id
	`, userID, reportReq.ReportedID, reportReq.ItemID, reportReq.Reason, "pending", now, now).Scan(&reportID)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Error creating report",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message":   "Report submitted successfully",
		"report_id": reportID,
	})
}
