package models

import "errors"

// ErrNoRecord ...
var ErrNoRecord = errors.New("models: no suitable entry was found")

// User ...
type User struct {
	ID              int64  `json:"id"`
	Login           string `json:"login"`
	Email           string `json:"email"`
	Password        string `json:"password"`
	ConfirmPassword string `json:"confirm"`
	Created         string `json:"created"`
}

// Post ...
type Post struct {
	ID        int64  `json:"id"`
	UserID    int64  `json:"userId"`
	UserLogin string `json:"userLogin"`
	Title     string `json:"title"`
	Text      string `json:"text"`
	Tags      []string
	Comments  []Comment
	Created   string `json:"created"`
}

// Comment ...
type Comment struct {
	ID      int64  `json:"id"`
	PostID  int64  `json:"postId"`
	UserID  int64  `json:"userId"`
	Login   string `json:"login"`
	Text    string `json:"text"`
	Created string `json:"created"`
}

// Like ...
type Like struct {
	ID        int64 `json:"id"`
	UserID    int64 `json:"userId"`
	PostID    int64 `json:"postId"`
	IsLike    bool  `json:"like"`
	IsDislike bool  `json:"dislike"`
}

// Errors ...
type Errors struct {
	IsPassNotMatch           bool
	IsInvalidForm            bool
	IsInvalidLoginOrPassword bool
	IsAlreadyExist           bool
	IsInvalidLogin           bool
	IsInvalidEmail           bool
}
