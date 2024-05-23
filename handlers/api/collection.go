package api

import (
	"brucheion/models"

	"github.com/vedicsociety/platform/authorization/identity"
	"github.com/vedicsociety/platform/http/actionresults"
	"github.com/vedicsociety/platform/http/handling"
	"github.com/vedicsociety/platform/logging"
)

// CollectionHandler is a struct that encapsulates methods for handling operations related to a single collection.
type CollectionHandler struct {
	Repository models.Repository
	identity.User
	logging.Logger
	handling.URLGenerator
}

// GetCollection is a method that retrieves a specific collection based on the collection id provided.
// It verifies the user's ID first, then attempts to load the collection from the repository.
// If the collection does not exist or the user does not have access, an error message is returned.
// Otherwise, a JSON response containing the collection data is returned.
func (h CollectionHandler) GetCollection(colid int) actionresults.ActionResult {
	userid := h.User.GetID()
	//h.Logger.Debugf("api.CollectionHandler.GetCollection: userid", userid)
	if userid == 0 {
		return actionresults.NewErrorJsonAction("Sorry, but your query is bad or you have not access to asked collection")
	}
	collection := h.Repository.LoadCollection(colid, userid)
	//h.Logger.Debugf("api.CollectionHandler.GetCollection: collection", collection)
	if collection.Id == 0 {
		//	h.Logger.Info("api.CollectionHandler.GetCollection: No collection")
		return actionresults.NewErrorJsonAction("Sorry, but your query is bad or you have not access to asked collection")
	}

	resp := models.JSONResponse{
		Status:  "success",
		Message: "",
		Data:    collection, // identity.User{Name: user},
	}

	return actionresults.NewJsonAction(resp)
}
