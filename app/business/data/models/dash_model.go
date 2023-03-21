// Package models implements core user and kubernetes
// models, which are contained in the postgresql database.
package models

import (
	"context"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/joshua-seals/gopherhelx/app/business/sys/database"
)

// The DashID is the UserId.
type Dashboard struct {
	DashID      int    `db:"users_dash_id" json:"dash_id,omitempty"`
	UserSession string `db:"users_session" json:"user_session"`
	AppID       int    `db:"apps_app_id" json:"app_id"`
}

// The dashboard struct is good for one new entry
// But the UserDashboard will hold all app entries
// for any given user and their associated sessions.
type UserDashboard []Dashboard

func GetDashboard(ctx context.Context, db *sqlx.DB, userID int) (UserDashboard, error) {
	if err := database.StatusCheck(ctx, db); err != nil {
		return UserDashboard{}, fmt.Errorf("status check database: %w", err)
	}

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	const getDashboard = `SELECT * from dashboard where users_dash_id = $1`
	rows, err := db.QueryContext(ctx, getDashboard, userID)
	if err != nil {
		return UserDashboard{}, err
	}
	defer rows.Close()
	var userDash UserDashboard
	for rows.Next() {
		d := Dashboard{}
		err := rows.Scan(&d.DashID, &d.UserSession, &d.AppID)
		if err != nil {
			return UserDashboard{}, err
		}
		userDash = append(userDash, d)
	}
	return userDash, nil
}
