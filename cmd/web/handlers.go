package main

import (
	"fmt"
	"html/template"
	"net/http"
	"strconv"
)

// * Setting the handlers as a application method we can use the application dependencies injected
// on its properties
func (app *application) home(w http.ResponseWriter, r *http.Request) {
	// To avoid the servemux behavior for subtree paths,
	// check if the path is exactly the same to continue or redirect
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	// Parse a template file using the absoluty path or the relative path from
	// the path that you start the application (usually the root)
	templatefilesPaths := []string{
		"./ui/html/base.tmpl",
		"./ui/html/partials/nav.tmpl",
		"./ui/html/pages/home.tmpl",
	}
	// The function "template.ParseFiles" accepts many template files paths (variadic function)
	ts, err := template.ParseFiles(templatefilesPaths...)
	if err != nil {
		// Using the application dependency `errorLog` injected on its initialization in main()
		app.errorLog.Println(err.Error())
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// * Write a response using the template file passing the dynamic data used inside template
	// * ts.Execute(w, dynamicData) to not use named template
	// * ts.ExecuteTemplate(w, "templateName", dynamicData) to use a named template
	// * Instead read files from disk using the package "html/template" another way is using the
	// package "embed" to embed files in Go
	var dynamicData interface{} = nil
	err = ts.ExecuteTemplate(w, "base", dynamicData)
	if err != nil {
		app.errorLog.Println(err.Error())
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}

func (app *application) snippetView(w http.ResponseWriter, r *http.Request) {
	// Getting query parameters
	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil || id < 1 {
		http.NotFound(w, r)
		return
	}

	// Fprintf can receive any instance that implements the io.Writer interface
	// working as a helper for write responses
	fmt.Fprintf(w, "Display a specific snippet with id %d", id)
}

func (app *application) snippetCreate(w http.ResponseWriter, r *http.Request) {
	// * Always use http constants as the best practice: http.MethodPost = "POST".
	// See all http constants in https://pkg.go.dev/net/http#pkg-constants
	if r.Method != http.MethodPost {
		// * You must set headers, if you need, before call functions to write response (w.WriteHeader(...),
		// w.Write(...), http.Error(...), fmt.Fprintf(w,...), etc) to have effect.
		w.Header().Set("Allow", http.MethodPost)

		// * For errors use the helper "http.Error(w,...)" to call "w.WriteHeader" and "w.Write"
		// indirectly to write the response;
		// * Always use http constants as the best practice: http.StatusMethodNotAllowed = 405.
		// See all http constants in https://pkg.go.dev/net/http#pkg-constants
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	// * Go automatically set the headers "Content-type", "Content-Length" and "Date", but
	// the function http.DetectContentType() can't distinguish JSON from plain text, so
	// you need to set manually.
	w.Header().Set("Content-type", "application/json")

	// To supress system generated headers you need to set "nil" for the value
	w.Header()["Date"] = nil

	// * Instead use "Set" you can use "Add" to append values for the same header
	// * When editing the header throught Set(), Add(), Del(), Get() and Values(), it
	// will be automatically canonicalized (first letters in keys names putted in upper case for
	// HTTP/1 or all keys and values in down case for HTTP/2)
	// Result: "Cache-Control": "public,mac-age=32645000"
	w.Header().Add("cache-control", "public")
	w.Header().Add("Cache-Control", "mac-age=32645000")

	// * Modifying Header directly, keep the the same case (upper or lower)
	w.Header()["key-in-down-case"] = []string{"keeps-in-down-case"}

	// * The function "w.WriteHeader" should be called only once per response;
	// * When not called after the function "w.Write", the response status code will be 200;
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("Create a new snippet..."))

}

// * Passing dependencies from main() directly to the handler function using "closure"
func closureExample(app *application) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		app.infoLog.Println("Handler using closure as a example")
		w.Write([]byte("Handler using closure as a example."))
	})
}
