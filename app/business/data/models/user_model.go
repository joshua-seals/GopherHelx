// Package models implements core user and kubernetes
// models, which are contained in the postgresql database.
package models

type User struct {
	UserID   int    `json:"user_id,omitempty"`
	UserName string `json:"user_name"`
	Session  string `json:"session,omitempty"`
}
