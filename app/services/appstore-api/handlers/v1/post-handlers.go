// All PUT method handlers are located herein which service the
// appstore-api
package v1

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

var apps = map[int]string{1: "Webtop", 2: "Filebrowser", 3: "Jupyter", 4: "Balsam", 5: "PGAdmin"}

// AppInstall installs a selected app into user dashboard.
// We move the app into user table in db as well.
func AppInstall(w http.ResponseWriter, r *http.Request) {
	// Dummy Data
	user := chi.URLParam(r, "userId")
	appId, err := strconv.Atoi(chi.URLParam(r, "appId"))
	if err != nil {
		fmt.Println("Error: Cannot convert string to int")
	}
	if app, ok := apps[appId]; ok {
		fmt.Fprintf(w, "<h1> Hello %s, installing app %s</h1>", user, app)
	} else {
		fmt.Fprintf(w, "<h1> Requested app was not found %s</h1>", app)
	}

}

// StartApp deploys an installed application from the user dashboard.
// A sessionId with "app-name"-"session" is generated first and
// added to the current user db table.
// The sessionId is then passed to Scheduling Package.
// The "Scheduling Package" is used
// to leverage kubernetes client functions
// specifically to provision a deployment and service resource
// corresponding to the desired application.
func StartApp(w http.ResponseWriter, r *http.Request) {
	// k8s.CreateDeployment()
	// appId := chi.URLParam(r, "appId")
	// k8s.CreateDeploymentFromFile(appId)
}
