// Package root provides API handlers for managing collections.
package root

import (
	"brucheion/models"

	"github.com/vedicsociety/platform/http/actionresults"
	"github.com/vedicsociety/platform/http/handling"
	"github.com/vedicsociety/platform/logging"
)

// CollectionBulkEditHandler is a struct that embeds the necessary dependencies for bulk editing a collection.
// It implements methods that handle HTTP requests for bulk operations on collections.
// The struct embeds a Repository from the models package for interacting with the database,
// a URLGenerator from the handling package to generate URLs for the collections,
// and a Logger from the logging package to log the operations.
type CollectionBulkEditHandler struct {
	models.Repository
	handling.URLGenerator
	logging.Logger
}

// GetBulkEdit is a method on CollectionBulkEditHandler that handles a GET request for bulk editing a collection.
// It logs the operation and then returns a ActionResult that represents a server-side rendered
// template response. The template file used is 'root_svelte.html', and it receives a struct with
// a single field, 'ColId', that is set to the passed collection id (colid).
func (h CollectionBulkEditHandler) GetBulkEdit(colid int) actionresults.ActionResult {
	// handler.Logger.Debugf("PassageOverviewHandler.urn:", handler.MiddlewareComponent())

	h.Logger.Debugf("CollectionBulkEditHandler.colid:", colid)
	return actionresults.NewTemplateAction("root_svelte.html",
		struct {
			ColId int
		}{
			ColId: colid,
		})
}
