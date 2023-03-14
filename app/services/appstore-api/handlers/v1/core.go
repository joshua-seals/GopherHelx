package v1

import (
	"encoding/json"
	"net/http"

	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
)

// CoreHandler can hold middleware, context, and logger, etc.
// This struct will bind the API associated methods
// found in apps.go, dashboard.go and any helpers needed for
// standardization of data flow.
type CoreHandler struct {
	Log *zap.SugaredLogger
	DB  *sqlx.DB
}

type envelope map[string]any

// Method for maintaining a standard response to the client.
func (c *CoreHandler) writeJSON(w http.ResponseWriter, status int, data envelope, headers http.Header) error {
	js, err := json.MarshalIndent(data, "", "\t")
	if err != nil {
		return err
	}
	js = append(js, '\n')
	for key, value := range headers {
		w.Header()[key] = value
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write(js)
	return nil
}
