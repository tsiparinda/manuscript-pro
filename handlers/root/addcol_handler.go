// Package root provides API handlers for managing collections.
package root

import (
	"brucheion/models"

	"github.com/vedicsociety/platform/http/actionresults"
	"github.com/vedicsociety/platform/http/handling"
	"github.com/vedicsociety/platform/logging"
)


// AddColHandler is a struct that embeds the necessary dependencies for adding a new collection.
// It implements methods that handle HTTP requests for adding collections.
// The struct embeds a Repository from the models package for interacting with the database,
// a URLGenerator from the handling package to generate URLs for the collections,
// and a Logger from the logging package to log the operations.
type AddColHandler struct {
	models.Repository
	handling.URLGenerator
	logging.Logger
}

// GetAddCollection is a method on AddColHandler that handles a GET request for adding a new collection.
// It logs the operation and then returns a ActionResult that represents a server-side rendered 
// template response. The template file used is 'root_svelte.html', and it receives a struct with 
// a single field, 'TitleWin', that is set to 'AddNewCollection'.
func (h AddColHandler) GetAddCollection() actionresults.ActionResult {
	h.Logger.Debugf("GetAddNewCollection!!!")
	return actionresults.NewTemplateAction("root_svelte.html",
		struct {
			TitleWin string
		}{
			TitleWin: "AddNewCollection",
		})
}
