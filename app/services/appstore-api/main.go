package main

import (
	"context"
	"errors"
	"expvar"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/ardanlabs/conf"
	"github.com/joshua-seals/gopherhelx/app/business/sys/database"
	"github.com/joshua-seals/gopherhelx/app/foundation/logger"
	"github.com/joshua-seals/gopherhelx/app/services/appstore-api/handlers"
	"go.uber.org/zap"
)

var build = "develop"

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
	// necessary structures like context, logger, db, etc
	// to the correspondent exepcting struct in the packages
	// imported into main.
	cfg := struct {
		conf.Version
		Web struct {
			APIHOST         string        `conf:"default:0.0.0.0:3000"`
			DebugHost       string        `conf:"default:0.0.0.0:4000"`
			Swagger         string        `conf:"default:0.0.0.0:1323"`
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
		Version: conf.Version{
			SVN:  build,
			Desc: "MIT License - Copyright (c) 2023 RENCI - Renaissance Computing Institute",
		},
	}

	const prefix = "APPSTORE"
	help, err := conf.ParseOSArgs(prefix, &cfg)
	if err != nil {
		if errors.Is(err, conf.ErrHelpWanted) {
			fmt.Println(help)
			return nil
		}
		return fmt.Errorf("parsing config: %w", err)
	}

	// =========================================================================
	// App Starting

	log.Infow("starting service", "version", build)
	defer log.Infow("shutdown complete")

	out, err := conf.String(&cfg)
	if err != nil {
		return fmt.Errorf("generating config for output: %w", err)
	}
	log.Infow("startup", "config", out)

	// =========================================================================
	// Database Support

	// Create connectivity to the database.
	log.Infow("startup", "status", "initializing database support", "host", cfg.DB.Host)

	db, err := database.Open(database.Config{
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

	expvar.NewString("build").Set(build)

	// Create custom debugServer
	debugMux := handlers.DebugMux(build, log, db)

	// Start the service listening for debug requests.
	// Not concerned with shutting this down with load shedding, ie don't need custom http.Server object
	go func() {
		if err := http.ListenAndServe(cfg.Web.DebugHost, debugMux); err != nil {
			log.Errorw("shutdown", "status", "debug router closed", "host", cfg.Web.DebugHost, "ERROR", err)
		}
	}()
	//================================================================
	// Start Swagger

	log.Infow("startup", "status", "debug router started", "host", cfg.Web.Swagger)

	// Create Swagger Router
	swagMux := handlers.SwaggerRoutes()

	// Start the service listening for debug requests.
	// Not concerned with shutting this down with load shedding, ie don't need custom http.Server object
	go func() {
		if err := http.ListenAndServe(cfg.Web.Swagger, swagMux); err != nil {
			log.Errorw("shutdown", "status", "debug router closed", "host", cfg.Web.Swagger, "ERROR", err)
		}
	}()

	//==============================================================
	// Start apiServer

	log.Infow("startup", "status", "initializing API support")

	// Make a channel to listen for an interrupt or terminate signal from the OS
	// Use buffered channel because the signal package requires it.
	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, syscall.SIGINT, syscall.SIGTERM)

	// We could pass shutdown through routes to encourage
	// graceful shutdown.
	apiMux := handlers.APIRoutes(log, db)
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

}
