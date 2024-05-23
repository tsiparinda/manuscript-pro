package api

import (
	"brucheion/models"
	"brucheion/utils"
	"strconv"

	//"github.com/ThomasK81/gocite"
	"github.com/vedicsociety/platform/authorization/identity"
	"github.com/vedicsociety/platform/http/actionresults"
	"github.com/vedicsociety/platform/http/handling"
	"github.com/vedicsociety/platform/logging"
)

// CollectionsHandler is a struct that encapsulates methods for handling operations related to collections.
type CollectionsHandler struct {
	Repository   models.Repository
	URLGenerator handling.URLGenerator
	identity.User
	logging.Logger
}

// GetCollections is a method that retrieves a list of collections based on the author id provided.
// It generates a list of collections for a particular user associated with an author id.
// Additional processing is done on each collection to generate specific URLs and fill sample names.
// Finally, a JSON response is created and returned with the collection data.
func (h CollectionsHandler) GetCollections(authorid int) actionresults.ActionResult {

	// get userid
	userid := h.User.GetID()
	user := h.User.GetDisplayName()
	h.Logger.Debugf("api.CollectionHandler.GetCollection: colid, urn, userid, user ", userid, user)

	colls, _ := h.Repository.LoadCollectionsPageAuthor(userid, authorid, 1, 1000)
	h.Logger.Debugf("api.CollectionHandler.GetCollection: colid, urn, userid, user ", colls)

	for c, _ := range colls {
		// colls[c].CollectionURL, _ = h.URLGenerator.GenerateUrl(root.CollectionOverviewHandler.GetView, colls[c].Id)
		colls[c].CollectionURL = "/view/" + strconv.Itoa(colls[c].Collection.Id) + "/home/"
		colls[c].Author.AuthorURL, _ = h.URLGenerator.GenerateUrl(CollectionsHandler.GetCollections, colls[c].Author.Id, 1)
		helper := utils.HandlersHelper{
			Repository:   h.Repository,
			URLGenerator: h.URLGenerator,
			User:         h.User,
			Logger:       h.Logger,
		}
		utils.HandlersHelper.FillSampleNameCollection(helper, &colls[c], userid)
	}

	h.Logger.Debugf("api.CollectionHandler.GetCollection: colls.size ", len(colls), colls[0].CollectionURL)

	// generate responce in json format
	resp := models.JSONResponse{
		Status:  "success",
		Message: "",
		Data:    colls,
	}

	return actionresults.NewJsonAction(resp)
}
