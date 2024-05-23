// Package root provides API handlers for managing and viewing collections.
package root

import (
	"brucheion/handlers/api"
	"brucheion/models"
	"brucheion/utils"
	"math"

	"github.com/vedicsociety/platform/authorization/identity"
	"github.com/vedicsociety/platform/http/actionresults"
	"github.com/vedicsociety/platform/http/handling"
	"github.com/vedicsociety/platform/logging"
)

const pageSize = 1000

// CollectionsHandler is a struct that embeds the necessary dependencies for handling requests
// related to multiple collections. It implements methods that handle HTTP requests for retrieving
// collections. The struct embeds a Repository from the models package for interacting with the database,
// a URLGenerator from the handling package to generate URLs for the collections,
// a User from the identity package to identify the current user, and a Logger from the logging package to log the operations.
type CollectionsHandler struct {
	Repository   models.Repository
	URLGenerator handling.URLGenerator
	identity.User
	logging.Logger
}

// CollectionsTemplateContext is a struct that represents the context to be passed to the template
// for rendering collections.
type CollectionsTemplateContext struct {
	Collections    []models.CollectionPage
	Page           int
	PageCount      int
	PageNumbers    []int
	PageUrlFunc    func(int) string
	SelectedAuthor int
	User           *identity.User
}

// GetCollections is a method on CollectionsHandler that handles a GET request for retrieving a page of collections.
// It logs the operation and then returns an ActionResult that represents a server-side rendered
// template response. The template file used is 'root_collections.html', and it receives a struct of type
// CollectionsTemplateContext which encapsulates information about the collections and the current user.
func (h CollectionsHandler) GetCollections(author, page int) actionresults.ActionResult {
	userid := h.User.GetID()
	h.Logger.Debugf("root.GetCollections input ", userid, author, page)
	colls, total := h.Repository.LoadCollectionsPageAuthor(userid, author, page, pageSize)

	for c, _ := range colls {
		colls[c].CollectionURL, _ = h.URLGenerator.GenerateUrl(CollectionOverviewHandler.GetView, colls[c].Collection.Id, "default")
		colls[c].EditCollectionURL = "#"
		//colls[c].CanEditCollection = false
		colls[c].SharingCollectionURL = "#"
		colls[c].CanSharingCollection = false
		colls[c].DropCollectionURL = "#"
		colls[c].CanDropCollection = false

		if userid == colls[c].Author.Id {

			colls[c].SharingCollectionURL, _ = h.URLGenerator.GenerateUrl(api.ShareCollectionHandler.PostShareCollection, colls[c].Collection.Id)
			colls[c].CanSharingCollection = true
			colls[c].DropCollectionURL, _ = h.URLGenerator.GenerateUrl(api.DropCollectionHandler.PostDropCollection, colls[c].Collection.Id)
			colls[c].CanDropCollection = true
		}

		if colls[c].CanEditCollection == true {
			colls[c].EditCollectionURL, _ = h.URLGenerator.GenerateUrl(api.EditPassageHandler.PostSaveTranscription, colls[c].Collection.Id)
		}

		colls[c].Author.AuthorURL, _ = h.URLGenerator.GenerateUrl(CollectionsHandler.GetCollections, colls[c].Author.Id, 1)

		helper := utils.HandlersHelper{
			Repository:   h.Repository,
			URLGenerator: h.URLGenerator,
			User:         h.User,
			Logger:       h.Logger,
		}
		utils.HandlersHelper.FillSampleNameCollection(helper, &colls[c], userid)
	}
	h.Logger.Debugf("root.GetCollections colls ", colls)
	pageCount := int(math.Ceil(float64(total) / float64(pageSize)))
	return actionresults.NewTemplateAction("root_collections.html",
		CollectionsTemplateContext{
			Collections:    colls,
			Page:           page,
			PageCount:      pageCount,
			PageNumbers:    h.generatePageNumbers(pageCount),
			PageUrlFunc:    h.createPageUrlFunction(author),
			SelectedAuthor: author,
			User:           &h.User,
		})
}

// createPageUrlFunction is a method on CollectionsHandler that returns a function for generating
// URLs for different pages of collections. This function is intended to be passed to the template
// context so it can generate the URLs as needed.
func (h CollectionsHandler) createPageUrlFunction(author int) func(int) string {
	return func(page int) string {
		url, _ := h.URLGenerator.GenerateUrl(CollectionsHandler.GetCollections,
			author, page)
		return url
	}
}

// generatePageNumbers is a method on CollectionsHandler that generates a slice of integers
// representing page numbers. The number of pages is determined by the provided pageCount.
func (h CollectionsHandler) generatePageNumbers(pageCount int) (pages []int) {
	pages = make([]int, pageCount)
	for i := 0; i < pageCount; i++ {
		pages[i] = i + 1
	}
	return
}

// mustGenerateUrl is a helper function that generates a URL using the provided URLGenerator and
// the target and data parameters. If there is an error during URL generation, it will panic.
// This function is intended to be used when it is guaranteed that the provided data will not
// cause an error.
func mustGenerateUrl(generator handling.URLGenerator, target interface{}, data ...interface{}) string {
	url, err := generator.GenerateUrl(target, data)
	if err != nil {
		panic(err)
	}
	return url
}
