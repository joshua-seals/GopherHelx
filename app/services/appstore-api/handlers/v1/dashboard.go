// This file contains all handlers associated with the appstore-api
// Additionally only the GET method routes are located here.
package v1

import (
	"context"
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/joshua-seals/gopherhelx/app/business/data/models"
	"github.com/joshua-seals/gopherhelx/app/business/k8s"
)

/*
*	TODO: Dynamically pull namespace info. Currently hard coded.
 */

// Dashboard shows the installed applications in the user
// specific dashboard. Dashboard is the entrypoint for
// users to start and stop applications.
func (c CoreHandler) Dashboard(w http.ResponseWriter, r *http.Request) {

	// userId, err := strconv.Atoi(chi.URLParam(r, "userId"))
	userId := chi.URLParam(r, "userId")
	// if err != nil {
	// 	c.serverErrorResponse(w, r, err)
	// 	return
	// }
	ctx := context.Background()
	userDash, err := models.GetDashboard(ctx, c.DB, userId)

	if err != nil {
		c.serverErrorResponse(w, r, err)
		return
	}

	status := http.StatusAccepted
	data := envelope{"dashboard": userDash}
	if err := c.writeJSON(w, status, data, nil); err != nil {
		c.serverErrorResponse(w, r, err)
		return
	}
}

// StartApp deploys an application from the user dashboard to kubernetes env.
func (c CoreHandler) StartApp(w http.ResponseWriter, r *http.Request) {
	userId := chi.URLParam(r, "userId")

	appId := chi.URLParam(r, "appId")

	ctx := context.Background()
	d, a, err := models.GetDeploymentInfo(ctx, c.DB, userId, appId)
	if err != nil {
		c.serverErrorResponse(w, r, err)
		return
	}
	// TODO: Need to dynamically pull namespace
	// Ensure data follows kubernetes requirements
	// ie. lowercase naming and remove ToLower mess.
	depName := strings.ToLower(a.AppName + "-" + d.UserSession)
	deployment := k8s.Deployment{
		DName:      depName,
		DNamespace: "appstore-system",
		DLabels: map[string]string{
			"apps": strings.ToLower(a.AppName),
		},
		AName:  strings.ToLower(a.AppName),
		AImage: a.Image,
		APort:  a.Port,
	}
	err = deployment.CreateDeployment()
	if err != nil {
		c.serverErrorResponse(w, r, err)
	}
}

// StopApp will delete the currently deployed application
// from the user table and dashboard.
func (c CoreHandler) StopApp(w http.ResponseWriter, r *http.Request) {

}

// ViewApp is the user's connection point
// to a specific service & app
// governed by a session id.
func (c CoreHandler) ViewApp(w http.ResponseWriter, r *http.Request) {

	k8s.ListDeployment()
}

func (c CoreHandler) AddToDashboard(w http.ResponseWriter, r *http.Request) {
	// userId := chi.URLParam(r, "userId")
	// appId := chi.URLParam(r, "appId")
}

// RemoveApp will uninstall the applicaiton from the user dashboard.
// Subsequently, the app is removed from the user db table purview.
func (c CoreHandler) RemoveApp(w http.ResponseWriter, r *http.Request) {
	app := chi.URLParam(r, "appId")
	k8s.DeleteDeployment(app)
}
