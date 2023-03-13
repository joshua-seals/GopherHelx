// All PUT method handlers are located herein which service the
// appstore-api
package v1

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

var apps = map[int]string{1: "Webtop", 2: "Filebrowser", 3: "Jupyter", 4: "Balsam", 5: "PGAdmin"}

// AddToDashboard installs a selected app from app/list into user dashboard.
// This action also triggers an update to the user db table
// where the app is then added to specific user dashboard.
func AddToDashboard(w http.ResponseWriter, r *http.Request) {
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

// Applist shows the list of applications available
// for a user to install in their dashboard.
func AppList(w http.ResponseWriter, r *http.Request) {
	apps := apps
	w.Header().Set("Content-Type", "application/json")
	err := json.NewEncoder(w).Encode(apps)
	if err != nil {
		fmt.Println(err)
	}
}

// NewApplication supports the installation of a new application.
func NewApplication(w http.ResponseWriter, r *http.Request) {

}
