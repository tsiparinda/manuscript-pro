package api

import (
	"brucheion/models"
	"net/http"

	//"github.com/ThomasK81/gocite"
	"github.com/vedicsociety/platform/authorization/identity"
	"github.com/vedicsociety/platform/http/actionresults"
	"github.com/vedicsociety/platform/logging"
)

// DropCollectionHandler is a struct that encapsulates methods for handling operations related to deleting a collection
type DropCollectionHandler struct {
	Repository models.Repository
	identity.User
	logging.Logger
}

// PostDropCollection removes the collection identified by the response structure.
// It logs the incoming response and extracts the user ID from it.
// It then drops the collection by ID from the repository, verifying the user is the author of the collection.
// It returns an HTTP status code depending on the result of the operation.
func (h DropCollectionHandler) PostDropCollection(resp struct{ Colid int }) actionresults.ActionResult {
	h.Logger.Debugf("api.PostDropCollection: ", resp)
	userid := h.User.GetID()
	// only collection's author can drop collection!!!

	err := h.Repository.DropCollection(resp.Colid, userid)
	if err != nil {
		return &StatusCodeResult{http.StatusBadRequest}
	}

	return &StatusCodeResult{http.StatusOK}
}
