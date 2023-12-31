package handlers

import (
	"expvar"
	"github.com/jmoiron/sqlx"
	"github.com/yourusername/basic-a/app/services/sales-api/handlers/debug/checkgrp"
	"github.com/yourusername/basic-a/app/services/sales-api/handlers/v1/testgrp"
	"github.com/yourusername/basic-a/app/services/sales-api/handlers/v1/usergrp"
	"github.com/yourusername/basic-a/business/core/user"
	userStore "github.com/yourusername/basic-a/business/data/store/user"
	"github.com/yourusername/basic-a/business/data/store/user/usercache"
	"go.opentelemetry.io/otel/trace"

	// 	"github.com/yourusername/basic-a/business/sys/auth"
	"github.com/yourusername/basic-a/business/web/auth"
	"github.com/yourusername/basic-a/business/web/mid"
	"github.com/yourusername/basic-a/foundation/web"

	"go.uber.org/zap"
	"net/http"
	"net/http/pprof"
	"os"
)

// DebugStandardLibraryMux registers all the debug routes from the standard library
// into a new mux bypassing the use of the DefaultServerMux. Using the
// DefaultServerMux would be a security risk since a dependency could inject a
// handler into our service without us knowing it.
func DebugStandardLibraryMux() *http.ServeMux {
	mux := http.NewServeMux()

	// Register all the standard library debug endpoints.
	mux.HandleFunc("/debug/pprof/", pprof.Index)
	mux.HandleFunc("/debug/pprof/cmdline", pprof.Cmdline)
	mux.HandleFunc("/debug/pprof/profile", pprof.Profile)
	mux.HandleFunc("/debug/pprof/symbol", pprof.Symbol)
	mux.HandleFunc("/debug/pprof/trace", pprof.Trace)
	mux.Handle("/debug/vars", expvar.Handler())

	return mux
}

// vars 包的相关作用～
// 	mux.Handle("/debug/vars", expvar.Handler())

// Do calls f for each exported variable.
// The global variable map is locked during the iteration,
// but existing entries may be concurrently updated.
//func Do(f func(KeyValue)) {
//	varKeysMu.RLock()
//	defer varKeysMu.RUnlock()
//	for _, k := range varKeys {
//		val, _ := vars.Load(k)
//		f(KeyValue{k, val.(Var)})
//	}
//}
//
//func expvarHandler(w http.ResponseWriter, r *http.Request) {
//	w.Header().Set("Content-Type", "application/json; charset=utf-8")
//	fmt.Fprintf(w, "{\n")
//	first := true
//	Do(func(kv KeyValue) {
//		if !first {
//			fmt.Fprintf(w, ",\n")
//		}
//		first = false
//		fmt.Fprintf(w, "%q: %s", kv.Key, kv.Value)
//	})
//	fmt.Fprintf(w, "\n}\n")
//}

// DebugMux registers all the debug standard library routes and then custom
// debug application routes for the service. This bypassing the use of the
// DefaultServerMux. Using the DefaultServerMux would be a security risk since
// a dependency could inject a handler into our service without us knowing it.
func DebugMux(build string, log *zap.SugaredLogger, db *sqlx.DB) http.Handler {
	mux := DebugStandardLibraryMux()

	// Register debug check endpoints.
	cgh := checkgrp.Handlers{
		Build: build,
		Log:   log,
		DB:    db,
	}
	mux.HandleFunc("/debug/readiness", cgh.Readiness)
	mux.HandleFunc("/debug/liveness", cgh.Liveness)

	return mux
}

// APIMuxConfig contains all the mandatory systems required by handlers.
type APIMuxConfig struct {
	Shutdown chan os.Signal
	Log      *zap.SugaredLogger
	Auth     *auth.Auth
	DB       *sqlx.DB
	Tracer   trace.Tracer
}

// Options represent optional parameters.
type Options struct {
	corsOrigin string
}

// APIMux constructs a http.Handler with all application routes defined.
func APIMux(cfg APIMuxConfig, options ...func(opts *Options)) *web.App {
	app := web.NewApp(
		cfg.Shutdown,
		mid.Logger(cfg.Log),
		mid.Errors(cfg.Log),
		mid.Panics(),
		mid.Metrics(),
	)

	// Load the routes for the difference versions API
	v1(app, cfg)

	return app
}

func v1(app *web.App, cfg APIMuxConfig) {
	const version = "v1"
	authen := mid.Authenticate(cfg.Auth)
	admin := mid.Authorize(auth.RuleAdminOnly)

	// Register debug check endpoints.
	tgh := testgrp.Handlers{
		Log: cfg.Log,
		// DB:    db,
	}

	ugh := usergrp.Handlers{
		User: user.NewCore(usercache.NewStore(cfg.Log, userStore.NewStore(cfg.Log, cfg.DB))),
		Auth: cfg.Auth,
	}
	app.Handle(http.MethodGet, version, "/users/token/:kid", ugh.Token)
	app.Handle(http.MethodGet, version, "/users/:page/:rows", ugh.Query, authen, admin)
	app.Handle(http.MethodGet, version, "/users/:id", ugh.QueryByID, authen)
	app.Handle(http.MethodPost, version, "/users", ugh.Create, authen, admin)
	app.Handle(http.MethodPut, version, "/users/:id", ugh.Update, authen, admin)
	app.Handle(http.MethodDelete, version, "/users/:id", ugh.Delete, authen, admin)

	// handle path
	app.Handle(http.MethodGet, version, "/test", tgh.Test)
	app.Handle(http.MethodGet, version, "/testauth", tgh.Test, authen)
	app.Handle(http.MethodGet, version, "/testadmin", tgh.Test, authen, mid.Authorize("ADMIN"))
}
