package web

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/Li-Khan/forum/pkg/models"
	"golang.org/x/crypto/bcrypt"
)

// FORMAT - time format
const FORMAT string = "01-02-2006 15:04:05"

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

		if data.User.Password != data.User.ConfirmPassword {
			data.Errors.IsError = true
			data.Errors.IsPassNotMatch = true
		}

		if data.User.Login == "" || data.User.Email == "" || data.User.Password == "" || data.User.ConfirmPassword == "" {
			data.Errors.IsError = true
			data.Errors.IsInvalidForm = true
		}

		hashPass, err := bcrypt.GenerateFromPassword([]byte(data.User.Password), 14)
		if err != nil {
			// error handle ...
		}

		data.User.Password = string(hashPass)
		data.User.ID, err = app.Snippet.CreateUser(&data.User)

		if err != nil {
			// error handle ...
		}

		if data.Errors.IsError {
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

	switch r.Method {
	case http.MethodGet:
		// Обработка страницы
		app.render(w, r, "signin.page.html", &templateData{})
	case http.MethodPost:
		// Получение данных
		user := models.User{
			Login:    r.FormValue("login"),
			Password: r.FormValue("password"),
		}
		fmt.Println(user)
		addCookie(w, r, user.Login)

		// Обработка данных

		// Добавление данных в бд

		// Перенаправление
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
		http.Redirect(w, r, "/user/signin", http.StatusSeeOther)
		return
	}

	if !isSession(r) {
		// Если пользователь не в сессии
		// обработать это ...
	}

	// получение данных

	// обработка данных

	// добавление данных в бд

	// перенаправление
	http.Redirect(w, r, "/", http.StatusSeeOther)
}
