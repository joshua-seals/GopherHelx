// Package models implements core user and kubernetes
// models, which are contained in the postgresql database.
package models

// Application only needs a few things to describe it.
// We add the Deployment, Namespace, PodName and labels
// based on the specific user context parsed when a
// user requests to start an app.
type Application struct {
	AppID int
	Name  string
	Image string
	Port  string
}

type User struct {
	UserID  int
	Name    string
	Session string
}

type Dashboard struct {
	DashID int
	UID    int
	Apps   int
}
