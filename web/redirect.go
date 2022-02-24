package web

import "net/http"

func (app *Application) userRedirect(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/user/" {
		app.notFound(w)
		return
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (app *Application) createRedirect(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/create/" {
		app.notFound(w)
		return
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}
