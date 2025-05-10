package controller

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/omniflare/campus-lostandfound/internal/database"
	"github.com/omniflare/campus-lostandfound/internal/models"
)

// SendMessage handles sending a message between users
func SendMessage(c *fiber.Ctx) error {
	// Get sender ID from JWT context
	senderID := c.Locals("user_id").(int)

	// Parse request body
	var messageReq models.MessageRequest
	if err := c.BodyParser(&messageReq); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request format",
		})
	}

	// Validate required fields
	if messageReq.Content == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Message content is required",
		})
	}

	// Check if receiver exists
	var receiverExists bool
	err := database.DB.Get(&receiverExists, "SELECT EXISTS(SELECT 1 FROM users WHERE id = $1)", messageReq.ReceiverID)
	if err != nil || !receiverExists {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Receiver not found",
		})
	}

	// Check if item exists (if item_id is provided)
	if messageReq.ItemID != 0 {
		var itemExists bool
		err := database.DB.Get(&itemExists, "SELECT EXISTS(SELECT 1 FROM items WHERE id = $1)", messageReq.ItemID)
		if err != nil || !itemExists {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "Item not found",
			})
		}
	}

	// Insert message into database
	var messageID int
	err = database.DB.QueryRow(`
		INSERT INTO messages (sender_id, receiver_id, item_id, content, created_at) 
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id
	`, senderID, messageReq.ReceiverID, messageReq.ItemID, messageReq.Content, time.Now()).Scan(&messageID)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Error sending message",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message":    "Message sent successfully",
		"message_id": messageID,
	})
}

// GetConversations gets all conversations for the current user
func GetConversations(c *fiber.Ctx) error {
	// Get user ID from JWT context
	userID := c.Locals("user_id").(int)

	// This query groups messages by the other user involved and returns the latest message
	query := `
	WITH conversations AS (
		SELECT 
			CASE 
				WHEN sender_id = $1 THEN receiver_id 
				ELSE sender_id 
			END AS other_user_id,
			MAX(created_at) as latest_time
		FROM messages
		WHERE sender_id = $1 OR receiver_id = $1
		GROUP BY other_user_id
	)
	SELECT 
		c.other_user_id, 
		u.username as other_username, 
		m.id as latest_message_id, 
		m.content as latest_message, 
		m.created_at as latest_message_time,
		(SELECT COUNT(*) FROM messages WHERE receiver_id = $1 AND sender_id = c.other_user_id AND read = false) as unread_count
	FROM conversations c
	JOIN users u ON c.other_user_id = u.id
	JOIN messages m ON (
		(m.sender_id = $1 AND m.receiver_id = c.other_user_id) OR 
		(m.sender_id = c.other_user_id AND m.receiver_id = $1)
	)
	WHERE m.created_at = c.latest_time
	ORDER BY m.created_at DESC
	`

	// Get conversations from database
	var conversations []struct {
		OtherUserID       int       `db:"other_user_id" json:"other_user_id"`
		OtherUsername     string    `db:"other_username" json:"other_username"`
		LatestMessageID   int       `db:"latest_message_id" json:"latest_message_id"`
		LatestMessage     string    `db:"latest_message" json:"latest_message"`
		LatestMessageTime time.Time `db:"latest_message_time" json:"latest_message_time"`
		UnreadCount       int       `db:"unread_count" json:"unread_count"`
	}

	err := database.DB.Select(&conversations, query, userID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Error retrieving conversations",
		})
	}

	return c.Status(fiber.StatusOK).JSON(conversations)
}

// GetMessages gets all messages between the current user and another user
func GetMessages(c *fiber.Ctx) error {
	// Get user ID from JWT context
	userID := c.Locals("user_id").(int)

	// Get other user ID from URL parameter
	otherUserID := c.Params("id")
	if otherUserID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Other user ID is required",
		})
	}

	// Optional item ID for filtering
	itemID := c.Query("item_id", "")

	// Parse pagination parameters
	page := c.QueryInt("page", 1)
	limit := c.QueryInt("limit", 50)
	offset := (page - 1) * limit

	// Build the query
	query := `
		SELECT m.*, u.username as sender_username 
		FROM messages m
		JOIN users u ON m.sender_id = u.id
		WHERE (m.sender_id = $1 AND m.receiver_id = $2) OR (m.sender_id = $2 AND m.receiver_id = $1)
	`
	countQuery := `
		SELECT COUNT(*) 
		FROM messages
		WHERE (sender_id = $1 AND receiver_id = $2) OR (sender_id = $2 AND receiver_id = $1)
	`
	args := []interface{}{userID, otherUserID}
	argCount := 3

	// Add item filter if provided
	if itemID != "" {
		query += " AND m.item_id = $" + string(rune('0'+argCount))
		countQuery += " AND item_id = $" + string(rune('0'+argCount))
		args = append(args, itemID)
		argCount++
	}

	// Add ordering and pagination
	query += " ORDER BY m.created_at DESC LIMIT $" + string(rune('0'+argCount)) + " OFFSET $" + string(rune('0'+argCount+1))
	args = append(args, limit, offset)

	// Get messages from database
	var messages []struct {
		models.Message
		SenderUsername string `db:"sender_username" json:"sender_username"`
	}
	err := database.DB.Select(&messages, query, args...)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Error retrieving messages",
		})
	}

	// Get total count for pagination
	var total int
	err = database.DB.Get(&total, countQuery, args[:argCount-1]...)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Error retrieving message count",
		})
	}

	// Mark messages as read
	_, err = database.DB.Exec(`
		UPDATE messages 
		SET read = true 
		WHERE receiver_id = $1 AND sender_id = $2 AND read = false
	`, userID, otherUserID)
	if err != nil {
		// Log the error but don't fail the request
		// Ideally, we'd use a logger here
		// log.Printf("Error marking messages as read: %v", err)
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"messages": messages,
		"meta": fiber.Map{
			"total": total,
			"page":  page,
			"limit": limit,
			"pages": (total + limit - 1) / limit,
		},
	})
}

// GetUnreadMessageCount gets the count of unread messages for the current user
func GetUnreadMessageCount(c *fiber.Ctx) error {
	// Get user ID from JWT context
	userID := c.Locals("user_id").(int)

	// Get unread count from database
	var unreadCount int
	err := database.DB.Get(&unreadCount, "SELECT COUNT(*) FROM messages WHERE receiver_id = $1 AND read = false", userID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Error retrieving unread message count",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"unread_count": unreadCount,
	})
}
