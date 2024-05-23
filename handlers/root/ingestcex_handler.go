// Package root provides API handlers for managing collections.
package root

import (
	"brucheion/models"

	"github.com/vedicsociety/platform/http/actionresults"
	"github.com/vedicsociety/platform/http/handling"
	"github.com/vedicsociety/platform/logging"
)

// IngestCEXHandler is a struct that embeds the necessary dependencies for handling the ingestion of CEX files.
// It implements methods that handle HTTP requests related to this process.
// The struct embeds a Repository from the models package for interacting with the database,
// a URLGenerator from the handling package to generate URLs for the collections,
// and a Logger from the logging package to log the operations.
type IngestCEXHandler struct {
	models.Repository
	handling.URLGenerator
	logging.Logger
}

// GetIngestCEX is a method on IngestCEXHandler that handles a GET request for ingesting a CEX file.
// It logs the operation and then returns an ActionResult that represents a server-side rendered 
// template response. The template file used is 'root_svelte.html', and it receives a struct with 
// a single field, 'TitleWin', which is set to 'IngestCEX'.
func (h IngestCEXHandler) GetIngestCEX() actionresults.ActionResult {
	h.Logger.Debugf("GetIngestCEX!!!")
	return actionresults.NewTemplateAction("root_svelte.html",
		struct {
			TitleWin string
		}{
			TitleWin: "IngestCEX",
		})
}
