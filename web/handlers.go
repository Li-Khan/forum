package web

import (
	"fmt"
	"net/http"
	"strconv"
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
	switch r.Method {
	case http.MethodGet:
		// Обработка страницы
		app.render(w, r, "signup.page.html", &templateData{})
	case http.MethodPost:
		// Получение данных
		login := r.FormValue("login")
		email := r.FormValue("email")
		password := r.FormValue("password")
		confirm := r.FormValue("confirm")
		fmt.Println(login, email, password, confirm)

		// Обработка данных

		// Добавление данных в бд

		// Перенаправление
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	default:
		app.methodNotAllowed(w)
	}
}

func (app *Application) signin(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		// Обработка страницы
		app.render(w, r, "signin.page.html", &templateData{})
	case http.MethodPost:
		// Получение данных
		login := r.FormValue("login")
		password := r.FormValue("password")
		fmt.Println(login, password)

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
	// Обработка выхода
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

	// получение данных

	// обработка данных

	// добавление данных в бд

	// перенаправление
	http.Redirect(w, r, "/", http.StatusSeeOther)
}
