package handlers

import (
	"expvar"
	"net/http"
	"net/http/pprof"
	"os"

	"github.com/jmoiron/sqlx"
	v1 "github.com/joshua-seals/gopherhelx/app/services/appstore-api/handlers/v1"
	"github.com/joshua-seals/gopherhelx/foundation/web"
	"go.uber.org/zap"
)

// APIMuxConfig contains all the mandatory systems required by handlers.
type APICoreConfig struct {
	Shutdown chan os.Signal
	Log      *zap.SugaredLogger
	DB       *sqlx.DB
}

func APIRoutes(cfg APICoreConfig) *web.APIRouter {
	// This router is used in the srv (http.Server) created
	// as the Handler and is where all api routes are located.
	// The corresponding functions however will be located in
	// handlers package.
	router := web.NewAPIRouter(cfg.Shutdown, *cfg.Log, cfg.DB)

	router.ApiMux.Get("/app/list", v1.AppList)
	router.ApiMux.Post("/app/new", v1.NewApplication)
	router.ApiMux.Post("/app/install/{appId}/{userId}", v1.AddToDashboard)

	router.ApiMux.Get("/dashboard/{userId}", v1.Dashboard)
	router.ApiMux.Post("/dashboard/{userId}/start/{appId}", v1.StartApp)
	router.ApiMux.Get("/dashboard/{userId}/session/{appId}/{sessionId}", v1.ViewApp)
	router.ApiMux.Delete("/dashboard/{userId}/stop/{appId}", v1.StopApp)
	router.ApiMux.Delete("/dashboard/{userId}/remove/{appId}", v1.RemoveApp)

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
