package sqlite

import (
	"database/sql"
	"log"

	"github.com/Li-Khan/forum/pkg/models"
)

// ForumModel - define the type that wraps the sql.DB connection pool
type ForumModel struct {
	DB *sql.DB
}

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

// CreatePost - adds a post to the database
func (m *ForumModel) CreatePost(post *models.Post) (int64, error) {
	stmt := `INSERT INTO "main"."posts"
	("user_id", "user_login", "title", "text", "created")
	VALUES (?, ?, ?, ?, ?);`
	result, err := m.DB.Exec(stmt, post.UserID, post.UserLogin, post.Title, post.Text, post.Created)
	if err != nil {
		return 0, err
	}

	return result.LastInsertId()
}

// GetAllPosts - retrieves information about all posts from the database
func (m *ForumModel) GetAllPosts() (*[]models.Post, error) {
	stmt := `SELECT * FROM "main"."posts" ORDER BY "created" DESC`
	rows, err := m.DB.Query(stmt)
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

		votesPost, err := m.getVotesPost(post.ID)
		if err != nil {
			return nil, err
		}

		post.Votes = *votesPost
		post.Tags = append(post.Tags, tags...)
		post.Comments = comments
		posts = append(posts, post)
		post.Tags = nil
	}

	return &posts, nil
}

func (m *ForumModel) getVotesPost(id int64) (*models.Vote, error) {
	stmt := `SELECT vote FROM vote_post WHERE post_id = ?`
	rows, err := m.DB.Query(stmt, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var votes models.Vote
	var vote int
	for rows.Next() {
		err := rows.Scan(&vote)
		if err != nil {
			log.Println(err)
			continue
		}
		if vote == 1 {
			votes.Like++
		} else {
			votes.Dislike++
		}
	}
	return &votes, nil
}

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

func (m *ForumModel) getVotesComment(id int64) (*models.Vote, error) {
	stmt := `SELECT vote FROM vote_comment WHERE comment_id = ?`
	rows, err := m.DB.Query(stmt, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var votes models.Vote
	var vote int
	for rows.Next() {
		err := rows.Scan(&vote)
		if err != nil {
			log.Println(err)
			continue
		}
		if vote == 1 {
			votes.Like++
		} else {
			votes.Dislike++
		}
	}
	return &votes, nil
}

// GetPostComments ...
func (m *ForumModel) getPostComments(postID int64) ([]models.Comment, error) {
	stmt := `SELECT * FROM "main"."comments" WHERE "post_id" = ? ORDER BY "created" DESC`

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

		votes, err := m.getVotesComment(comment.ID)
		if err != nil {
			log.Println(err)
			continue
		}

		comment.Votes = *votes
		comments = append(comments, comment)
	}

	return comments, nil
}

// CreateTags - adds tags to the database
func (m *ForumModel) CreateTags(tags []string) error {
	stmt := `INSERT INTO "main"."tags" (Tag) VALUES (?)`

	for _, tag := range tags {
		_, err := m.DB.Exec(stmt, tag)
		if err != nil {
			continue
		}
	}
	return nil
}

func (m *ForumModel) getTagsID(tags []string) ([]int64, error) {
	stmt := `SELECT "id" FROM "main"."tags" WHERE "tag" = ?`

	tagsID := []int64{}
	var id int64
	for _, tag := range tags {
		row := m.DB.QueryRow(stmt, tag)
		err := row.Scan(&id)
		if err != nil {
			return nil, err
		}
		tagsID = append(tagsID, id)
	}
	return tagsID, nil
}

// CreatePostsAndTags - links tags to posts
func (m *ForumModel) CreatePostsAndTags(postID int64, tags []string) error {
	tagsID, err := m.getTagsID(tags)
	if err != nil {
		return err
	}

	stmt := `INSERT INTO "main"."posts_and_tags"
	("post_id", "tag_id")
	VALUES (?, ?);`

	for _, tagID := range tagsID {
		_, err := m.DB.Exec(stmt, postID, tagID)
		if err != nil {
			return err
		}
	}

	return nil
}

func (m *ForumModel) getPostIDbyTagID(id int64) ([]int64, error) {
	stmt := `SELECT "post_id" FROM "main"."posts_and_tags" WHERE "tag_id" = ?`

	rows, err := m.DB.Query(stmt, id)
	if err != nil {
		return nil, err
	}

	var postsID []int64
	var postID int64

	for rows.Next() {
		err = rows.Scan(&postID)
		if err != nil {
			return nil, err
		}
		postsID = append(postsID, postID)
	}

	return postsID, nil
}

func (m *ForumModel) getTagIDbyPostID(postID int64) ([]int64, error) {
	stmt := `SELECT "tag_id" FROM "main"."posts_and_tags" WHERE "post_id" = ?`

	rows, err := m.DB.Query(stmt, postID)
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

func (m *ForumModel) getTagsByTagID(tagsID []int64) ([]string, error) {
	stmt := `SELECT "tag" FROM "main"."tags" WHERE "id" = ?`

	tags := []string{}
	var tag string
	for _, tagID := range tagsID {
		row := m.DB.QueryRow(stmt, tagID)
		err := row.Scan(&tag)
		if err != nil {
			return nil, err
		}
		tags = append(tags, tag)
	}
	return tags, nil
}

// GetAllTags - retrieves all tags from the database
func (m *ForumModel) GetAllTags() ([]string, error) {
	stmt := `SELECT tag FROM tags`

	rows, err := m.DB.Query(stmt)
	if err != nil {
		return nil, err
	}

	var tag string
	var tags []string
	for rows.Next() {
		err = rows.Scan(&tag)
		if err != nil {
			return nil, err
		}
		tags = append(tags, tag)
	}

	return tags, nil
}

// CreatePostVote - adds a vote for the post
func (m *ForumModel) CreatePostVote(userID, postID int64, vote int) error {
	stmtSelect := `SELECT id, vote FROM vote_post WHERE user_id = ? AND post_id = ?;`
	stmtExec := `INSERT INTO "main"."vote_post" (
		"user_id",
		"post_id",
		"vote")
		VALUES (?, ?, ?)`
	stmtDelete := `DELETE FROM "main"."vote_post" WHERE "id" = ?`

	var like int64
	var id int64

	row := m.DB.QueryRow(stmtSelect, userID, postID)
	err := row.Scan(&id, &like)
	if err == sql.ErrNoRows {
		_, err := m.DB.Exec(stmtExec, userID, postID, vote)
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
		_, err = m.DB.Exec(stmtExec, userID, postID, vote)
		if err != nil {
			return err
		}
	}

	return nil
}

// CreateCommentVote - adds a vote to a comment
func (m *ForumModel) CreateCommentVote(userID, commentID int64, vote int) error {
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
