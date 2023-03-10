package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/joshua-seals/gopherhelx/app/services/appstore-api/handlers"
	databse "github.com/joshua-seals/gopherhelx/foundation/database"
	"github.com/joshua-seals/gopherhelx/foundation/logger"
	"go.uber.org/zap"
)

type Config struct {
	apiPort   int
	debugPort int
	env       string
}

func main() {

	// Construct the application logger.
	log, err := logger.New("APPSTORE-API")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer log.Sync()

	if err := run(log); err != nil {
		log.Errorw("startup", "ERROR", err)
		os.Exit(1)
	}

}

// Run handles the initialization of major components of the app
// Run will start both debug and appstore apis.
func run(log *zap.SugaredLogger) error {

	// We only pull in packages to main, careful to avoid global level access.
	// To do so a common pattern here is to mirror the package level
	// structures in this annonomous struct, then to pass down
	// necessary api libraries like context, logger, db, etc
	// to the correspondent exepcting struct in the packages
	// imported into main.
	cfg := struct {
		Config
		Web struct {
			APIHOST         string        `conf:"default:0.0.0.0:3000"`
			DebugHost       string        `conf:"default:0.0.0.0:4000"`
			ReadTimeout     time.Duration `conf:"default:5s"`
			WriteTimeout    time.Duration `conf:"default:10s"`
			IdleTimeout     time.Duration `conf:"default:120s"`
			ShutdownTimeout time.Duration `conf:"default:20s,mask"` //mask can be used as well as noprint here
		}
		DB struct {
			User         string `conf:"default:postgres"`
			Password     string `conf:"default:postgres,mask"`
			Host         string `conf:"default:localhost"`
			Name         string `conf:"default:postgres"`
			MaxIdleConns int    `conf:"default:0"`
			MaxOpenConns int    `conf:"default:0"`
			DisableTLS   bool   `conf:"default:true"`
		}
	}{
		Config: Config{
			apiPort:   3000,
			debugPort: 4000,
			env:       "dev",
		},
	}

	// =========================================================================
	// Database Support

	// Create connectivity to the database.
	log.Infow("startup", "status", "initializing database support", "host", cfg.DB.Host)

	db, err := databse.Open(databse.Config{
		User:         cfg.DB.User,
		Password:     cfg.DB.Password,
		Host:         cfg.DB.Host,
		Name:         cfg.DB.Name,
		MaxIdleConns: cfg.DB.MaxIdleConns,
		MaxOpenConns: cfg.DB.MaxOpenConns,
		DisableTLS:   cfg.DB.DisableTLS,
	})
	if err != nil {
		return fmt.Errorf("connecting to db: %w", err)
	}
	defer func() {
		log.Infow("shutdown", "status", "stopping database support", "host", cfg.DB.Host)
		db.Close()
	}()

	//================================================================
	// Start debugServer

	log.Infow("startup", "status", "debug router started", "host", cfg.Web.DebugHost)

	// Create custom debugServer
	debugMux := handlers.DebugMux(log, db)

	// Start the service listening for debug requests.
	// Not concerned with shutting this down with load shedding, ie don't need custom http.Server object
	go func() {
		if err := http.ListenAndServe(cfg.Web.DebugHost, debugMux); err != nil {
			log.Errorw("shutdown", "status", "debug router closed", "host", cfg.Web.DebugHost, "ERROR", err)
		}
	}()

	//==============================================================
	// Start apiServer

	log.Infow("startup", "status", "initializing API support")

	// Make a channel to listen for an interrupt or terminate signal from the OS
	// Use buffered channel because the signal package requires it.
	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, syscall.SIGINT, syscall.SIGTERM)

	apiMux := handlers.APIMux(handlers.APIMuxConfig{
		Shutdown: shutdown,
		Log:      log,
		Auth:     auth,
		DB:       db,
	})
	// Construct a server to service requests against the mux
	api := http.Server{
		Addr:         cfg.Web.APIHOST,
		Handler:      apiMux,
		ReadTimeout:  cfg.Web.ReadTimeout,
		WriteTimeout: cfg.Web.WriteTimeout,
		IdleTimeout:  cfg.Web.IdleTimeout,
		ErrorLog:     zap.NewStdLog(log.Desugar()),
	}

	// Make a channel to listen for errors coming from listener. Use a
	// buffered channel so the goroutine can exit if we don't get an error.
	serverErrors := make(chan error, 1)

	// Start the service listening for api requests
	go func() {
		log.Infow("startup", "status", "api router started", "host", api.Addr)
		// blocking on ListenAndserver() here.
		serverErrors <- api.ListenAndServe()
	}()

	// =============================================================
	// Shutdown

	// Blocking main and waiting for shutdown.
	select {
	case err := <-serverErrors:
		return fmt.Errorf("server error %w", err)
	case sig := <-shutdown:
		log.Infow("shutdown", "status", "shutdown started", "signal", sig)
		defer log.Infow("shutdown", "status", "shutdown complete", "signal", sig)

		// Give outstanding requests a deadline for completion
		ctx, cancel := context.WithTimeout(context.Background(), cfg.Web.ShutdownTimeout)
		defer cancel()

		// Asking listener to shutdown and shed load
		if err := api.Shutdown(ctx); err != nil {
			api.Close()
			return fmt.Errorf("could not stop server gracefully: %w", err)
		}
	}

	return nil

	// // Create custom apiServer mux
	// apiServer := handlers.APIRouter(log, db)

	// shutdown := make(chan os.Signal, 1)
	// signal.Notify(shutdown, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	// serverErrors := make(chan error, 1)
	// go func() {
	// 	app.logger.infoLog.Printf("starting server on %d", app.config.apiPort)
	// 	serverErrors <- apiServer.ListenAndServe()
	// }()

	// // Listen on error and shutdown channels
	// select {
	// case err := <-serverErrors:
	// 	app.logger.errLog.Println("Received error signal: ", err)
	// case sig := <-shutdown:
	// 	app.logger.infoLog.Println("Received signal: ", sig)
	// 	app.logger.infoLog.Println("Performing shutdown operations.")
	// }

}
