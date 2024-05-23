// Package root provides API handlers for managing collections.
package root

import (
	"brucheion/models"

	"github.com/vedicsociety/platform/http/actionresults"
	"github.com/vedicsociety/platform/http/handling"
	"github.com/vedicsociety/platform/logging"
)

// IngestImageHandler is a struct that embeds the necessary dependencies for handling the ingestion of images.
// It implements methods that handle HTTP requests related to this process.
// The struct embeds a Repository from the models package for interacting with the database,
// a URLGenerator from the handling package to generate URLs for the collections,
// and a Logger from the logging package to log the operations.
type IngestImageHandler struct {
	models.Repository
	handling.URLGenerator
	logging.Logger
}

// GetIngestImage is a method on IngestImageHandler that handles a GET request for ingesting an image.
// It logs the operation and then returns an ActionResult that represents a server-side rendered 
// template response. The template file used is 'root_svelte.html', and it receives a struct with 
// a single field, 'ColId', which is set to the passed collection id (colid).
func (h IngestImageHandler) GetIngestImage(colid int) actionresults.ActionResult {
	// handler.Logger.Debugf("PassageOverviewHandler.urn:", handler.MiddlewareComponent())
	h.Logger.Debugf("GetIngestImage.colid:", colid)
	return actionresults.NewTemplateAction("root_svelte.html",
		struct {
			ColId int
		}{
			ColId: colid,
		})
}

