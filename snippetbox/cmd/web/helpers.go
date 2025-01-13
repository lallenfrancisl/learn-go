package main

import (
	"fmt"
	"net/http"
	"runtime/debug"
)

// Logs an error and sends a generic 500 Internal Server Error response
// to the user
func (app *Application) serverError(w http.ResponseWriter, r *http.Request, err error) {
	app.logger.Error(
		err.Error(),
		"method",
		r.Method,
		"url",
		r.URL.RequestURI(),
		"trace",
		debug.Stack(),
	)
	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}

// Sends a specific status code and corresponding description to the user
func (app *Application) clientError(w http.ResponseWriter, status int) {
	text := fmt.Sprintf("%d %s", status, http.StatusText(status))
	app.logger.Error(text)
	http.Error(w, text, status)
}

func (app *Application) render(
	w http.ResponseWriter, r *http.Request,
	status int, page string,
	data templateData,
) {
	ts, ok := app.templateCache[page]
	if !ok {
		err := fmt.Errorf("the template %s does not exist", page)
		app.serverError(w, r, err)

		return
	}

	w.WriteHeader(status)

	err := ts.ExecuteTemplate(w, "base", data)
	if err != nil {
		app.serverError(w, r, err)

		return
	}
}
