package models

// ErrNoRecord ...
// var ErrNoRecord = errors.New("models: no suitable entry was found")

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
	Votes     Vote
	Comments  []Comment
	Created   string `json:"created"`
}

// Vote ...
type Vote struct {
	Like    uint64
	Dislike uint64
}

// Comment ...
type Comment struct {
	ID      int64  `json:"id"`
	PostID  int64  `json:"postId"`
	UserID  int64  `json:"userId"`
	Login   string `json:"login"`
	Text    string `json:"text"`
	Created string `json:"created"`
	Votes   Vote
}

// Like ...
type Like struct {
	ID     int64 `json:"id,omitempty"`
	UserID int64 `json:"user_id,omitempty"`
	PostID int64 `json:"post_id,omitempty"`
	IsLike bool  `json:"is_like,omitempty"`
}

// Errors ...
type Errors struct {
	IsInvalidLoginOrPassword bool
	IsAlreadyExist           bool
}
