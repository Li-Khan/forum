package web

import (
	"fmt"
	"net/http"
)

func (app *Application) home(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	if r.Method != http.MethodGet {
		w.Header().Set("Allow", http.MethodGet)
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return
	}

	fmt.Fprint(w, "home page")
}

func (app *Application) signin(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "signin page")
}

func (app *Application) signup(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "signup page")
}

func (app *Application) signout(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "signout page")
}

func (app *Application) profile(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "profile page")
}

func (app *Application) createPost(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "create post page")
}

func (app *Application) createComment(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "create comment page")
}
