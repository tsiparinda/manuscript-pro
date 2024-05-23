// Package api provides a set of HTTP APIs to handle passage editing operations.
package api

import (
	"brucheion/models"
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/vedicsociety/platform/authorization/identity"
	"github.com/vedicsociety/platform/http/actionresults"
	"github.com/vedicsociety/platform/http/handling"
	"github.com/vedicsociety/platform/logging"
)

// EditPassageHandler is a struct that has necessary dependencies for handling 
// passage editing requests. It implements various methods to edit and save passages.
type EditPassageHandler struct {
	Repository models.Repository
	identity.User
	logging.Logger
	handling.URLGenerator
}

// GetEditTranscription loads a transcription for editing identified by colid.
// It returns a JSONActionResult containing a JSON object with the transcription data.
func (h EditPassageHandler) GetEditTranscription(colid int) actionresults.ActionResult {
	userid := h.User.GetID()
	return actionresults.NewJsonAction(h.Repository.LoadTranscriptionForEdit(colid, userid))
}

// PostSaveTranscription saves a provided transcription to the repository.
// It returns a RedirectActionResult pointing to the home view of the collection.
func (h EditPassageHandler) PostSaveTranscription(coltrans models.Transcription) actionresults.ActionResult {

	h.Logger.Debugf("PostEditTranscription input: ", coltrans)

	userid := h.User.GetID()

	h.Repository.SaveTranscription(coltrans, userid)

	return actionresults.NewRedirectAction("/view/" + strconv.Itoa(coltrans.ColID) + "/home/")
}

// PostSaveReference saves a provided reference to the repository.
// It returns a RedirectActionResult pointing to the home view of the collection.
func (h EditPassageHandler) PostSaveReference(coltrans models.Transcription) actionresults.ActionResult {

	h.Logger.Debugf("PostSaveReference input: ", coltrans)

	userid := h.User.GetID()

	h.Repository.SaveReference(coltrans, userid)

	return actionresults.NewRedirectAction("/view/" + strconv.Itoa(coltrans.ColID) + "/home/")
}

// PostSavePassageText saves the provided passage text to the repository.
// It returns a RedirectActionResult pointing to the home view of the collection.
func (h EditPassageHandler) PostSavePassageText(passagetext models.PassageText) actionresults.ActionResult {
	// Block editor
	h.Logger.Debugf("PostSavePassageText input: ", passagetext)

	userid := h.User.GetID()

	h.Repository.SavePassageText(passagetext, userid)

	return actionresults.NewRedirectAction("/view/" + strconv.Itoa(passagetext.ColId) + "/home/")
}

// PostSavePassage saves a provided passage transcription to the repository.
// It returns a JSONActionResult indicating successful operation.
func (h EditPassageHandler) PostSavePassage(passage models.Passage) actionresults.ActionResult {

	h.Logger.Debugf("PostSavePassage input: ", passage)

	userid := h.User.GetID()

	h.Repository.SavePassageTranscription(passage, userid)

	return actionresults.NewJsonAction("OK")
}

// PostSaveMetadata saves a provided catalog metadata to the repository.
// If an error occurs during the save operation, it returns an ErrorActionResult.
// Otherwise, it returns a RedirectActionResult pointing to the home view of the collection.
func (h EditPassageHandler) PostSaveMetadata(cat models.BoltCatalog) actionresults.ActionResult {
	h.Logger.Debugf("PostEditMetadata.TranscriptionCatalog: ", cat)

	catvalue, _ := json.Marshal(cat)

	//put bucket data (hstore)
	h.Logger.Debugf("PostEditMetadata: id_col, newbucket, catkey, catvalue", cat.ColId, cat.URN, cat.URN)

	err := h.Repository.SaveCiteDataDict(cat.ColId, cat.URN, cat.URN, catvalue)

	if err != nil {
		h.Logger.Debugf("PostEditMetadata: Cannot save metadata", err.Error())
		return &actionresults.ErrorActionResult{}
	}

	return actionresults.NewRedirectAction("/view/" + fmt.Sprintf("%v", cat.ColId) + "/home/")
}
