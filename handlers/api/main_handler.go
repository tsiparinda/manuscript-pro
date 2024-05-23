// Package api provides a set of HTTP APIs for handling various server responses.
package api

import (
	"github.com/vedicsociety/platform/http/actionresults"
)

// StatusCodeResult is a struct that holds an HTTP status code.
type StatusCodeResult struct {
	// code represents the HTTP status code to be returned in the response.
	code int
}

// Execute sets the HTTP status code of the ResponseWriter in the provided ActionContext to the code
// stored in the StatusCodeResult instance. It returns nil if the operation is successful.
func (action *StatusCodeResult) Execute(ctx *actionresults.ActionContext) error {
	ctx.ResponseWriter.WriteHeader(action.code)
	return nil
}
