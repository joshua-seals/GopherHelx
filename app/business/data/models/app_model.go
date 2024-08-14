// Package models implements core user and kubernetes
// models, which are contained in the postgresql database.
package models

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/joshua-seals/gopherhelx/app/business/sys/database"
	_ "github.com/lib/pq"
)

// Application only needs a few things to describe it.
// We add the Deployment, Namespace, PodName and labels
// based on the specific user context parsed when a
// user requests to start an app.
//

// Application defines the structure for an Application.
// swagger:model
type Application struct {
	// the id for the product
	//
	// required: false
	// min: 1
	AppID int `json:"app_id,omitempty" db:"app_id"`

	// the name of the application
	//
	// required: true
	// example: scipy
	AppName string `json:"app_name" db:"app_name"`

	// the image of the application
	//
	// required: true
	// example: helxplatform/scipy
	Image string `json:"image" db:"image"`

	// the port for the application
	//
	// required: true
	// example: 8888
	Port int `json:"port" db:"port"`
}

type Apps []Application

// AddNewApplication installs new applications into the database
// returning the id of app_id upon success.
func (a *Application) AddNewApplication(ctx context.Context, db *sqlx.DB) (string, error) {

	if err := database.StatusCheck(ctx, db); err != nil {
		return "", fmt.Errorf("status check database: %w", err)
	}

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	tx, err := db.Begin()
	if err != nil {
		return "", err
	}

	var id int
	const insertApp = `
	INSERT INTO applications 
	(app_name, image, port) VALUES ($1, $2, $3) 
	RETURNING app_id;`
	if err := tx.QueryRowContext(ctx,
		insertApp,
		a.AppName,
		a.Image,
		a.Port).Scan(&id); err != nil {
		return "", err
	}
	err = tx.Commit()
	if err != nil {
		return "", err
	}
	// Convert id to string
	return strconv.Itoa(id), nil

}

// AppList returns a list of applications available
// to be installed in the user dashboard.
func AppList(ctx context.Context, db *sqlx.DB) (Apps, error) {

	if err := database.StatusCheck(ctx, db); err != nil {
		return Apps{}, fmt.Errorf("status check database: %w", err)
	}

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	const getApps = `SELECT * from applications;`
	var appList Apps
	err := db.SelectContext(ctx, &appList, getApps)
	if err != nil {
		return Apps{}, err
	}

	return appList, nil
}
