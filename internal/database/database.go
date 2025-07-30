package database

import (
	"database/sql"
	"fmt"
	"sync"
	
	_ "github.com/mattn/go-sqlite3" // SQLite driver
	"shien/internal/paths"
)

// DB represents the database connection
type DB struct {
	conn   *sql.DB
	path   string
	mu     sync.RWMutex
}

// New creates a new database connection
func New() (*DB, error) {
	// Database file path
	dbPath := paths.DatabaseFile()
	
	// Open database connection
	conn, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}
	
	// Enable foreign keys
	if _, err := conn.Exec("PRAGMA foreign_keys = ON"); err != nil {
		conn.Close()
		return nil, fmt.Errorf("failed to enable foreign keys: %w", err)
	}
	
	// Set connection pool settings
	conn.SetMaxOpenConns(1) // SQLite doesn't benefit from multiple connections
	
	db := &DB{
		conn: conn,
		path: dbPath,
	}
	
	// Run migrations
	if err := db.Migrate(); err != nil {
		conn.Close()
		return nil, fmt.Errorf("failed to run migrations: %w", err)
	}
	
	return db, nil
}

// Close closes the database connection
func (db *DB) Close() error {
	return db.conn.Close()
}

// Path returns the database file path
func (db *DB) Path() string {
	return db.path
}

// Conn returns the underlying database connection
func (db *DB) Conn() *sql.DB {
	return db.conn
}

// Transaction executes a function within a database transaction
func (db *DB) Transaction(fn func(*sql.Tx) error) error {
	tx, err := db.conn.Begin()
	if err != nil {
		return err
	}
	
	if err := fn(tx); err != nil {
		tx.Rollback()
		return err
	}
	
	return tx.Commit()
}