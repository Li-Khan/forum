package sqlite

import "github.com/Li-Khan/forum/pkg/models"

// CreateComment - adds a comment to the database
func (m *ForumModel) CreateComment(comment *models.Comment) error {
	stmt := `INSERT INTO "main"."comments"
	("user_id", "post_id", "login", "text", "created")
	VALUES (?, ?, ?, ?, ?);`

	_, err := m.DB.Exec(stmt, comment.UserID, comment.PostID, comment.Login, comment.Text, comment.Created)
	if err != nil {
		return err
	}
	return nil
}

// AddVoteToComment - adds a vote to a comment
func (m *ForumModel) AddVoteToComment(userID, commentID int64, vote int) error {
	stmtSelect := `SELECT id, vote FROM vote_comment WHERE user_id = ? AND comment_id = ?;`
	stmtExec := `INSERT INTO "main"."vote_comment" (
		"user_id",
		"comment_id",
		"vote")
		VALUES (?, ?, ?)`
	stmtDelete := `DELETE FROM "main"."vote_comment" WHERE "id" = ?`

	var like int64
	var id int64

	row := m.DB.QueryRow(stmtSelect, userID, commentID)
	err := row.Scan(&id, &like)
	if err != nil {
		_, err := m.DB.Exec(stmtExec, userID, commentID, vote)
		if err != nil {
			return err
		}
		return nil
	}

	_, err = m.DB.Exec(stmtDelete, id)
	if err != nil {
		return err
	}

	if int64(vote) != like {
		_, err = m.DB.Exec(stmtExec, userID, commentID, vote)
		if err != nil {
			return err
		}
	}

	return nil
}

// GetCommentByID - gets a comment from the database
func (m *ForumModel) GetCommentByID(id int64) (*models.Comment, error) {
	stmt := `SELECT * FROM comments WHERE id = ?`

	comment := models.Comment{}

	row := m.DB.QueryRow(stmt, id)
	err := row.Scan(&comment.ID, &comment.UserID, &comment.PostID, &comment.Login, &comment.Text, &comment.Created)
	if err != nil {
		return nil, err
	}

	return &comment, nil
}
