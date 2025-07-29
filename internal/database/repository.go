package database

import (
	"shien/internal/database/repository"
)

// Repository aggregates all repositories
type Repository struct {
	db       *DB
	activity *repository.ActivityRepo
}

// NewRepository creates a new repository manager
func NewRepository(db *DB) *Repository {
	return &Repository{
		db:       db,
		activity: repository.NewActivityRepo(db.Conn()),
	}
}

// Activity returns the activity repository
func (r *Repository) Activity() *repository.ActivityRepo {
	return r.activity
}
