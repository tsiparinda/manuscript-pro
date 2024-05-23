// Package api provides API handlers for managing individual user data.
package api

import (
	"brucheion/models"

	"github.com/vedicsociety/platform/authorization/identity"
	"github.com/vedicsociety/platform/http/actionresults"
	"github.com/vedicsociety/platform/logging"
)

// UserHandler is a struct that embeds the necessary dependencies for handling individual user-related requests.
// It includes models.Repository for interacting with the data store,
// identity.User for accessing the user's identity, and logging.Logger for logging the operations.
type UserHandler struct {
	Repository models.Repository
	identity.User
	logging.Logger
}

// GetUser is a method on UserHandler that retrieves the identity of the current user.
// It retrieves the user's ID and display name from the User object, logs the operation,
// and then returns an ActionResult that represents a JSON response containing the username.
// The JSONResponse is defined in the models package and consists of a status, a message, and the data itself.
func (h UserHandler) GetUser() actionresults.ActionResult {
	// should be: {"status":"success","message":"","data":{"name":"albatros"}}
	// get userid
	userid := h.User.GetID()
	username := h.User.GetDisplayName()
	h.Logger.Debugf("Userid:", userid)
	resp := models.JSONResponse{
		Status:  "success",
		Message: "",
		Data:    username, // identity.User{Name: user},
	}

	return actionresults.NewJsonAction(resp)
}

