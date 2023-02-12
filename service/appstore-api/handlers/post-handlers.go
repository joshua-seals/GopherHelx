// All PUT method handlers are located herein which service the
// appstore-api
package handlers

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
)

// AppInstall installs a selected app into user dashboard.
// We move the app into user table in db as well.
func AppInstall(w http.ResponseWriter, r *http.Request) {
	// Dummy Data
	user := chi.URLParam(r, "userId")
	app := chi.URLParam(r, "appId")
	fmt.Fprintf(w, "<h1> Hello %s, installing app %s</h1>", user, app)
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

}
