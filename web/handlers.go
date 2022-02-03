package web

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/Li-Khan/forum/pkg/models"
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

	switch r.Method {
	case http.MethodGet:
		// Обработка страницы
		app.render(w, r, "signup.page.html", &templateData{})
	case http.MethodPost:
		// Получение данных
		time := time.Now().Format(FORMAT)
		user := &models.User{
			Login:           r.FormValue("login"),
			Email:           r.FormValue("email"),
			Password:        r.FormValue("password"),
			ConfirmPassword: r.FormValue("confirm"),
			Created:         time,
		}

		// data := &templateData{User: user}
		fmt.Println(user)
		// app.Snippet.InsertUser(user)
		// Обработка данных

		// Добавление данных в бд

		// Создание куки
		addCookie(w, r, user.Login)

		// Перенаправление
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
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
	}

	deleteCookie(w, r)

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
