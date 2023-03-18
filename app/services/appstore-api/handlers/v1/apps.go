// All PUT method handlers are located herein which service the
// appstore-api
package v1

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/joshua-seals/gopherhelx/app/business/data/models"
)

// AddToDashboard installs a selected app from app/list into user dashboard.
// This action also triggers an update to the user db table
// where the app is then added to specific user dashboard.
func (c CoreHandler) AddToDashboard(w http.ResponseWriter, r *http.Request) {

}

// Applist shows the list of applications available
// for a user to install in their dashboard.
func (c CoreHandler) AppList(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	apps, err := models.AppList(ctx, c.DB)

	if err != nil {
		c.serverErrorResponse(w, r, err)
	}

	status := http.StatusAccepted
	data := envelope{"applist": apps}
	if err := c.writeJSON(w, status, data, nil); err != nil {
		c.serverErrorResponse(w, r, err)
	}

}

// NewApplication supports the installation of a new application.
func (c CoreHandler) NewApplication(w http.ResponseWriter, r *http.Request) {

	newApp := models.Application{}
	err := json.NewDecoder(r.Body).Decode(&newApp)
	if err != nil {
		// c.logError(r, err, newApp)
		c.serverErrorResponse(w, r, err)
	}
	// Add contexting information
	ctx := context.TODO()
	appId, err := newApp.AddNewApplication(ctx, c.DB)
	if err != nil {
		c.serverErrorResponse(w, r, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	status := http.StatusAccepted
	data := envelope{"success": appId}

	if err := c.writeJSON(w, status, data, nil); err != nil {
		c.serverErrorResponse(w, r, err)
	}
}
