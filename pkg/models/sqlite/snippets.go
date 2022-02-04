package sqlite

import (
	"database/sql"
	"fmt"

	"github.com/Li-Khan/forum/pkg/models"
)

// SnippetModel - define the type that wraps the sql.DB connection pool
type SnippetModel struct {
	DB *sql.DB
}

/* ===== METHODS FOR THE USER ===== */

// CreateUser ...
func (m *SnippetModel) CreateUser(user *models.User) (int64, error) {
	info, err := m.DB.Exec(`INSERT INTO users(
		Login,
		Email,
		Created,
		Password
	) VALUES (?, ?, ?, ?)
	`, user.Login, user.Email, user.Created, user.Password)
	if err != nil {
		return 0, fmt.Errorf("create user: %w", err)
	}

	return info.LastInsertId()
}

// GetUser ...
func (m *SnippetModel) GetUser(id int) (*models.User, error) {
	return nil, nil
}

/* ===== METHODS FOR THE POST ===== */

// CreatePost ...
func (m *SnippetModel) CreatePost(user *models.Post) (int, error) {
	return 0, nil
}

// GetPost ...
func (m *SnippetModel) GetPost(id int) (*models.Post, error) {
	return nil, nil
}

/* ===== METHODS FOR THE COMMENT ===== */

// CreateComment ...
func (m *SnippetModel) CreateComment(user *models.Comment) (int, error) {
	return 0, nil
}

// GetComment ...
func (m *SnippetModel) GetComment(id int) (*models.Comment, error) {
	return nil, nil
}

/* ===== METHODS FOR THE LIKE ===== */

// InsertLike ...
func (m *SnippetModel) Like(user *models.Comment) (int, error) {
	return 0, nil
}

// GetLike ...
func (m *SnippetModel) GetLike(id int) (*models.Comment, error) {
	return nil, nil
}

/* ===== METHODS FOR THE TAG ===== */
