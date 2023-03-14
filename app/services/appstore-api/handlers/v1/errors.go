package v1

import "net/http"

var (
	ErrInternalServer = http.StatusInternalServerError
)

// This collection of methods will standardize our ErrorResponses
// as well as error logging that clients see.

func (c *CoreHandler) logError(r *http.Request, err error) {
	c.Log.Errorln(err)
}
