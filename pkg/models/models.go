package models

import "errors"

// ErrNoRecord ...
var ErrNoRecord = errors.New("models: no suitable entry was found")

// User ...
type User struct {
	ID       int
	Login    string
	Email    string
	Password string
	Created  string
}

// Post ...
type Post struct {
	ID      int
	UserID  int
	Title   string
	Text    string
	Created string
}

// Comment ...
type Comment struct {
	ID      int
	PostID  int
	UserID  int
	Text    string
	Created string
}

// Like ...
type Like struct {
	ID        int
	UserID    int
	PostID    int
	IsLike    bool
	IsDislike bool
}
