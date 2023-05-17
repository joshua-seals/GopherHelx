package v1

import (
	"fmt"
	"net/http"
)

// This collection of methods will standardize our errorResponses
// as well as error logging that clients see. logError, errorResponse, &
// serverErrorResponse are currently unexported.

func (c *CoreHandler) logError(r *http.Request, err error, message any) {
	c.Log.Errorln(err, "url-string: ", r.URL.String())
}

func (c *CoreHandler) errorResponse(w http.ResponseWriter, r *http.Request, status int, message any) {
	env := envelope{"error": message}

	err := c.writeJSON(w, status, env, nil)
	if err != nil {
		c.logError(r, err, nil)
		w.WriteHeader(500)
	}

}

func (c *CoreHandler) serverErrorResponse(w http.ResponseWriter, r *http.Request, err error) {
	c.logError(r, err, nil)
	message := "the server encountered a problem and could not process request"
	c.errorResponse(w, r, http.StatusInternalServerError, message)
}

// NotFoundRequest is exported to routes.go to handle notfound cases.
func (c *CoreHandler) NotFoundResponse(w http.ResponseWriter, r *http.Request) {
	message := "the requested resource could not be found"
	c.errorResponse(w, r, http.StatusNotFound, message)
}

// MethodNotAllowedResponse is exported to routes.go
// to handle MethodNotAllowed cases in a structured way.
func (c *CoreHandler) MethodNotAllowedResponse(w http.ResponseWriter, r *http.Request) {
	message := fmt.Sprintf("the %s method is not supported for this resource", r.Method)
	c.errorResponse(w, r, http.StatusMethodNotAllowed, message)
}
