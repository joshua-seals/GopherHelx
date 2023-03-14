// All PUT method handlers are located herein which service the
// appstore-api
package v1

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/joshua-seals/gopherhelx/app/business/data/models"
)

var apps = map[int]string{1: "Webtop", 2: "Filebrowser", 3: "Jupyter", 4: "Balsam", 5: "PGAdmin"}

// AddToDashboard installs a selected app from app/list into user dashboard.
// This action also triggers an update to the user db table
// where the app is then added to specific user dashboard.
func (c CoreHandler) AddToDashboard(w http.ResponseWriter, r *http.Request) {
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
func (c CoreHandler) AppList(w http.ResponseWriter, r *http.Request) {
	appsList := apps
	w.Header().Set("Content-Type", "application/json")
	err := json.NewEncoder(w).Encode(appsList)
	if err != nil {
		c.Log.Error(err)
	}
}

// NewApplication supports the installation of a new application.
func (c CoreHandler) AddNewApplication(w http.ResponseWriter, r *http.Request) {

	newApp := models.Application{}
	err := json.NewDecoder(r.Body).Decode(&newApp)
	if err != nil {

		c.Log.Errorln(err, "Decoding New App ", newApp)
	}
	// Add contexting information
	ctx := context.TODO()
	appId, err := newApp.AddNewApplication(ctx, c.DB)
	if err != nil {
		status := http.StatusInternalServerError
		c.logError(r, err)
		json.NewEncoder(w).Encode(status)
	}

	w.Header().Set("Content-Type", "application/json")
	status := http.StatusAccepted
	response := fmt.Sprintf("status: %d, app_id: %s", status, appId)
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		status :=
			c.Log.Errorln("Error: ", err, "status: ", http.StatusInternalServerError, "Encoding New App ID Response", response)
	}
}
