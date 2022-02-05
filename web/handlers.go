package web

import (
	"fmt"
	"net/http"
	"net/mail"
	"regexp"
	"strconv"
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

	app.render(w, r, "home.page.html", &templateData{})
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
		}

		err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
		if err != nil {
			data.Errors.IsInvalidLoginOrPassword = true
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

	cookie.mx.Lock()
	defer cookie.mx.Unlock()
	// r.Cookie - не вернет ошибку потому что перед тем как вызвать deleteCookie
	// вызывается isSession который уже проверяет его на ошибку
	c, _ := r.Cookie(cookieName)

	http.SetCookie(w, &http.Cookie{
		Name:     cookieName,
		Value:    "",
		Path:     "/",
		HttpOnly: true,
		Expires:  time.Unix(0, 0),
	})

	delete(cookie.mapCookie, c.Value)

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (app *Application) profile(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.Header().Set("Allow", http.MethodGet)
		app.methodNotAllowed(w)
		return
	}

	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil || id < 1 {
		app.notFound(w)
		return
	}

	app.render(w, r, "profile.page.html", &templateData{})
}

func (app *Application) createPost(w http.ResponseWriter, r *http.Request) {
	if !isSession(r) {
		http.Redirect(w, r, "/user/signin", http.StatusSeeOther)
		return
	}

	switch r.Method {
	case http.MethodGet:
		// Обработка страницы
		app.render(w, r, "create.post.page.html", &templateData{})
	case http.MethodPost:
		// Получение данных
		title := r.FormValue("title")
		text := r.FormValue("text")
		fmt.Println(title, text)

		// Обработка данных

		// Добавление данных в бд

		// Перенаправление
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

	// получение данных

	// обработка данных

	// добавление данных в бд

	// перенаправление
	http.Redirect(w, r, "/", http.StatusSeeOther)
}
