package web

import (
	"fmt"
	"net/http"
)

func (app *Application) serverError(w http.ResponseWriter, err error) {
	// trace - получаю трассировку стека для текущей горутины
	// trace := fmt.Sprintf("%s\n%s", err.Error(), debug.Stack())
	// app.ErrorLog.Output(2, trace)
	app.ErrorLog.Println(err)

	app.render(w, nil, "error.page.html", &templateData{Error: http.StatusText(http.StatusInternalServerError)})
}

// func (app *Application) clientError(w http.ResponseWriter, status int) {
// 	w.WriteHeader(status)
// 	app.render(w, nil, "error.page.html", &templateData{Error: http.StatusText(status)})
// 	// http.Error(w, http.StatusText(status), status)
// }

func (app *Application) notFound(w http.ResponseWriter) {
	w.WriteHeader(http.StatusNotFound)
	app.render(w, nil, "error.page.html", &templateData{Error: http.StatusText(http.StatusNotFound)})
	// app.clientError(w, http.StatusNotFound)
}

func (app *Application) badRequest(w http.ResponseWriter) {
	w.WriteHeader(http.StatusBadRequest)
	app.render(w, nil, "error.page.html", &templateData{Error: http.StatusText(http.StatusBadRequest)})
	// app.clientError(w, http.StatusBadRequest)
}

func (app *Application) methodNotAllowed(w http.ResponseWriter) {
	w.WriteHeader(http.StatusMethodNotAllowed)
	app.render(w, nil, "error.page.html", &templateData{Error: http.StatusText(http.StatusMethodNotAllowed)})
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
