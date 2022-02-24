package sqlite

import (
	"database/sql"
)

// ForumModel - define the type that wraps the sql.DB connection pool
type ForumModel struct {
	DB *sql.DB
}
