package web

import (
	"net/http"
)

// Routes - initializes routes
func (app *Application) Routes() *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("/", app.home)
	// user handlers
	mux.HandleFunc("/user/", userRedirect)
	mux.HandleFunc("/user/signin", app.signin)
	mux.HandleFunc("/user/signup", app.signup)
	mux.HandleFunc("/user/logout", app.signout)
	mux.HandleFunc("/user/profile", app.profile)

	// create handlers
	mux.HandleFunc("/create/", createRedirect)
	mux.HandleFunc("/create/post", app.createPost)
	mux.HandleFunc("/create/comment", app.createComment)

	fileServer := http.FileServer(neuteredFileSystem{http.Dir("./ui/static")})
	mux.Handle("/static", http.NotFoundHandler())
	mux.Handle("/static/", http.StripPrefix("/static", fileServer))

	return mux
}
