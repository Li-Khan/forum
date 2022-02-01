package web

import "net/http"

func userRedirect(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/user/" {
		http.NotFound(w, r)
		return
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func createRedirect(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/create/" {
		http.NotFound(w, r)
		return
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}
