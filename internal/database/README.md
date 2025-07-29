# Database Architecture

This directory contains the database layer for Shien, using SQLite with a repository pattern.

## Structure

- `database.go` - Database connection and initialization
- `migrations.go` - Migration runner (executes migrations)
- `repository.go` - Repository manager (aggregates all repositories)
- `schema.sql` - Reference schema (documentation only)
- `repository/` - Domain-specific repository implementations
  - `activity.go` - Activity log repository (includes model definition)
- `migrations/` - Migration files
  - `migrations.go` - Migration registry (lists all migrations)
  - `001_activity_logs.go` - First migration

## Adding New Tables

1. Create a new migration file:
```go
// migrations/002_users.go
var Migration002Users = Migration{
    Version:     2,
    Description: "Create users table",
    Up: func(tx *sql.Tx) error {
        // Create table
    },
}
```

2. Add to migration registry in `migrations/migrations.go`

3. Create repository in `repository/user.go`:
```go
package repository

type User struct {
    ID    int64  `json:"id"`
    Name  string `json:"name"`
}

type UserRepo struct {
    conn *sql.DB
}

func NewUserRepo(conn *sql.DB) *UserRepo {
    return &UserRepo{conn: conn}
}

func (r *UserRepo) Create(user *User) error {
    // Implementation
}
```

4. Add to repository manager in `repository.go`:
```go
type Repository struct {
    db       *DB
    activity *repository.ActivityRepo
    user     *repository.UserRepo  // Add this
}
```

## Usage

```go
// Record activity
repo.Activity().RecordActivity()

// Get activity logs
logs, err := repo.Activity().GetActivityLogs(from, to)
```