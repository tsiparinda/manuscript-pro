package api

import (
	"brucheion/models"

	"github.com/vedicsociety/platform/authorization/identity"
	"github.com/vedicsociety/platform/http/actionresults"
	"github.com/vedicsociety/platform/http/handling"
	"github.com/vedicsociety/platform/logging"
)

// ShareCollectionHandler is a struct that encapsulates methods for handling operations related to sharing a collection.
type ShareCollectionHandler struct {
	Repository models.Repository
	identity.User
	logging.Logger
	handling.URLGenerator
}

// GetShareCollection retrieves the shared collection by its ID.
// It checks the user's ID, and if it is invalid, an error message is returned.
// If the collection exists and the user has access, it is loaded and returned in JSON format.
func (h ShareCollectionHandler) GetShareCollection(colid int) actionresults.ActionResult {
	userid := h.User.GetID()
	h.Logger.Debugf("api.CollectionHandler.GetShareCollection: userid", userid)
	if userid == 0 {
		return actionresults.NewErrorJsonAction("Sorry, but your query is bad or you have not access to asked collection")
	}
	sharecol := h.Repository.LoadCollectionForShare(colid, userid)
	h.Logger.Debugf("api.CollectionHandler.GetShareCollection: sharecol", sharecol)
	if sharecol.Collection.Id == 0 {
		h.Logger.Info("api.CollectionHandler.GetShareCollection: No collection")
		return actionresults.NewErrorJsonAction("Sorry, but your query is bad or you have not access to asked collection")
	}

	return actionresults.NewJsonAction(sharecol)
}

// PostShareCollection updates the shared collection with the provided data.
// It logs the input, gets the user's ID and the collection ID, and loads the current state of the collection from the database.
// If there are any changes to be made, they are saved in the repository.
// If there is an error during this process, it is logged and returned as a JSON response.
func (h ShareCollectionHandler) PostShareCollection(p models.ShareCollection) actionresults.ActionResult {
	h.Logger.Debugf("PostShareCollection input: ", p)

	userid := h.User.GetID()
	colid := p.Id
	// get the current state of models.ShareCollection from DB
	sharecol := h.Repository.LoadCollectionForShare(colid, userid)
	// compare and update public flag
	if sharecol.Collection.IsPublic != p.Collection.IsPublic {
		// update Collection.IsPublic flag
		if err := h.Repository.SaveCollection(p.Collection, userid); err != nil {

			h.Logger.Debugf("PostShareCollection.SaveCollection error: ", err.Error())
			return actionresults.NewErrorJsonAction(err.Error())
		}
	}
	//	h.Logger.Debugf("PostShareCollection. sharecol.ColUsers != nil: ", p.ColUsers[0].Id_Col, len(p.ColUsers))
	if len(p.ColUsers) > 0 {
		h.Logger.Debugf("PostShareCollection. sharecol.ColUsers != nil: ")
		// delete userscol's
		err := h.Repository.DropCollectionUsers(colid, userid)
		if err != nil {
			return actionresults.NewErrorJsonAction(err.Error())
		}
		if len(p.ColUsers) > 0 {
			h.Logger.Debugf("PostShareCollection. sharecol.ColUsers != nil: ")
			// delete userscol's
			err := h.Repository.DropCollectionUsers(colid, userid)
			if err != nil {
				return actionresults.NewErrorJsonAction(err.Error())
			}
	
			// insert  userscol's
			for _, cu := range p.ColUsers {
				err = h.Repository.AddColUsers(cu, userid)
				if err != nil {
					return actionresults.NewErrorJsonAction(err.Error())
				}
			}
		}		// insert  userscol's
		for _, cu := range p.ColUsers {
			err = h.Repository.AddColUsers(cu, userid)
			if err != nil {
				return actionresults.NewErrorJsonAction(err.Error())
			}
		}
	}
	return actionresults.NewRedirectAction("/")
}

// DeleteCollectionsUser deletes the association between a collection and a user.
// It logs the input and attempts to remove the user from the collection.
// If there is an error during this process, it is returned as a JSON response.
func (h ShareCollectionHandler) DeleteCollectionsUser(p struct {
	ColId  int
	UserId int
}) actionresults.ActionResult {
	h.Logger.Debugf("DeleteCollectionsUser input: ", p)
	err := h.Repository.DropCollectionsUser(p.ColId, p.UserId)
	if err != nil {
		return actionresults.NewErrorJsonAction(err.Error())
	}

	// userid := h.User.GetID()
	// colid := p.Id
	return actionresults.NewJsonAction(p)
}
