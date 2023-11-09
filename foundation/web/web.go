// Package web contains a small web framework extension.
package web

import (
	"context"
	"github.com/dimfeld/httptreemux/v5"
	"github.com/google/uuid"
	"net/http"
	"os"
	"syscall"
	"time"
)

// A Handler is a type that handles a http request within our own little mini
// framework.
type Handler func(ctx context.Context, w http.ResponseWriter, r *http.Request) error

// App is the entrypoint into our application and what configures our context
// object for each of our http handlers. Feel free to add any configuration
// data/logic on this App struct.
type App struct {
	*httptreemux.ContextMux
	// otmux    http.Handler
	shutdown chan os.Signal

	mw []Middleware
}

// NewApp creates an App value that handle a set of routes for the application.
func NewApp(shutdown chan os.Signal, mw ...Middleware) *App {

	// Create an OpenTelemetry HTTP Handler which wraps our router. This will start
	// the initial span and annotate it with information about the request/response.
	//
	// This is configured to use the W3C TraceContext standard to set the remote
	// parent if a client request includes the appropriate headers.
	// https://w3c.github.io/trace-context/

	mux := httptreemux.NewContextMux()

	return &App{
		ContextMux: mux,
		shutdown:   shutdown,
		mw:         mw,
	}
}

// SignalShutdown is used to gracefully shut down the app when an integrity
// issue is identified.
func (a *App) SignalShutdown() {
	a.shutdown <- syscall.SIGTERM
}

// Handle sets a handler function for a given HTTP method and path pair
// to the application server mux.
func (a *App) Handle(method string, group string, path string, handler Handler, mw ...Middleware) {

	// middleware 的顺序非常重要‼️
	// 1. First wrap handler specific middleware around this handler.
	handler = wrapMiddleware(mw, handler)

	// 2. Add the application's general middleware to the handler chain.
	handler = wrapMiddleware(a.mw, handler)

	h := func(w http.ResponseWriter, r *http.Request) {
		// PRE CODE PROCESSING
		// Pull the context from the request and
		// use it as a separate parameter.
		ctx := r.Context()

		// Capture the parent request span from the context.
		// span := trace.SpanFromContext(ctx)

		// Set the context with the required values to
		// process the request.

		uuid, _ := uuid.NewUUID()
		v := Values{
			// TraceID: span.SpanContext().TraceID().String(),
			TraceID: uuid.String(),
			// Tracer:  a.tracer,
			Now: time.Now().UTC(),
		}
		ctx = context.WithValue(ctx, key, &v)

		// INJECT CODE, can only exist in the business layer
		// Logging Started
		if err := handler(ctx, w, r); err != nil {
			// Logging error - handle it

			// ERROR HANDLING
			return
		}

		// Logging Ended
		// POST CODE PROCESSING
	}

	finalPath := path
	if group != "" {
		finalPath = "/" + group + path
	}
	a.ContextMux.Handle(method, finalPath, h)
}
