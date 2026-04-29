package database

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/mattn/go-sqlite3" 
)

// InitDB initializes the SQLite database and creates necessary tables 
func InitDB(dataSourceName string) (*sql.DB, error) {
	db, err := sql.Open("sqlite3", dataSourceName)
	if err != nil {
		return nil, err
	}

	// Enable connection pooling as required by performance specs [cite: 771]
	db.SetMaxOpenConns(25)

	// Create Tables based on Project Spec 
	schema := `
	CREATE TABLE IF NOT EXISTS users (
		id TEXT PRIMARY KEY,
		username TEXT UNIQUE,
		password_hash TEXT,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);

	CREATE TABLE IF NOT EXISTS manga (
		id TEXT PRIMARY KEY,
		title TEXT,
		author TEXT,
		genres TEXT, 
		status TEXT,
		total_chapters INTEGER,
		description TEXT
	);

	CREATE TABLE IF NOT EXISTS user_progress (
		user_id TEXT,
		manga_id TEXT,
		current_chapter INTEGER,
		status TEXT,
		updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		PRIMARY KEY (user_id, manga_id)
	);`

	_, err = db.Exec(schema)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize schema: %w", err)
	}

	log.Println("Database tables initialized successfully.")
	return db, nil
}