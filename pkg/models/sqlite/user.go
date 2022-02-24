package sqlite

import "github.com/Li-Khan/forum/pkg/models"

// CreateUser - adds a user to the database
func (m *ForumModel) CreateUser(user *models.User) (int64, error) {
	stmt := `INSERT INTO "main"."users"(
		"login",
		"email",
		"created",
		"password"
	) VALUES (?, ?, ?, ?)`

	result, err := m.DB.Exec(stmt, user.Login, user.Email, user.Created, user.Password)
	if err != nil {
		return 0, err
	}

	return result.LastInsertId()
}

// GetUser - retrieves user information from the database
func (m *ForumModel) GetUser(login string) (*models.User, error) {
	stmt := `SELECT * FROM "main"."users" WHERE "login" = ?`

	row := m.DB.QueryRow(stmt, login)

	user := models.User{}
	err := row.Scan(&user.ID, &user.Login, &user.Email, &user.Created, &user.Password)
	if err != nil {
		return nil, err
	}

	return &user, nil
}
