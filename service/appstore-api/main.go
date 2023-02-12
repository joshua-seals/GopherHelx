package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

type config struct {
	apiPort   int
	debugPort int
	env       string
}

// powLogger is simple logger created via two std library
// loggers customized with INFO and ERROR strings accordingly.
type powLogger struct {
	infoLog *log.Logger
	errLog  *log.Logger
}
type Application struct {
	config config
	logger *powLogger
}

func main() {
	// Create our logger, which will be passed through application.
	// Zap logger may provide better functionality in the future.
	infoLog := log.New(os.Stdout, "INFO: ", log.Ldate|log.Ltime)
	errLog := log.New(os.Stdout, "ERROR: ", log.Ldate|log.Ltime)
	powlog := &powLogger{
		infoLog: infoLog,
		errLog:  errLog,
	}

	run(powlog)

}

// Run handles the initialization of major components of the app
// Run will start both debug and appstore apis.
func run(log *powLogger) {

	cfg := config{
		apiPort:   3000,
		debugPort: 4000,
		env:       "dev",
	}

	// Initialize the app.
	app := &Application{
		config: cfg,
		logger: log,
	}

	//================================================================
	// Start debugServer

	// Create custom debugServer
	debugServer := &http.Server{
		Addr:     fmt.Sprintf(":%d", cfg.debugPort),
		Handler:  app.DebugRouter(),
		ErrorLog: app.logger.errLog,
	}

	go func() {
		if err := debugServer.ListenAndServe(); err != nil {
			app.logger.errLog.Fatal("shutdown", "debug router closed", err)
		}
	}()

	//==============================================================
	// Start apiServer

	// Create custom apiServer mux
	apiServer := &http.Server{
		Addr:     fmt.Sprintf(":%d", cfg.apiPort),
		Handler:  app.APIRouter(),
		ErrorLog: app.logger.errLog,
	}

	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	serverErrors := make(chan error, 1)
	go func() {
		app.logger.infoLog.Printf("starting server on %d", app.config.apiPort)
		serverErrors <- apiServer.ListenAndServe()
	}()

	// Listen on error and shutdown channels
	select {
	case err := <-serverErrors:
		app.logger.errLog.Println("Received error signal: ", err)
	case sig := <-shutdown:
		app.logger.infoLog.Println("Received signal: ", sig)
		app.logger.infoLog.Println("Performing shutdown operations.")
	}

}
