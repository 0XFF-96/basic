package handlers

import (
	"expvar"
	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
	"net/http"
	"net/http/pprof"
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
	//cgh := checkgrp.Handlers{
	//	Build: build,
	//	Log:   log,
	//	DB:    db,
	//}
	//mux.HandleFunc("/debug/readiness", cgh.Readiness)
	//mux.HandleFunc("/debug/liveness", cgh.Liveness)

	return mux
}