package sqlite

import (
	"database/sql"

	"github.com/Li-Khan/forum/pkg/models"
)

// ForumModel - define the type that wraps the sql.DB connection pool
type ForumModel struct {
	DB *sql.DB
}

/* ===== METHODS FOR THE USER ===== */

// InsertUser ...
func (m *ForumModel) InsertUser(user *models.User) (int, error) {
	return 0, nil
}

// GetUser ...
func (m *ForumModel) GetUser(id int) (*models.User, error) {
	return nil, nil
}

/* ===== METHODS FOR THE POST ===== */

// InsertPost ...
func (m *ForumModel) InsertPost(user *models.Post) (int, error) {
	return 0, nil
}

// GetPost ...
func (m *ForumModel) GetPost(id int) (*models.Post, error) {
	return nil, nil
}

/* ===== METHODS FOR THE COMMENT ===== */

// InsertComment ...
func (m *ForumModel) InsertComment(user *models.Comment) (int, error) {
	return 0, nil
}

// GetComment ...
func (m *ForumModel) GetComment(id int) (*models.Comment, error) {
	return nil, nil
}

/* ===== METHODS FOR THE LIKE ===== */

// InsertLike ...
func (m *ForumModel) InsertLike(user *models.Comment) (int, error) {
	return 0, nil
}

// GetLike ...
func (m *ForumModel) GetLike(id int) (*models.Comment, error) {
	return nil, nil
}

/* ===== METHODS FOR THE TAG ===== */
