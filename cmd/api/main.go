package main

import (
	"expvar"
	"flag"
	"fmt"
	"net/http"
	"os"
	"runtime"
	"strconv"
	"sync"
	"time"

	log "mooveit-backend.mooveit.com/internal/jsonlog"
	"mooveit-backend.mooveit.com/internal/vcs"
)

var version = vcs.Version()

type appConfig struct {
	port int
	env  string
}

type application struct {
	config appConfig
	wg     sync.WaitGroup // Include a sync.WaitGroup in the application struct. The zero-value for a sync.WaitGroup type is a valid, useable, sync.WaitGroup with a 'counter' value of 0, so we don't need to do anything else to initialize it before we can use it.
}

func main() {
	// Declare an instance of the appConfig struct.
	var cfg appConfig
	parseFlags(&cfg)

	// Set metrics parameters for the debug/vars endpoint
	setMetricsParameters()

	// Declare an instance of the application struct, containing the appConfig struct and the log.
	app := &application{
		config: cfg,
	}

	// Start the server
	err := app.serve()
	if err != nil {
		log.Fatal(err)
	}
}

func parseFlags(cfg *appConfig) {
	// Read the command-line flags into the appConfig struct
	// Server
	// Default port is 4000, but check for PORT environment variable first (Railway requirement)
	defaultPort := 4000
	if portEnv := os.Getenv("PORT"); portEnv != "" {
		if port, err := strconv.Atoi(portEnv); err == nil {
			defaultPort = port
		}
	}
	flag.IntVar(&cfg.port, "port", defaultPort, "API server port")

	// Default environment is development, but check for ENV environment variable
	defaultEnv := "development"
	if envEnv := os.Getenv("ENV"); envEnv != "" {
		defaultEnv = envEnv
	}
	flag.StringVar(&cfg.env, "env", defaultEnv, "Environment (development|staging|production)")

	// Create a new version boolean flag with the default value of false.
	displayVersion := flag.Bool("version", false, "Display version and exit")

	flag.Parse()
	log.Info("parseFlags() - command-line flags have been parsed")

	// If the version flag value is true, then print out the version number and
	// immediately exit.>
	if *displayVersion {
		log.Info("Version:\t%s", version)
		os.Exit(0)
	}
}

func setMetricsParameters() {
	// Publish a new "version" variable in the expvar handler containing our application
	// version number (currently the constant "1.0.0").
	expvar.NewString("version").Set(version)

	// Publish the number of active goroutines.
	expvar.Publish("goroutines", expvar.Func(func() any {
		return runtime.NumGoroutine()
	}))

	// Publish the current Unix timestamp.
	expvar.Publish("timestamp", expvar.Func(func() any {
		return time.Now().Unix()
	}))
}

func (app *application) serve() error {
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", app.config.port),
		Handler: app.routes(),
	}

	log.Info("Starting server on port %d", app.config.port)
	return srv.ListenAndServe()
}
