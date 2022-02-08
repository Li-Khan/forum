package web

import (
	"fmt"
	"net/http"
	"net/mail"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/Li-Khan/forum/pkg/models"
	"golang.org/x/crypto/bcrypt"
)

const (
	// FORMAT - time format.
	FORMAT string = "01-02-2006 15:04:05"
	// Pass - adds additional characters to the user's password to make the password more complex.
	Pass string = "758+}+%s#=?.^$69,"
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

	data := &templateData{}

	posts, err := app.Snippet.GetAllPosts()
	if err != nil {
		app.serverError(w, err)
		return
	}
	data.Posts = *posts

	data.IsSession = isSession(r)

	app.render(w, r, "home.page.html", data)
}

func (app *Application) signup(w http.ResponseWriter, r *http.Request) {
	if isSession(r) {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	data := &templateData{}

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

		// Valid characters for the login
		loginConvention := "^[a-zA-Z0-9]*[-]?[a-zA-Z0-9]*$"
		// Valid characters for the email
		emailRegex := regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,4}$`)

		if re, _ := regexp.Compile(loginConvention); !re.MatchString(data.User.Login) {
			data.Errors.IsInvalidLogin = true
		}

		_, err := mail.ParseAddress(data.User.Email)
		if !emailRegex.MatchString(data.User.Email) || err != nil {
			data.Errors.IsInvalidEmail = true
		}

		if data.User.Password != data.User.ConfirmPassword {
			data.Errors.IsPassNotMatch = true
		}

		if data.User.Login == "" || data.User.Email == "" || data.User.Password == "" || data.User.ConfirmPassword == "" {
			data.Errors.IsInvalidForm = true
		}

		data.User.Password = fmt.Sprintf(Pass, data.User.Password)

		hashPass, err := bcrypt.GenerateFromPassword([]byte(data.User.Password), 14)
		if err != nil {
			app.serverError(w, err)
			return
		}

		data.User.Password = string(hashPass)
		if (templateData{}.Errors) == data.Errors {
			data.User.ID, err = app.Snippet.CreateUser(&data.User)
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

	data := &templateData{}

	switch r.Method {
	case http.MethodGet:
		app.render(w, r, "signin.page.html", data)
	case http.MethodPost:
		login := r.FormValue("login")
		password := fmt.Sprintf(Pass, r.FormValue("password"))

		user, err := app.Snippet.GetUser(login)
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
		c, _ := r.Cookie(cookieName)
		value, _ := cookie.Load(c.Value)
		login = fmt.Sprint(value)
		http.Redirect(w, r, "/user/profile?login="+login, http.StatusSeeOther)
		return
	}

	user, err := app.Snippet.GetUser(login)
	if err != nil {
		app.notFound(w)
		return
	}

	app.render(w, r, "profile.page.html", &templateData{User: *user})
}

func (app *Application) createPost(w http.ResponseWriter, r *http.Request) {
	if !isSession(r) {
		http.Redirect(w, r, "/user/signin", http.StatusSeeOther)
		return
	}

	data := &templateData{}

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

		c, _ := r.Cookie(cookieName)
		value, _ := cookie.Load(c.Value)
		login := fmt.Sprint(value)

		user, err := app.Snippet.GetUser(login)
		if err != nil {
			app.serverError(w, err)
			return
		}

		for _, tag := range data.Post.Tags {
			if strings.Contains(tag, " ") {
				data.Errors.IsInvalidForm = true
			}
		}

		if data.Post.Text == "" || data.Post.Title == "" || len(data.Post.Tags) == 1 {
			data.Errors.IsInvalidForm = true
		} else {
			data.Post.UserID = user.ID
			data.Post.UserLogin = login

			postID, err := app.Snippet.CreatePost(&data.Post)
			if err != nil {
				app.serverError(w, err)
				return
			}
			err = app.Snippet.CreateTags(data.Post.Tags)
			if err != nil {
				app.serverError(w, err)
				return
			}

			err = app.Snippet.CreatePostsAndTags(postID, data.Post.Tags)
			if err != nil {
				fmt.Println(err)
			}
		}

		if data.Errors.IsInvalidForm {
			app.render(w, r, "create.post.page.html", data)
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

	c, _ := r.Cookie(cookieName)
	value, _ := cookie.Load(c.Value)
	login := fmt.Sprint(value)

	user, err := app.Snippet.GetUser(login)
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

	if comment.Text == "" {
		app.badRequest(w)
		return
	}

	err = app.Snippet.CreateComment(comment)
	if err != nil {
		app.serverError(w, err)
		return
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}
