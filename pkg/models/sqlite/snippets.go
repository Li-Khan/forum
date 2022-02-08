package sqlite

import (
	"database/sql"
	"log"

	"github.com/Li-Khan/forum/pkg/models"
)

// SnippetModel - define the type that wraps the sql.DB connection pool
type SnippetModel struct {
	DB *sql.DB
}

/* ===== METHODS FOR THE USER ===== */

// CreateUser ...
func (m *SnippetModel) CreateUser(user *models.User) (int64, error) {
	result, err := m.DB.Exec(`INSERT INTO users(
		Login,
		Email,
		Created,
		Password
	) VALUES (?, ?, ?, ?)
	`, user.Login, user.Email, user.Created, user.Password)
	if err != nil {
		return 0, err
	}

	return result.LastInsertId()
}

// GetUser ...
func (m *SnippetModel) GetUser(login string) (*models.User, error) {
	row := m.DB.QueryRow("SELECT * FROM users WHERE Login = ?", login)

	user := models.User{}
	err := row.Scan(&user.ID, &user.Login, &user.Email, &user.Created, &user.Password)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

/* ===== METHODS FOR THE POST ===== */

// CreatePost ...
func (m *SnippetModel) CreatePost(post *models.Post) (int64, error) {
	result, err := m.DB.Exec(`INSERT INTO posts (
		UserID,
		UserLogin,
		Title,
		Text,
		Created
	) VALUES (?, ?, ?, ?, ?)`, post.UserID, post.UserLogin, post.Title, post.Text, post.Created)
	if err != nil {
		return 0, err
	}

	return result.LastInsertId()
}

// GetPostById ...
func (m *SnippetModel) GetPostById(id int) (*models.Post, error) {
	return nil, nil
}

// GetAllPosts ...
func (m *SnippetModel) GetAllPosts() (*[]models.Post, error) {
	rows, err := m.DB.Query(`SELECT * FROM posts ORDER BY Created DESC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	post := models.Post{}
	posts := []models.Post{}

	for rows.Next() {
		err = rows.Scan(&post.ID, &post.UserID, &post.UserLogin, &post.Title, &post.Text, &post.Created)
		if err != nil {
			log.Println(err)
			continue
		}

		tagsID, err := m.getTagIDbyPostID(post.ID)
		if err != nil {
			log.Println(err)
			continue
		}

		tags, err := m.getTagsByTagID(tagsID)
		if err != nil {
			log.Println(err)
			continue
		}

		comments, err := m.getPostComments(post.ID)

		post.Tags = append(post.Tags, tags...)
		post.Comments = comments
		posts = append(posts, post)
		post.Tags = nil
	}

	return &posts, nil
}

/* ===== METHODS FOR THE COMMENT ===== */

// CreateComment ...
func (m *SnippetModel) CreateComment(comment *models.Comment) error {
	stmt := `INSERT INTO comments (
		UserID,
		PostID,
		Login,
		Text, 
		Created
	) VALUES (?, ?, ?, ?, ?)`
	_, err := m.DB.Exec(stmt, comment.UserID, comment.PostID, comment.Login, comment.Text, comment.Created)
	if err != nil {
		return err
	}
	return nil
}

// GetPostComments ...
func (m *SnippetModel) getPostComments(postID int64) ([]models.Comment, error) {
	stmt := "SELECT * FROM comments WHERE PostID = ? ORDER BY Created DESC"

	rows, err := m.DB.Query(stmt, postID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var comment models.Comment
	var comments []models.Comment

	for rows.Next() {
		err := rows.Scan(&comment.ID, &comment.PostID, &comment.UserID, &comment.Login, &comment.Text, &comment.Created)
		if err != nil {
			log.Println(err)
			continue
		}
		comments = append(comments, comment)
	}

	return comments, nil
}

/* ===== METHODS FOR THE TAGS ===== */

// CreateTags ...
func (m *SnippetModel) CreateTags(tags []string) error {
	for _, tag := range tags {
		_, err := m.DB.Exec(`INSERT INTO tags (Tag) VALUES (?)`, tag)
		if err != nil {
			continue
		}
	}
	return nil
}

func (m *SnippetModel) getTagsID(tags []string) ([]int, error) {
	tagsID := []int{}
	var id int
	for _, tag := range tags {
		row := m.DB.QueryRow("SELECT ID FROM tags WHERE Tag = ?", tag)
		err := row.Scan(&id)
		if err != nil {
			return nil, err
		}
		tagsID = append(tagsID, id)
	}
	return tagsID, nil
}

// CreatePostsAndTags ...
func (m *SnippetModel) CreatePostsAndTags(postId int64, tags []string) error {
	tagsID, err := m.getTagsID(tags)
	if err != nil {
		return err
	}

	for _, tagID := range tagsID {
		_, err := m.DB.Exec("INSERT INTO postsAndTags (PostID, TagID) VALUES(?, ?)", postId, tagID)
		if err != nil {
			return err
		}
	}

	return nil
}

func (m *SnippetModel) getTagIDbyPostID(postID int64) ([]int64, error) {
	rows, err := m.DB.Query("SELECT TagID FROM postsAndTags WHERE PostID = ?", postID)
	if err != nil {
		return nil, err
	}

	var tagsID []int64
	var tagID int64
	for rows.Next() {
		err = rows.Scan(&tagID)
		if err != nil {
			return nil, err
		}
		tagsID = append(tagsID, tagID)
	}

	return tagsID, nil
}

func (m *SnippetModel) getTagsByTagID(tagsID []int64) ([]string, error) {
	tags := []string{}
	var tag string
	for _, tagID := range tagsID {
		row := m.DB.QueryRow("SELECT Tag FROM tags WHERE ID = ?", tagID)
		err := row.Scan(&tag)
		if err != nil {
			return nil, err
		}
		tags = append(tags, tag)
	}
	return tags, nil
}

/* ===== METHODS FOR THE LIKE ===== */

// Like ...
func (m *SnippetModel) Like(user *models.Comment) (int, error) {
	return 0, nil
}

// GetLike ...
func (m *SnippetModel) GetLike(id int) (*models.Comment, error) {
	return nil, nil
}

/* ===== METHODS FOR THE TAG ===== */
