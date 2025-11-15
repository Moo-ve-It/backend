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
	// Log application startup
	log.Info("Application starting...")
	log.InfoWithProperties("Application version", map[string]string{
		"version": version,
	})

	// Declare an instance of the appConfig struct.
	var cfg appConfig
	parseFlags(&cfg)

	// Log configuration
	log.InfoWithProperties("Application configuration loaded", map[string]string{
		"environment": cfg.env,
		"port":        fmt.Sprintf("%d", cfg.port),
	})

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

	// Construct server URL based on environment
	serverURL := app.getServerURL()

	// Log detailed server startup information
	log.InfoWithProperties("Server starting", map[string]string{
		"port":        fmt.Sprintf("%d", app.config.port),
		"address":     fmt.Sprintf("0.0.0.0:%d", app.config.port),
		"url":         serverURL,
		"environment": app.config.env,
	})

	log.Info("Server is ready to accept connections")
	log.Info("Health check endpoint available at: %s/healthcheck", serverURL)
	log.Info("Metrics endpoint available at: %s/debug/vars", serverURL)

	return srv.ListenAndServe()
}

// getServerURL constructs the full server URL based on the deployment environment
func (app *application) getServerURL() string {
	// Check for Railway public domain (Railway sets this automatically)
	if railwayDomain := os.Getenv("RAILWAY_PUBLIC_DOMAIN"); railwayDomain != "" {
		return fmt.Sprintf("https://%s", railwayDomain)
	}

	// Check for Railway service URL
	if railwayServiceURL := os.Getenv("RAILWAY_STATIC_URL"); railwayServiceURL != "" {
		return railwayServiceURL
	}

	// Check for custom domain environment variable
	if customDomain := os.Getenv("PUBLIC_DOMAIN"); customDomain != "" {
		scheme := "https"
		if app.config.env == "development" {
			scheme = "http"
		}
		return fmt.Sprintf("%s://%s", scheme, customDomain)
	}

	// Default to localhost for development
	if app.config.env == "development" {
		return fmt.Sprintf("http://localhost:%d", app.config.port)
	}

	// For production without domain info, return generic URL
	return fmt.Sprintf("https://0.0.0.0:%d", app.config.port)
}
