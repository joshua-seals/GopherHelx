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
type Application struct {
	AppID   int    `db:"app_id" json:"app_id,omitempty"`
	AppName string `db:"app_name" json:"app_name"`
	Image   string `db:"image" json:"image"`
	Port    int    `db:"port" json:"port"`
}

// Apps is a slice of Applicaiton.
type Apps []Application

// Returns the newly installed appId and error.
func (a *Application) AddNewApplication(ctx context.Context, db *sqlx.DB) (string, error) {
	// https://go.dev/doc/database/execute-transactions
	if err := database.StatusCheck(ctx, db); err != nil {
		return "", fmt.Errorf("status check database: %w", err)
	}

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	tx, err := db.Begin()
	if err != nil {
		return "", err
	}
	// Check to make sure app doesn't already exist.
	// If so, return the response error.
	// var checkApp = Application{}
	// checkQuery := "SELECT * from applications where app_name=$1"
	// err = tx.QueryRow(checkQuery, a.AppName).Scan(&checkApp)

	// if ok := errors.Is(err, ErrNoRows); !ok {
	// 	return "", err
	// }
	// if checkApp.AppName == a.AppName {
	// 	response := fmt.Errorf("Application by this name is instaled already %s", a.AppName)
	// 	return "", response
	// }
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

func AppList(ctx context.Context, db *sqlx.DB) (Apps, error) {

	if err := database.StatusCheck(ctx, db); err != nil {
		return Apps{}, fmt.Errorf("status check database: %w", err)
	}

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	// tx, err := db.Begin()
	// if err != nil {
	// 	return Apps{}, err
	// }
	const getApps = `SELECT * from applications;`
	rows, err := db.QueryContext(ctx, getApps)
	if err != nil {
		return Apps{}, err
	}
	// No errors we commit the Select statement.
	// err = tx.Commit()
	// if err != nil {
	// 	return Apps{}, err
	// }
	defer rows.Close()
	var appList Apps
	for rows.Next() {
		a := Application{}
		err := rows.Scan(&a.AppID, &a.AppName, &a.Image, &a.Port)
		if err != nil {
			return Apps{}, err
		}
		appList = append(appList, a)
	}
	return appList, nil
}
