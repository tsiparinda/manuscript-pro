package root

import (
	"brucheion/models"

	"github.com/vedicsociety/platform/http/actionresults"
	"github.com/vedicsociety/platform/http/handling"
	"github.com/vedicsociety/platform/logging"
)

type ShareCollectionHandler struct {
	models.Repository
	handling.URLGenerator
	logging.Logger
}

func (h ShareCollectionHandler) GetShareCollection(colid int) actionresults.ActionResult {
	h.Logger.Debugf("ShareCollectionHandler.colid:", colid)
	return actionresults.NewTemplateAction("root_svelte.html",
		struct {
			ColId int
		}{
			ColId: colid,
		})
}

