package database

import (
	"log"
	"os"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

// DB is the database connection
var DB *sqlx.DB

// ConnectDB establishes a connection with the PostgreSQL database
func ConnectDB() {
	// Use Neon tech connection string
	connStr := getEnv("DATABASE_URL", "postgresql://neondb_owner:npg_top1bMIYlw9Z@ep-billowing-sky-a139oeav-pooler.ap-southeast-1.aws.neon.tech/neondb?sslmode=require")

	// Open connection to the database
	db, err := sqlx.Connect("postgres", connStr)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Test the connection
	err = db.Ping()
	if err != nil {
		log.Fatalf("Failed to ping database: %v", err)
	}

	// Set max open connections
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(25)

	DB = db
	log.Println("Connected to database successfully!")
}

// InitDB initializes the database schema if it doesn't exist
func InitDB() {
	// Create tables if they don't exist
	createTables()
}

// createTables creates all necessary tables for the application
func createTables() {
	// Create users table
	_, err := DB.Exec(`
	CREATE TABLE IF NOT EXISTS users (
		id SERIAL PRIMARY KEY,
		username VARCHAR(50) UNIQUE NOT NULL,
		email VARCHAR(255) UNIQUE NOT NULL,
		password_hash VARCHAR(255) NOT NULL,
		role VARCHAR(20) NOT NULL DEFAULT 'student',
		first_name VARCHAR(50),
		last_name VARCHAR(50),
		phone VARCHAR(20),
		created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
	)`)
	if err != nil {
		log.Fatalf("Failed to create users table: %v", err)
	}

	// Create items table
	_, err = DB.Exec(`
	CREATE TABLE IF NOT EXISTS items (
		id SERIAL PRIMARY KEY,
		title VARCHAR(100) NOT NULL,
		description TEXT,
		category VARCHAR(50) NOT NULL,
		status VARCHAR(20) NOT NULL DEFAULT 'lost',
		location VARCHAR(255),
		lost_time TIMESTAMP WITH TIME ZONE,
		report_time TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
		claimed_time TIMESTAMP WITH TIME ZONE,
		reporter_id INTEGER REFERENCES users(id),
		finder_id INTEGER REFERENCES users(id),
		image_url VARCHAR(255),
		created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
	)`)
	if err != nil {
		log.Fatalf("Failed to create items table: %v", err)
	}

	// Create images table
	_, err = DB.Exec(`
	CREATE TABLE IF NOT EXISTS images (
		id SERIAL PRIMARY KEY,
		item_id INTEGER REFERENCES items(id),
		image_url VARCHAR(255) NOT NULL,
		timestamp TIMESTAMP WITH TIME ZONE,
		latitude DOUBLE PRECISION,
		longitude DOUBLE PRECISION,
		created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
	)`)
	if err != nil {
		log.Fatalf("Failed to create images table: %v", err)
	}

	// Create messages table
	_, err = DB.Exec(`
	CREATE TABLE IF NOT EXISTS messages (
		id SERIAL PRIMARY KEY,
		sender_id INTEGER REFERENCES users(id),
		receiver_id INTEGER REFERENCES users(id),
		item_id INTEGER REFERENCES items(id),
		content TEXT NOT NULL,
		read BOOLEAN DEFAULT FALSE,
		created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
	)`)
	if err != nil {
		log.Fatalf("Failed to create messages table: %v", err)
	}

	// Create reports table
	_, err = DB.Exec(`
	CREATE TABLE IF NOT EXISTS reports (
		id SERIAL PRIMARY KEY,
		reporter_id INTEGER REFERENCES users(id),
		reported_id INTEGER REFERENCES users(id),
		item_id INTEGER REFERENCES items(id),
		reason TEXT NOT NULL,
		status VARCHAR(20) DEFAULT 'pending',
		created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
	)`)
	if err != nil {
		log.Fatalf("Failed to create reports table: %v", err)
	}

	log.Println("Database schema initialized")
}

// getEnv gets the environment variable or returns a default value
func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}
