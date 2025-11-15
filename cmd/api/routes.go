package main

import (
	"expvar"
	"fmt"
	"net/http"

	"github.com/julienschmidt/httprouter"
	jsonlog "mooveit-backend.mooveit.com/internal/jsonlog"
)

func (app *application) routes() http.Handler {
	router := httprouter.New()

	// Convert httprouter.Handler to http.Handler
	router.HandlerFunc(http.MethodGet, "/healthcheck", app.healthcheckHandler)

	// Register the expvar handler for metrics
	router.Handler(http.MethodGet, "/debug/vars", expvar.Handler())

	// Create a middleware chain
	return app.recoverPanic(app.logRequest(router))
}

// recoverPanic middleware recovers from panics and logs the error
func (app *application) recoverPanic(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				w.Header().Set("Connection", "close")
				app.serverErrorResponse(w, r, fmt.Errorf("%s", err))
			}
		}()

		next.ServeHTTP(w, r)
	})
}

// logRequest middleware logs HTTP requests
func (app *application) logRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		jsonlog.InfoWithProperties("request received", map[string]string{
			"method": r.Method,
			"url":    r.URL.String(),
		})

		next.ServeHTTP(w, r)
	})
}
