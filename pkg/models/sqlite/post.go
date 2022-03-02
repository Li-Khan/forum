package sqlite

import (
	"database/sql"
	"log"

	"github.com/Li-Khan/forum/pkg/models"
)

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

// GetPostByID ...
func (m *ForumModel) GetPostByID(id int) (*models.Post, error) {
	stmt := `SELECT * FROM "main"."posts" WHERE "id" = ?`
	row := m.DB.QueryRow(stmt, id)
	post := models.Post{}
	err := row.Scan(&post.ID, &post.UserID, &post.UserLogin, &post.Title, &post.Text, &post.Created)
	if err != nil {
		return nil, err
	}
	return &post, nil
}

// LinkTagsToAPost - links tags to posts
func (m *ForumModel) LinkTagsToAPost(postID int64, tags []string) error {
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

// AddVoteToPost - adds a vote for the post
func (m *ForumModel) AddVoteToPost(userID, postID int64, vote int) error {
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
