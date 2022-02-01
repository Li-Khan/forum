package web

import (
	"fmt"
	"html/template"
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

	files := []string{
		"./ui/html/home.page.html",
		"./ui/html/base.layout.html",
		"./ui/html/footer.partial.html",
	}

	ts, err := template.ParseFiles(files...)
	if err != nil {
		app.ErrorLog.Println(err.Error())
		app.serverError(w, err)
		return
	}

	err = ts.Execute(w, nil)
	if err != nil {
		app.ErrorLog.Println(err.Error())
		app.serverError(w, err)
	}
}

func (app *Application) signup(w http.ResponseWriter, r *http.Request) {
	files := []string{
		"./ui/html/signup.page.html",
		"./ui/html/base.layout.html",
		"./ui/html/footer.partial.html",
	}

	ts, err := template.ParseFiles(files...)
	if err != nil {
		app.ErrorLog.Println(err.Error())
		app.serverError(w, err)
		return
	}

	switch r.Method {
	case http.MethodGet:
		// Обработка страницы
		err = ts.Execute(w, nil)
		if err != nil {
			app.ErrorLog.Println(err.Error())
			app.serverError(w, err)
			return
		}
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
	files := []string{
		"./ui/html/signin.page.html",
		"./ui/html/base.layout.html",
		"./ui/html/footer.partial.html",
	}

	ts, err := template.ParseFiles(files...)
	if err != nil {
		app.ErrorLog.Println(err.Error())
		app.serverError(w, err)
		return
	}

	switch r.Method {
	case http.MethodGet:
		// Обработка страницы
		err = ts.Execute(w, nil)
		if err != nil {
			app.ErrorLog.Println(err.Error())
			app.serverError(w, err)
			return
		}
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

	files := []string{
		"./ui/html/profile.page.html",
		"./ui/html/base.layout.html",
		"./ui/html/footer.partial.html",
	}

	ts, err := template.ParseFiles(files...)
	if err != nil {
		app.ErrorLog.Println(err.Error())
		app.serverError(w, err)
		return
	}

	err = ts.Execute(w, nil)
	if err != nil {
		app.ErrorLog.Println(err.Error())
		app.serverError(w, err)
	}
}

func (app *Application) createPost(w http.ResponseWriter, r *http.Request) {
	files := []string{
		"./ui/html/create.post.page.html",
		"./ui/html/base.layout.html",
		"./ui/html/footer.partial.html",
	}

	ts, err := template.ParseFiles(files...)
	if err != nil {
		app.ErrorLog.Println(err.Error())
		app.serverError(w, err)
		return
	}

	switch r.Method {
	case http.MethodGet:
		// Обработка страницы
		err = ts.Execute(w, nil)
		if err != nil {
			app.ErrorLog.Println(err.Error())
			app.serverError(w, err)
			return
		}
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
