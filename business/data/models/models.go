// Package models implements core user and kubernetes
// models, which are contained in the postgresql database.
package models

type Application struct {
	AppID          int
	Namespace      string
	DeploymentName string
	PodName        string
	Image          string
	Port           string
}

type Dashboard struct {
	Apps []Application
}

type User struct {
	ID        int
	Name      string
	Session   string
	Dashboard Dashboard
}
