package main

import (
	"fmt"
	"net/http"
	"runtime/debug"
)

// Helper function to centralize error handling logging stack trace
func (app *application) serverError(w http.ResponseWriter, err error) {
	// debug.Stack() get the stack trace for the current goroutine
	trace := fmt.Sprintf("%s\n%s", err.Error(), debug.Stack())
	// * Using the application dependency `errorLog` injected on its initialization in main()
	// * Using Output() we can show the correct error origin ("handler.go:32") instead the
	// current file ("helpers.go:16") passing the correct number to back steps
	app.errorLog.Output(2, trace)

	// * For errors use the helper "http.Error(w,...)" to call "w.WriteHeader" and "w.Write"
	// indirectly to write the response;
	// * Always use http constants as the best practice: http.StatusMethodNotAllowed = 405.
	// See all http constants in https://pkg.go.dev/net/http#pkg-constants
	// * http.StatusText() return the status code description
	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}

// Helper function to centralize errors handling
func (app *application) clientError(w http.ResponseWriter, status int) {
	http.Error(w, http.StatusText(status), status)
}

// Helper function to centralize errors handling
func (app *application) notFound(w http.ResponseWriter) {
	app.clientError(w, http.StatusNotFound)
}
