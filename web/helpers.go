package web

import (
	"fmt"
	"net/http"
	"runtime/debug"
)

func (app *Application) serverError(w http.ResponseWriter, err error) {
	// trace - получаю трассировку стека для текущей горутины и добавляю ее в логер
	trace := fmt.Sprintf("%s\n%s", err.Error(), debug.Stack())
	app.ErrorLog.Output(2, trace)

	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}

func (app *Application) clientError(w http.ResponseWriter, status int) {
	w.WriteHeader(status)
	app.render(w, nil, "error.page.html", &templateData{Error: http.StatusText(status)})
	// http.Error(w, http.StatusText(status), status)
}

func (app *Application) notFound(w http.ResponseWriter) {
	w.WriteHeader(http.StatusNotFound)
	app.render(w, nil, "error.page.html", &templateData{Error: http.StatusText(404)})
	// app.clientError(w, http.StatusNotFound)
}

func (app *Application) badRequest(w http.ResponseWriter) {
	w.WriteHeader(http.StatusBadRequest)
	app.render(w, nil, "error.page.html", &templateData{Error: http.StatusText(400)})
	// app.clientError(w, http.StatusBadRequest)
}

func (app *Application) methodNotAllowed(w http.ResponseWriter) {
	w.WriteHeader(http.StatusMethodNotAllowed)
	app.render(w, nil, "error.page.html", &templateData{Error: http.StatusText(405)})
	// app.clientError(w, http.StatusMethodNotAllowed)
}

func (app *Application) render(w http.ResponseWriter, r *http.Request, name string, td *templateData) {
	ts, ok := app.TemplateCache[name]
	if !ok {
		app.serverError(w, fmt.Errorf("template %s does not exist", name))
		return
	}

	err := ts.Execute(w, td)
	if err != nil {
		app.serverError(w, err)
	}
}
