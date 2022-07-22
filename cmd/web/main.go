package main

import (
	"flag"
	"log"
	"net/http"
	"os"
)

// * Create a application struct to inject dependencies needed for the whole application
// * For handlers spread in different packages we can pass dependencies directly to handler
// functions using closure (ex: https://www.calhoun.io/5-useful-ways-to-use-closures-in-go)
// or creating a handler struct for each package to receive dependencies to be accessed by
// the handler functions as methods
type application struct {
	infoLog  *log.Logger
	errorLog *log.Logger
}

func main() {
	// * Map a flag passed from commands in terminal to a variable, storing the default value,
	// such as `go run ./cmd/web -addr=":9999"`
	// * Return a pointer
	// * flag.String(), flagInt(), flag.Bool(), flag.Float64() try to convert value automatically
	// * The default for boolean flags is true
	// * When use the flag "-help" on the terminal all the mapped flags will be showed
	// * To map flags in pre-existing variables use `flag.StringVar(&addr, "addr", ":4000", "...")`
	// * Instead flags you can use env variables to get values with `os.Getenv()`, but for this you
	// don't have implicity type conversion, default value and help like flags. Another way is to use
	// flags with envs to have all the functionalities. Ex: `go run ./cmd/web -addr=$MY_ADDR_ENV`
	addr := flag.String("addr", ":4000", "HTTP network address")
	// * Parse flags to set the received values in mapped variables
	flag.Parse()

	// * With `log.New(destination, prefix, additional information)` we can set different log levels
	// for diferent destinations and different additional information (flags)
	// * Use log constants (https://pkg.go.dev/log#pkg-constants) to set additional information like
	// `log.Ldate` to add date, `log.Ltime` to add time, `log.Lshortfile` to add filename and line,
	// log.LUTC to use UTC datetime instead local datetime, etc...
	// * Loggers created by log.New() are concurrency-safe, but if you have many loggers writing in the
	// same destination you need to ensure that the destination is safe for concurrent use too
	// * Output: INFO    2022/07/19 19:32:12 my custom message
	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	// * Output: ERROR   2022/07/19 19:32:12 main.go:85: error message
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	// Instantiate the application injecting dependencies
	app := &application{
		infoLog:  infoLog,
		errorLog: errorLog,
	}

	// The function http.ListenAndServe() should receive the TCP network
	// in the format "host:port" to specifics host and port number or
	// ":port" to any host and specific port number or
	// ":my-port-name" to use named port that Go will try to get the
	// correspondence from "/etc/services"
	infoLog.Println("Starting server on", *addr)
	// * The servemux implements the http.Handler interface, so we can pass it to the
	// http.ListenAndServe and for each request the handler method in servemux forward
	// to the correct handler method based in the registered routes
	// * Each request is handled concurrently in its own goroutine
	// * To use the custom errorLog for http errors we should to instantiate a new http
	// server struct passing the settings by properties and start the server from it
	srv := &http.Server{
		Addr:     *addr,
		Handler:  app.routes(),
		ErrorLog: errorLog,
	}
	err := srv.ListenAndServe()
	// * Is a good practice call Panic() ou Fatal() and their variations only from the main()
	// * More log methods good to know: https://pkg.go.dev/log#Logger
	errorLog.Fatal(err)
}
