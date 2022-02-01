package web

import "net/http"

// Routes - initializes routes
func (app *Application) Routes() *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("/", app.home)
	// user handlers
	mux.HandleFunc("/user/", redirect)
	mux.HandleFunc("/user/signin", app.signin)
	mux.HandleFunc("/user/signup", app.signup)
	mux.HandleFunc("/user/profile", app.profile)
	mux.HandleFunc("/user/logout", app.logout)

	// create handlers
	mux.HandleFunc("/create/", redirect)
	mux.HandleFunc("/create/post", app.createPost)
	mux.HandleFunc("/create/comment", app.createComment)

	return mux
}
