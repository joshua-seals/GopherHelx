// This file contains all handlers associated with the appstore-api
// Additionally only the GET method routes are located here.
package v1

import (
	"context"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/joshua-seals/gopherhelx/app/business/data/models"
	"github.com/joshua-seals/gopherhelx/app/business/k8s"
)

// Dashboard shows the installed applications in the user
// specific dashboard. Dashboard is the entrypoint for
// users to start and stop applications.
func (c CoreHandler) Dashboard(w http.ResponseWriter, r *http.Request) {

	userId, err := strconv.Atoi(chi.URLParam(r, "userId"))
	if err != nil {
		c.serverErrorResponse(w, r, err)
		return
	}
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

// ViewApp is the user's connection point
// to a specific service & app
// governed by a session id.
func (c CoreHandler) ViewApp(w http.ResponseWriter, r *http.Request) {

	k8s.ListDeployment()
}

// StartApp deploys an application from the user dashboard to kubernetes env.
func (c CoreHandler) StartApp(w http.ResponseWriter, r *http.Request) {
	// appId := chi.URLParam(r, "appId")
	// k8s.CreateDeploymentFromFile(appId)
}

// StopApp will delete the currently deployed application
// from the user table and dashboard.
func (c CoreHandler) StopApp(w http.ResponseWriter, r *http.Request) {

}

// RemoveApp will uninstall the applicaiton from the user dashboard.
// Subsequently, the app is removed from the user db table purview.
func (c CoreHandler) RemoveApp(w http.ResponseWriter, r *http.Request) {
	app := chi.URLParam(r, "appId")
	k8s.DeleteDeployment(app)
}

// =========== Dummy Data ===============
type User struct {
	Name      string   `json:"name"`
	Dashboard []string `json:"dashboard"` // Installed user apps view
}

type UserList map[string]User

// Simulate a database for testing purposes
var UserDb = UserList{
	"1": User{
		Name:      "Bob",
		Dashboard: []string{"Webtop", "R-Studio", "Atlas"},
	},
	"2": User{
		Name:      "Sally",
		Dashboard: []string{"Jupyter", "Webtop", "Filebrowser"},
	},
}
