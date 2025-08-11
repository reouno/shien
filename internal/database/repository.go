package database

import (
	"shien/internal/database/repository"
)

// Repository aggregates all repositories
type Repository struct {
	db           *DB
	activity     *repository.ActivityRepo
	gamification *repository.GamificationRepo
}

// NewRepository creates a new repository manager
func NewRepository(db *DB) *Repository {
	return &Repository{
		db:           db,
		activity:     repository.NewActivityRepo(db.Conn()),
		gamification: repository.NewGamificationRepo(db.Conn()),
	}
}

// Activity returns the activity repository
func (r *Repository) Activity() *repository.ActivityRepo {
	return r.activity
}

// Gamification returns the gamification repository
func (r *Repository) Gamification() *repository.GamificationRepo {
	return r.gamification
}
