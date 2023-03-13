// This file contains all handlers associated with the appstore-api
// Additionally only the GET method routes are located here.
package v1

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/joshua-seals/gopherhelx/business/k8s"
)

// Dashboard shows the installed applications in the user
// specific dashboard. Dashboard is the entrypoint for
// users to start and stop applications.
func Dashboard(w http.ResponseWriter, r *http.Request) {

	req_user := chi.URLParam(r, "userId")
	if user, ok := UserDb[req_user]; ok {
		w.Header().Set("Content-Type", "application/json")
		err := json.NewEncoder(w).Encode(user)
		if err != nil {
			fmt.Println(err)
		}

	} else {
		fmt.Fprintf(w, "User not found: %d", http.StatusNotFound)
	}

}

// ViewApp is the user's connection point
// to a specific service & app
// governed by a session id.
func ViewApp(w http.ResponseWriter, r *http.Request) {

	k8s.ListDeployment()
}

// StartApp deploys an installed application from the user dashboard.
func StartApp(w http.ResponseWriter, r *http.Request) {
	// k8s.CreateDeployment()
	// appId := chi.URLParam(r, "appId")
	// k8s.CreateDeploymentFromFile(appId)
}

// StopApp will delete the currently deployed application
// from the user table and dashboard.
func StopApp(w http.ResponseWriter, r *http.Request) {

}

// RemoveApp will uninstall the applicaiton from the user dashboard.
// Subsequently, the app is removed from the user db table purview.
func RemoveApp(w http.ResponseWriter, r *http.Request) {
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
