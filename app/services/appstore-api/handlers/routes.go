// package handlers creates all the routes for the program
// and returns the routers/muxes that manage them.
package handlers

import (
	"expvar"
	"net/http"
	"net/http/pprof"

	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
	apiMiddleware "github.com/go-openapi/runtime/middleware"
	"github.com/jmoiron/sqlx"
	v1 "github.com/joshua-seals/gopherhelx/app/services/appstore-api/handlers/v1"

	"go.uber.org/zap"
)

/*
TODO:
	- Chi needs to use zap logger instead of chi middleware.logger.
		" router.Use(middleware.Logger) "
	- Add User Authorization
	- Add 'stateful' authentication tokens.
*/

// APIRoutes holds all api routes currently served.
func APIRoutes(log *zap.SugaredLogger, db *sqlx.DB) *chi.Mux {
	// This router is used in the srv (http.Server) created
	// as the Handler and is where all api routes are located.
	// The corresponding functions however will be located in
	// v1 package. This allows for simultaneous work and testing of
	// new routes from a v2 api, without the need to change every
	// endpoint.
	router := chi.NewRouter()

	// Core binds the v1 handlers
	core := v1.CoreHandler{
		Log: log,
		DB:  db,
	}

	router.Use(middleware.Logger)
	router.NotFound(core.NotFoundResponse)
	router.MethodNotAllowed(core.MethodNotAllowedResponse)

	router.Get("/app/list", core.AppList)
	router.Post("/app/new", core.NewApplication)

	router.Get("/dashboard/{userId}", core.Dashboard)
	router.Post("/dashboard/{userId}/install/{appId}", core.AddToDashboard)
	router.Put("/dashboard/{userId}/start/{appId}", core.StartApp)
	router.Get("/dashboard/{userId}/{sessionId}/connect/{appId}", core.ViewApp)
	router.Delete("/dashboard/{userId}/stop/{appId}", core.StopApp)
	router.Delete("/dashboard/{userId}/remove/{appId}", core.RemoveApp)

	return router
}

// DebugStandardLibraryMux registers all the debug routes from the standard library
// into a new mux bypassing the use of the DefaultServerMux. Using the
// DefaultServerMux would be a security risk since a dependency could inject a
// handler into our service without us knowing it.
func DebugStandardLibraryMux() *http.ServeMux {
	mux := http.NewServeMux()

	// Register all the standard library debug endpoints.
	// These debug endpoints can be viewed at localhost:4000/debug/pprof
	mux.HandleFunc("/debug/pprof/", pprof.Index)
	mux.HandleFunc("/debug/pprof/cmdline", pprof.Cmdline)
	mux.HandleFunc("/debug/pprof/profile", pprof.Profile)
	mux.HandleFunc("/debug/pprof/symbol", pprof.Symbol)
	mux.HandleFunc("/debug/pprof/trace", pprof.Trace)
	mux.Handle("/debug/vars", expvar.Handler())

	return mux
}

func DebugMux(build string, log *zap.SugaredLogger, db *sqlx.DB) http.Handler {
	// Imbed a copy of the above function.
	mux := DebugStandardLibraryMux()

	// Here we reference a struct from handlers.debug
	// And Instantiate with app Applicaiton type variables.
	dbug := v1.DebugHandler{
		Build: build,
		Log:   log,
		DB:    db,
	}

	mux.HandleFunc("/debug/readiness", dbug.Readiness)
	mux.HandleFunc("/debug/liveness", dbug.Liveness)

	return mux
}

func SwaggerRoutes() *chi.Mux {
	swag_router := chi.NewRouter()

	ops := apiMiddleware.RedocOpts{SpecURL: "/swagger.yaml"}
	sh := apiMiddleware.Redoc(ops, nil)

	swag_router.Handle("/docs", sh)
	swag_router.Handle("/swagger.yaml", http.FileServer(http.Dir("./appstore-api/swagger/")))
	return swag_router
}
