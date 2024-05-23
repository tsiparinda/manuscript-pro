// Package root provides API handlers for managing and viewing collections.
package root

import (
	"brucheion/models"

	"github.com/vedicsociety/platform/http/actionresults"
	"github.com/vedicsociety/platform/http/handling"
	"github.com/vedicsociety/platform/logging"
)

// CollectionOverviewHandler is a struct that embeds the necessary dependencies for handling 
// requests related to the overview of a collection.
// It implements methods that handle HTTP requests for viewing collections.
// The struct embeds a Repository from the models package for interacting with the database,
// a URLGenerator from the handling package to generate URLs for the collections,
// and a Logger from the logging package to log the operations.
type CollectionOverviewHandler struct {
	models.Repository
	handling.URLGenerator
	logging.Logger
}

// GetView is a method on CollectionOverviewHandler that handles a GET request to view a collection.
// It logs the operation and then returns an ActionResult that represents a server-side rendered 
// template response. The template file used is 'root_svelte.html', and it receives a struct with 
// two fields, 'ColId' and 'URN', which are set to the passed collection id (colid) and URN respectively.
func (h CollectionOverviewHandler) GetView(colid int, urn string) actionresults.ActionResult {
	// handler.Logger.Debugf("PassageOverviewHandler.urn:", handler.MiddlewareComponent())
	h.Logger.Debugf("CollectionOverviewHandler.colid:", colid)
	return actionresults.NewTemplateAction("root_svelte.html",
		struct {
			ColId int
			URN string
		}{
			ColId: colid,
			URN: urn,
		})
}

