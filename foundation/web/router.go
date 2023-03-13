package web

import (
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
)

// We can obfuscate our router here by embedding it's type.
type APIRouter struct {
	shutdown chan os.Signal
	log      zap.SugaredLogger
	db       *sqlx.DB
	ApiMux   *chi.Mux
}

func NewAPIRouter(shutdown chan os.Signal, log zap.SugaredLogger, db *sqlx.DB) *APIRouter {

	r := chi.NewRouter()

	return &APIRouter{
		shutdown: shutdown,
		log:      log,
		db:       db,
		ApiMux:   r,
	}

}
