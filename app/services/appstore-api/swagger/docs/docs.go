// Package classification GopherHelx API.
//
// The purpose of this application is to provide a
// pluggable REST interface to the helxplatform that
// will act as orchestrator among varying microservices.
//
// Terms Of Service:
//
// there are no TOS at this moment, use at your own risk we take no responsibility
//
//	Schemes: http
//	Host: localhost:3000
//	BasePath: /
//	Version: 1.0.0
//	License: MIT http://opensource.org/licenses/MIT
//	Contact: Joshua Seals<jseals@renci.org>
//
//	Consumes:
//	- application/json
//	- application/yaml
//
//	Produces:
//	- application/json
//
// swagger:meta
package docs

import (
	"github.com/joshua-seals/gopherhelx/app/business/data/models"
)

// swagger:response appsResponse
type appsResponseWrapper struct {
	// List of Applications
	// in: body
	Body []models.Application
}

// swagger:parameters newApplication
type newApplicationWrapper struct {
	// in: query

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

// swagger:response newApplication
type newAppSuccessWrapper struct {
	//	in: body
	Body struct {
		// Success message with appID
		//
		// example: 117
		Success string
	}
}

// swagger:response userDashboard
type userDashboardWrapper struct {
	// User dashboard with associated apps installed
	// in: body
	Body []models.Dashboard
}

// The name of the deployment and service created.
// swagger:response startApplicaiton
type startApplicationWrapper struct {
	// in: body
	Body struct {
		// Deployment name
		//
		// required: true
		// example:  postgresql-pdry9f2
		Deployment string
		// Service name to connect to deployment
		//
		// required: true
		// example: postgresql-pdry9f2-service
		Services string
	}
}

// Successfully Added an App to the User Dashboard
// swagger:response addToDashSuccess
type addToDashSuccessWrapper struct {
	// in: body
	Body struct {
		// Success Response for added App
		//
		// required: true
		// example: App added to dashboard
		Success string
	}
}
