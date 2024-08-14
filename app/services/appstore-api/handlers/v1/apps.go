// All PUT method handlers are located herein which service the
// appstore-api
package v1

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/joshua-seals/gopherhelx/app/business/data/models"
)

// swagger:route GET /app/list apps appList
// Returns list of applications available for use on helxplatform.
// responses:
//
//	200: appsResponse
func (c CoreHandler) AppList(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	apps, err := models.AppList(ctx, c.DB)

	if err != nil {
		c.serverErrorResponse(w, r, err)
		return
	}

	status := http.StatusAccepted
	data := envelope{"applist": apps}
	if err := c.writeJSON(w, status, data, nil); err != nil {
		c.serverErrorResponse(w, r, err)
		return
	}

}

// swagger:route POST /app/new apps newApplication
// Expects an Application definition
// Responds Success or Error Json Response
// responses:
//
//	200: newApplication
func (c CoreHandler) NewApplication(w http.ResponseWriter, r *http.Request) {

	newApp := models.Application{}
	err := json.NewDecoder(r.Body).Decode(&newApp)
	if err != nil {
		// c.logError(r, err, newApp)
		c.serverErrorResponse(w, r, err)
		return
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
		return
	}
}
