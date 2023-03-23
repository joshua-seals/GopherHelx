// Package models implements core user and kubernetes
// models, which are contained in the postgresql database.
package models

import (
	"context"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/joshua-seals/gopherhelx/app/business/sys/database"
	_ "github.com/lib/pq"
)

// Dashboard is the model structure
// for our database transactions. DashId
// corresponds directly to user_id.
type Dashboard struct {
	DashID      int    `db:"users_dash_id" json:"dash_id,omitempty"`
	UserSession string `db:"users_session" json:"user_session"`
	AppID       int    `db:"apps_app_id" json:"app_id"`
}

// The dashboard struct is good for one new entry
// But the UserDashboard will hold all app entries
// for any given user and their associated sessions.
type UserDashboard []Dashboard

// GetDashboard uses the userId to return the corresponding user
// dashboard information, which includeds all installed apps
// and active session data/token.
func GetDashboard(ctx context.Context, db *sqlx.DB, userID string) (UserDashboard, error) {
	if err := database.StatusCheck(ctx, db); err != nil {
		return UserDashboard{}, fmt.Errorf("status check database: %w", err)
	}

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	const getDashboard = `SELECT * from dashboard where users_dash_id = $1`
	var userDash UserDashboard
	err := db.SelectContext(ctx, &userDash, getDashboard, userID)
	if err != nil {
		fmt.Println("Error with SelectContext")
		return UserDashboard{}, err
	}

	return userDash, nil
}

// GetDeploymentInfo manages the db queries and returns
// information for needed for creating a new kubernetes object, ie
// user starting the app, user session data, as well as the app
// information needed to populate k8s manifest image, port, etc.
func GetDeploymentInfo(ctx context.Context, db *sqlx.DB, userID string, appID string) (Dashboard, Application, error) {
	if err := database.StatusCheck(ctx, db); err != nil {
		return Dashboard{}, Application{}, fmt.Errorf("status check database: %w", err)
	}

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	var d Dashboard
	const getDashboard = `SELECT * from dashboard WHERE users_dash_id=$1 AND apps_app_id=$2;`
	err := db.GetContext(ctx, &d, getDashboard, userID, appID)
	if err != nil {
		return Dashboard{}, Application{}, err
	}

	var a Application
	const getApplication = `SELECT * from applications where app_id = $1;`
	err = db.GetContext(ctx, &a, getApplication, appID)
	if err != nil {
		return Dashboard{}, Application{}, err
	}

	return d, a, nil
}

// AddToDashboard adds an app to a users dashboard
// if no error returns then it is understood as success.
func (d *Dashboard) AddToDashboard(ctx context.Context, db *sqlx.DB, appID string) error {
	if err := database.StatusCheck(ctx, db); err != nil {
		return fmt.Errorf("status check database: %w", err)
	}

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	// Add logic to ensure this app exists first.

	const getDashboard = `
	INSERT into dashboard (users_dash_id, users_session, apps_app_id)
	VALUES ($1, $2, $3) ;`
	_, err := db.ExecContext(ctx, getDashboard, d.DashID, d.UserSession, appID)
	if err != nil {
		return err
	}

	return nil
}
