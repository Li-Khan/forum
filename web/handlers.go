package web

import (
	"database/sql"
	"fmt"
	"net/http"
	"net/mail"
	"regexp"
	"strconv"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/Li-Khan/forum/pkg/models"
	"golang.org/x/crypto/bcrypt"
)

const (
	// FORMAT - time format.
	FORMAT string = "01-02-2006 15:04:05"
	// Salt - adds additional characters to the user's password to make the password more complex.
	Salt string = "758+}+%s#=?.^$69,"
)

// home - main page handler
func (app *Application) home(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		app.notFound(w)
		return
	}

	if r.Method != http.MethodGet {
		w.Header().Set("Allow", http.MethodGet)
		app.methodNotAllowed(w)
		return
	}

	posts, err := app.Forum.GetAllPosts()
	if err != nil {
		app.serverError(w, err)
		return
	}

	tags, err := app.Forum.GetAllTags()
	if err != nil {
		app.serverError(w, err)
		return
	}

	app.render(w, r, "home.page.html", &templateData{
		Posts:     posts,
		Tags:      tags,
		IsSession: isSession(r),
	})
}

func (app *Application) signup(w http.ResponseWriter, r *http.Request) {
	if isSession(r) {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	data := &templateData{IsSession: isSession(r)}

	switch r.Method {
	case http.MethodGet:
		app.render(w, r, "signup.page.html", data)
	case http.MethodPost:
		time := time.Now().Format(FORMAT)

		data.User = models.User{
			Login:           r.FormValue("login"),
			Email:           r.FormValue("email"),
			Password:        r.FormValue("password"),
			ConfirmPassword: r.FormValue("confirm"),
			Created:         time,
		}

		if strings.TrimSpace(data.User.Email) == "" || strings.TrimSpace(data.User.Login) == "" {
			app.badRequest(w)
			return
		}

		if len(data.User.Login) > 16 {
			app.badRequest(w)
			return
		}

		if len(data.User.Password) < 6 || len(data.User.ConfirmPassword) < 6 {
			app.badRequest(w)
			return
		}

		// Valid characters for the login
		loginConvention := "^[a-zA-Z0-9]*$"
		// Valid characters for the email
		emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,4}$`)

		if re, _ := regexp.Compile(loginConvention); !re.MatchString(data.User.Login) {
			app.badRequest(w)
			return
		}

		_, err := mail.ParseAddress(data.User.Email)
		if !emailRegex.MatchString(data.User.Email) || err != nil {
			app.badRequest(w)
			return
		}

		if data.User.Password != data.User.ConfirmPassword {
			app.badRequest(w)
			return
		}

		if data.User.Login == "" || data.User.Email == "" || data.User.Password == "" || data.User.ConfirmPassword == "" {
			app.badRequest(w)
			return
		}

		data.User.Password = fmt.Sprintf(Salt, data.User.Password)

		hashPass, err := bcrypt.GenerateFromPassword([]byte(data.User.Password), 14)
		if err != nil {
			app.serverError(w, err)
			return
		}

		data.User.Password = string(hashPass)
		if (templateData{}.Errors) == data.Errors {
			data.User.ID, err = app.Forum.CreateUser(&data.User)
			if err != nil {
				data.Errors.IsAlreadyExist = true
			}
		}

		if (templateData{}.Errors) != data.Errors {
			app.render(w, r, "signup.page.html", data)
			return
		}

		addCookie(w, r, data.User.Login)
		http.Redirect(w, r, "/", http.StatusSeeOther)
	default:
		app.methodNotAllowed(w)
	}
}

func (app *Application) signin(w http.ResponseWriter, r *http.Request) {
	if isSession(r) {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	data := &templateData{IsSession: isSession(r)}

	switch r.Method {
	case http.MethodGet:
		app.render(w, r, "signin.page.html", data)
	case http.MethodPost:
		login := r.FormValue("login")
		password := fmt.Sprintf(Salt, r.FormValue("password"))

		if strings.TrimSpace(login) == "" || strings.TrimSpace(password) == "" {
			app.badRequest(w)
			return
		}

		user, err := app.Forum.GetUser(login)
		if err != nil {
			data.Errors.IsInvalidLoginOrPassword = true
		} else {
			err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
			if err != nil {
				data.Errors.IsInvalidLoginOrPassword = true
			}
		}

		if (templateData{}.Errors) != data.Errors {
			app.render(w, r, "signin.page.html", data)
			return
		}

		addCookie(w, r, user.Login)
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	default:
		app.methodNotAllowed(w)
	}
}

func (app *Application) signout(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.Header().Set("Allow", http.MethodGet)
		app.methodNotAllowed(w)
		return
	}

	if !isSession(r) {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	c, _ := r.Cookie(cookieName)

	http.SetCookie(w, &http.Cookie{
		Name:     cookieName,
		Value:    "",
		Path:     "/",
		HttpOnly: true,
		Expires:  time.Unix(0, 0),
	})

	cookie.Delete(c.Value)

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (app *Application) profile(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.Header().Set("Allow", http.MethodGet)
		app.methodNotAllowed(w)
		return
	}

	login := r.URL.Query().Get("login")
	if !isSession(r) && login == "" {
		http.Redirect(w, r, "/user/signin", http.StatusSeeOther)
		return
	}

	if isSession(r) && login == "" {
		login = getLogin(r)
		http.Redirect(w, r, "/user/profile?login="+login, http.StatusSeeOther)
		return
	}

	user, err := app.Forum.GetUser(login)
	if err != nil {
		app.notFound(w)
		return
	}

	posts, err := app.Forum.GetAllPosts()
	if err != nil {
		app.serverError(w, err)
		return
	}

	var userPosts []models.Post

	for _, post := range posts {
		if post.UserLogin == login {
			userPosts = append(userPosts, post)
		}
	}

	app.render(w, r, "profile.page.html", &templateData{
		User:       *user,
		IsSession:  isSession(r),
		Posts:      userPosts,
		NumOfPosts: len(userPosts),
	})
}

func (app *Application) createPost(w http.ResponseWriter, r *http.Request) {
	if !isSession(r) {
		http.Redirect(w, r, "/user/signin", http.StatusSeeOther)
		return
	}

	data := &templateData{IsSession: isSession(r)}

	switch r.Method {
	case http.MethodGet:
		app.render(w, r, "create.post.page.html", data)
	case http.MethodPost:
		time := time.Now().Format(FORMAT)

		data.Post = models.Post{
			Title:   r.FormValue("title"),
			Text:    r.FormValue("text"),
			Tags:    strings.Split(r.FormValue("tags"), " "),
			Created: time,
		}

		if strings.TrimSpace(data.Post.Title) == "" || strings.TrimSpace(data.Post.Text) == "" {
			app.badRequest(w)
			return
		}

		if utf8.RuneCountInString(data.Post.Title) > 58 || utf8.RuneCountInString(data.Post.Text) > 10000 {
			app.badRequest(w)
			return
		}

		if len(data.Post.Tags) > 10 {
			app.badRequest(w)
			return
		}

		for _, tag := range data.Post.Tags {
			if strings.Contains(tag, " ") || utf8.RuneCountInString(tag) > 16 {
				app.badRequest(w)
				return
			}
		}

		login := getLogin(r)

		user, err := app.Forum.GetUser(login)
		if err != nil {
			app.serverError(w, err)
			return
		}

		if (data.Post.Text == "" || data.Post.Title == "") || (len(data.Post.Tags) == 1 && data.Post.Tags[0] == "") {
			app.badRequest(w)
			return
		}

		data.Post.UserID = user.ID
		data.Post.UserLogin = login

		postID, err := app.Forum.CreatePost(&data.Post)
		if err != nil {
			app.serverError(w, err)
			return
		}
		err = app.Forum.CreateTags(data.Post.Tags)
		if err != nil {
			app.serverError(w, err)
			return
		}

		err = app.Forum.LinkTagsToAPost(postID, data.Post.Tags)
		if err != nil {
			app.serverError(w, err)
			return
		}

		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	default:
		app.methodNotAllowed(w)
	}
}

func (app *Application) createComment(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	if !isSession(r) {
		http.Redirect(w, r, "/user/signin", http.StatusSeeOther)
		return
	}

	postID, err := strconv.Atoi(r.URL.Query().Get("post"))
	if err != nil {
		app.notFound(w)
		return
	}

	login := getLogin(r)

	user, err := app.Forum.GetUser(login)
	if err != nil {
		app.serverError(w, err)
		return
	}

	time := time.Now().Format(FORMAT)
	comment := &models.Comment{
		PostID:  int64(postID),
		UserID:  user.ID,
		Login:   user.Login,
		Text:    r.FormValue("text"),
		Created: time,
	}

	if strings.TrimSpace(comment.Text) == "" {
		app.badRequest(w)
		return
	}

	if utf8.RuneCountInString(comment.Text) > 200 {
		app.badRequest(w)
		return
	}

	if comment.Text == "" {
		app.badRequest(w)
		return
	}

	err = app.Forum.CreateComment(comment)
	if err != nil {
		app.serverError(w, err)
		return
	}

	url := "/post?id=%v"
	http.Redirect(w, r, fmt.Sprintf(url, comment.PostID), http.StatusSeeOther)
}

func (app *Application) post(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.Header().Set("Allow", http.MethodGet)
		app.methodNotAllowed(w)
		return
	}

	key := r.URL.Query().Get("id")

	postID, err := strconv.Atoi(key)
	if err != nil {
		app.notFound(w)
		return
	}

	posts, err := app.Forum.GetAllPosts()
	if err != nil {
		app.serverError(w, err)
		return
	}

	data := &templateData{IsSession: isSession(r)}

	var isFound bool
	for _, post := range posts {
		if post.ID == int64(postID) {
			isFound = true
			data.Post = post
		}
	}

	if !isFound {
		app.notFound(w)
		return
	}

	app.render(w, r, "post.id.page.html", data)
}

func (app *Application) postVote(w http.ResponseWriter, r *http.Request) {
	if !isSession(r) {
		http.Redirect(w, r, "/user/signin", http.StatusSeeOther)
		return
	}

	if r.Method != http.MethodPost {
		w.Header().Set("Allow", http.MethodPost)
		app.methodNotAllowed(w)
		return
	}

	postID, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil {
		app.badRequest(w)
		return
	}

	vote, err := strconv.Atoi(r.URL.Query().Get("vote"))
	if err != nil {
		app.badRequest(w)
		return
	}

	if vote != 1 && vote != -1 {
		app.badRequest(w)
		return
	}

	login := getLogin(r)

	user, err := app.Forum.GetUser(login)
	if err != nil {
		app.badRequest(w)
		return
	}

	_, err = app.Forum.GetPostByID(postID)
	if err != nil {
		app.badRequest(w)
		return
	}

	err = app.Forum.AddVoteToPost(user.ID, int64(postID), vote)
	if err != nil {
		app.serverError(w, err)
		return
	}

	url := "/post?id=%v"
	http.Redirect(w, r, fmt.Sprintf(url, postID), http.StatusSeeOther)
}

func (app *Application) commentVote(w http.ResponseWriter, r *http.Request) {
	if !isSession(r) {
		http.Redirect(w, r, "/user/signin", http.StatusSeeOther)
		return
	}

	if r.Method != http.MethodPost {
		w.Header().Set("Allow", http.MethodPost)
		app.methodNotAllowed(w)
		return
	}

	commentID, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil {
		fmt.Println(513)
		app.badRequest(w)
		return
	}

	vote, err := strconv.Atoi(r.URL.Query().Get("vote"))
	if err != nil {
		fmt.Println(520)
		app.badRequest(w)
		return
	}

	if vote != 1 && vote != -1 {
		fmt.Println(526)
		app.badRequest(w)
		return
	}

	login := getLogin(r)

	user, err := app.Forum.GetUser(login)
	if err != nil {
		fmt.Println(535)
		app.badRequest(w)
		return
	}

	comment, err := app.Forum.GetCommentByID(int64(commentID))
	if err == sql.ErrNoRows {
		app.ErrorLog.Println(err)
		app.badRequest(w)
		return
	}
	if err != nil {
		app.serverError(w, err)
		return
	}

	err = app.Forum.AddVoteToComment(user.ID, int64(commentID), vote)
	if err != nil {
		app.serverError(w, err)
		return
	}

	url := "/post?id=%v"
	http.Redirect(w, r, fmt.Sprintf(url, comment.PostID), http.StatusSeeOther)
}

func (app *Application) filter(w http.ResponseWriter, r *http.Request) {
	isSession := isSession(r)

	tag := r.URL.Query().Get("tag")

	posts, err := app.Forum.GetAllPosts()
	if err != nil {
		app.serverError(w, err)
		return
	}

	tags, err := app.Forum.GetAllTags()
	if err != nil {
		app.serverError(w, err)
		return
	}

	tagsPosts := posts
	if tag != "" {
		tagsPosts, err = app.filterByTag(tag, posts)
		if err != nil {
			app.notFound(w)
			return
		}
	} else {
		app.notFound(w)
		return
	}

	app.render(w, r, "filter.tag.page.html", &templateData{
		Posts:     tagsPosts,
		IsSession: isSession,
		Tags:      tags,
	})
}

func (app *Application) filterByTag(tag string, posts []models.Post) ([]models.Post, error) {
	tagsPosts := []models.Post{}

	for _, post := range posts {
		for _, postTag := range post.Tags {
			if postTag == tag {
				tagsPosts = append(tagsPosts, post)
			}
		}
	}

	if len(tagsPosts) == 0 {
		return nil, fmt.Errorf("not found")
	}
	return tagsPosts, nil
}

func (app *Application) likedPosts(w http.ResponseWriter, r *http.Request) {
	if !isSession(r) {
		http.Redirect(w, r, "/user/signin", http.StatusSeeOther)
		return
	}

	if r.Method != http.MethodGet {
		app.methodNotAllowed(w)
		return
	}

	login := getLogin(r)

	posts, err := app.Forum.GetVotePosts(login, 1)
	if err != nil {
		app.serverError(w, err)
		return
	}

	tags, err := app.Forum.GetAllTags()
	if err != nil {
		app.serverError(w, err)
		return
	}

	app.render(w, r, "home.page.html", &templateData{
		Posts:     *posts,
		IsSession: true,
		Tags:      tags,
	})
}

func (app *Application) dislikedPosts(w http.ResponseWriter, r *http.Request) {
	if !isSession(r) {
		http.Redirect(w, r, "/user/signin", http.StatusSeeOther)
		return
	}

	if r.Method != http.MethodGet {
		app.methodNotAllowed(w)
		return
	}

	login := getLogin(r)

	posts, err := app.Forum.GetVotePosts(login, -1)
	if err != nil {
		app.serverError(w, err)
	}

	tags, err := app.Forum.GetAllTags()
	if err != nil {
		app.serverError(w, err)
		return
	}

	app.render(w, r, "home.page.html", &templateData{
		Posts:     *posts,
		IsSession: true,
		Tags:      tags,
	})
}
