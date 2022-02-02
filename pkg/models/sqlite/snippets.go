package sqlite

import (
	"database/sql"

	"github.com/Li-Khan/forum/pkg/models"
)

// SnippetModel - define the type that wraps the sql.DB connection pool
type SnippetModel struct {
	DB *sql.DB
}

/* ===== METHODS FOR THE USER ===== */

// InsertUser ...
func (m *SnippetModel) InsertUser(user *models.User) (int, error) {
	return 0, nil
}

// GetUser ...
func (m *SnippetModel) GetUser(id int) (*models.User, error) {
	return nil, nil
}

/* ===== METHODS FOR THE POST ===== */

// InsertPost ...
func (m *SnippetModel) InsertPost(user *models.Post) (int, error) {
	return 0, nil
}

// GetPost ...
func (m *SnippetModel) GetPost(id int) (*models.Post, error) {
	return nil, nil
}

/* ===== METHODS FOR THE COMMENT ===== */

// InsertComment ...
func (m *SnippetModel) InsertComment(user *models.Comment) (int, error) {
	return 0, nil
}

// GetComment ...
func (m *SnippetModel) GetComment(id int) (*models.Comment, error) {
	return nil, nil
}

/* ===== METHODS FOR THE LIKE ===== */

// InsertLike ...
func (m *SnippetModel) InsertLike(user *models.Comment) (int, error) {
	return 0, nil
}

// GetLike ...
func (m *SnippetModel) GetLike(id int) (*models.Comment, error) {
	return nil, nil
}

/* ===== METHODS FOR THE TAG ===== */
