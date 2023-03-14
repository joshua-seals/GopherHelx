// Package models implements core user and kubernetes
// models, which are contained in the postgresql database.
package models

import (
	"context"

	"github.com/jmoiron/sqlx"
)

// Application only needs a few things to describe it.
// We add the Deployment, Namespace, PodName and labels
// based on the specific user context parsed when a
// user requests to start an app.
type Application struct {
	AppID   int    `json:"app_id,omitempty"`
	AppName string `json:"app_name"`
	Image   string `json:"image"`
	Port    string `json:"port"`
}

// Unexported and protected access to db
// via the models method newApplication.
// Returns the newly installed appId or error
func (a *Application) AddNewApplication(ctx context.Context, db *sqlx.DB) (string, error) {

	return "123", nil
}
