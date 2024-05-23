package api

import (
	"brucheion/models"
	"net/http"
	"strconv"

	"github.com/vedicsociety/platform/authorization/identity"
	"github.com/vedicsociety/platform/http/actionresults"
	"github.com/vedicsociety/platform/http/handling"
	"github.com/vedicsociety/platform/http/handling/params"
	"github.com/vedicsociety/platform/logging"
)

// AddColHandler is a struct that encapsulates methods for handling operations related to adding a new collection.
type AddColHandler struct {
	Repository models.Repository
	identity.User
	logging.Logger
	handling.URLGenerator
}

// PostAddCollection adds a new collection using the input parameters and files provided.
// If the collection addition fails, it returns an HTTP 400 status code.
// If the addition is successful, it returns a success message in the response along with the ID of the created collection
func (h AddColHandler) PostAddCollection(params params.InputParams, files params.Files) actionresults.ActionResult {
	h.Logger.Debugf("api.PostAddCollection: params: %v", params.InputParam[0].Value[0], files)
	// get userid
	userid := h.User.GetID()
	title := params.InputParam[0].Value[0]
	colid, err := h.Repository.AddNewCollection(params, userid)
	if err != nil {
		h.Logger.Debugf("api.PostAddCollection: id_col: %v", colid)
		return &StatusCodeResult{http.StatusBadRequest}
	}
	resp := models.JSONResponse{
		Status:  "success",
		Message: "File was inserted successfully! Collection named " + title + " was created with ID " + strconv.Itoa(colid) + "",
		Data:    colid,
	}
	return actionresults.NewJsonAction(resp)
}
