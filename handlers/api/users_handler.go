// Package api provides API handlers for managing users.
package api

import (
	"brucheion/models"

	"github.com/vedicsociety/platform/authorization/identity"
	"github.com/vedicsociety/platform/http/actionresults"
	"github.com/vedicsociety/platform/logging"
)

// UsersHandler is a struct that embeds the necessary dependencies for handling user-related requests.
// It implements methods that handle HTTP requests related to this process.
// The struct embeds a Repository from the models package for interacting with the database,
// an identity.User to access the user's identity, and a Logger from the logging package to log the operations.
type UsersHandler struct {
	Repository models.Repository
	identity.User
	logging.Logger
}

// GetUsers is a method on UsersHandler that handles a GET request to retrieve a user by userid.
// It retrieves the user data from the repository, logs the operation, and then returns an ActionResult
// that represents a JSON response containing the user data.
// The JSONResponse is defined in the models package and consists of a status, a message, and the data itself.
func (h UsersHandler) GetUsers(userid int) actionresults.ActionResult {
	// should be: {"status":"success","message":"","data":{"name":"albatros"}}
	// get userid

	//username := h.User.GetDisplayName()

	users := h.Repository.GetUsers(userid)
	h.Logger.Debugf("api.CollectionHandler.GetCollection: colid, urn, userid, user ", users)
	// generate responce in json format
	resp := models.JSONResponse{
		Status:  "success",
		Message: "",
		Data:    users,
	}
	return actionresults.NewJsonAction(resp)
}
