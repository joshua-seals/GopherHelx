package handlers

import (
	"net/http"

	"github.com/joshua-seals/gopherhelx/business/k8s"
)

// StopApp will delete the currently deployed application
// from the user table and dashboard.
func StopApp(w http.ResponseWriter, r *http.Request) {

}

// RemoveApp will uninstall the applicaiton from the user dashboard.
// Subsequently, the app is removed from the user db table purview.
func RemoveApp(w http.ResponseWriter, r *http.Request) {
	k8s.DeleteDeployment()
}
