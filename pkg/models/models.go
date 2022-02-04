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
	ID      int    `json:"id"`
	UserID  int    `json:"userId"`
	Title   string `json:"title"`
	Text    string `json:"text"`
	Created string `json:"created"`
}

// Comment ...
type Comment struct {
	ID      int    `json:"id"`
	PostID  int    `json:"postId"`
	UserID  int    `json:"userId"`
	Text    string `json:"text"`
	Created string `json:"created"`
}

// Like ...
type Like struct {
	ID        int  `json:"id"`
	UserID    int  `json:"userId"`
	PostID    int  `json:"postId"`
	IsLike    bool `json:"like"`
	IsDislike bool `json:"dislike"`
}

// Errors ...
type Errors struct {
	IsPassNotMatch           bool
	IsInvalidForm            bool
	IsInvalidLoginOrPassword bool
	IsAlreadyExist           bool
}
