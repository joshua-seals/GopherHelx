// Package models implements core user and kubernetes
// models, which are contained in the postgresql database.
package models

type Dashboard struct {
	DashID      int    `json:"dash_id,omitempty"`
	UserID      int    `json:"user_id"`
	AppID       int    `json:"app_id"`
	UserSession string `json:"user_session"`
}
