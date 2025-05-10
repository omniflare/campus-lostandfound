package models

import "time"

// User represents a user in the system
type User struct {
	ID           int       `db:"id" json:"id"`
	Username     string    `db:"username" json:"username"`
	Email        string    `db:"email" json:"email"`
	PasswordHash string    `db:"password_hash" json:"-"`
	Role         string    `db:"role" json:"role"` // student, guard, admin
	FirstName    string    `db:"first_name" json:"first_name"`
	LastName     string    `db:"last_name" json:"last_name"`
	Phone        string    `db:"phone" json:"phone"`
	CreatedAt    time.Time `db:"created_at" json:"created_at"`
	UpdatedAt    time.Time `db:"updated_at" json:"updated_at"`
}

// Item represents an item in the lost and found system
type Item struct {
	ID          int        `db:"id" json:"id"`
	Title       string     `db:"title" json:"title"`
	Description string     `db:"description" json:"description"`
	Category    string     `db:"category" json:"category"`
	Status      string     `db:"status" json:"status"` // lost, found, claimed, returned
	Location    string     `db:"location" json:"location"`
	LostTime    *time.Time `db:"lost_time" json:"lost_time"`
	ReportTime  time.Time  `db:"report_time" json:"report_time"`
	ClaimedTime *time.Time `db:"claimed_time" json:"claimed_time"`
	ReporterID  *int       `db:"reporter_id" json:"reporter_id"`
	FinderID    *int       `db:"finder_id" json:"finder_id"`
	ImageURL    *string    `db:"image_url" json:"image_url"`
	CreatedAt   time.Time  `db:"created_at" json:"created_at"`
	UpdatedAt   time.Time  `db:"updated_at" json:"updated_at"`
}

// Image represents an image of an item with metadata
type Image struct {
	ID        int        `db:"id" json:"id"`
	ItemID    int        `db:"item_id" json:"item_id"`
	ImageURL  string     `db:"image_url" json:"image_url"`
	Timestamp *time.Time `db:"timestamp" json:"timestamp"`
	Latitude  *float64   `db:"latitude" json:"latitude"`
	Longitude *float64   `db:"longitude" json:"longitude"`
	CreatedAt time.Time  `db:"created_at" json:"created_at"`
}

// Message represents a message between users about an item
type Message struct {
	ID         int       `db:"id" json:"id"`
	SenderID   int       `db:"sender_id" json:"sender_id"`
	ReceiverID int       `db:"receiver_id" json:"receiver_id"`
	ItemID     int       `db:"item_id" json:"item_id"`
	Content    string    `db:"content" json:"content"`
	Read       bool      `db:"read" json:"read"`
	CreatedAt  time.Time `db:"created_at" json:"created_at"`
}

// Report represents a report of abuse or suspicious activity
type Report struct {
	ID         int       `db:"id" json:"id"`
	ReporterID *int      `db:"reporter_id" json:"reporter_id"`
	ReportedID int       `db:"reported_id" json:"reported_id"`
	ItemID     *int      `db:"item_id" json:"item_id"`
	Reason     string    `db:"reason" json:"reason"`
	Status     string    `db:"status" json:"status"` // pending, resolved, dismissed
	CreatedAt  time.Time `db:"created_at" json:"created_at"`
	UpdatedAt  time.Time `db:"updated_at" json:"updated_at"`
}

// Login represents the login request payload
type Login struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// Register represents the registration request payload
type Register struct {
	Username  string `json:"username"`
	Email     string `json:"email"`
	Password  string `json:"password"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Phone     string `json:"phone"`
}

// ItemRequest represents the item request payload
type ItemRequest struct {
	Title       string     `json:"title"`
	Description string     `json:"description"`
	Category    string     `json:"category"`
	Location    string     `json:"location"`
	LostTime    *time.Time `json:"lost_time"`
	Status      string     `json:"status"` // For found items, this would be "found"
}

// MessageRequest represents a message request payload
type MessageRequest struct {
	ReceiverID int    `json:"receiver_id"`
	ItemID     int    `json:"item_id"`
	Content    string `json:"content"`
}

// ReportRequest represents a report request payload
type ReportRequest struct {
	ReportedID int    `json:"reported_id"`
	ItemID     *int   `json:"item_id"`
	Reason     string `json:"reason"`
}

// TokenResponse is the response containing the JWT token
type TokenResponse struct {
	Token string `json:"token"`
}
