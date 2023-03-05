package main

import (
	"expvar"
	"net/http"
	"net/http/pprof"

	"github.com/go-chi/chi/v5"
	"github.com/joshua-seals/gopherhelx/app/services/appstore-api/handlers"
)

func (app *Application) APIRouter() *chi.Mux {
	// This router is used in the srv (http.Server) created
	// as the Handler and is where all api routes are located.
	// The corresponding functions however will be located in
	// handlers package.
	router := chi.NewRouter()

	// router.Get("/login", handlers.Login)
	// router.Post("/Userlogin", handlers.UserLogin)

	router.Get("/app/list", handlers.AppList)
	router.Get("/{userId}/dashboard", handlers.Dashboard)
	router.Get("/{userId}/dashboard/session", handlers.Session)
	// /{userId}/dashboard/session/{appId}/{sessionId}"

	router.Post("/{userId}/app/install/{appId}", handlers.AppInstall)
	router.Post("/{userId}/dashboard/start/{appId}", handlers.StartApp)

	router.Delete("/{userId}/dashboard/stop/{appId}", handlers.StopApp)
	router.Delete("/{userId}/dashboard/remove/{appId}", handlers.RemoveApp)

	return router
}

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

func (app *Application) DebugRouter() http.Handler {
	// Imbed a copy of the above function.
	mux := DebugStandardLibraryMux()

	// Here we reference a struct from handlers.debug
	// And Instantiate with app Applicaiton type variables.
	dbug := handlers.DebugHandler{
		InfoLog: app.logger.infoLog,
		ErrLog:  app.logger.errLog,
	}

	mux.HandleFunc("/debug/readiness", dbug.Readiness)
	mux.HandleFunc("/debug/liveness", dbug.Liveness)

	return mux
}
