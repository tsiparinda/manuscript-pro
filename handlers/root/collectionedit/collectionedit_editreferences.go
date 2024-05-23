package collectionedit

import (
	"brucheion/gocite"
	"brucheion/models"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/vedicsociety/platform/http/actionresults"
)

// GetEditReferences is a method of CollectionEditHandler, that 
// handles a GET request for editing references in a collection.
// It fetches a list of buckets from a collection for a given user id 
// and collection id, retrieves specific data from a particular bucket 
// based on a URN, unmarshals the JSON data into a gocite.Passage struct, 
// retrieves all passages from the user's bucket in sorted state, prepares 
// a models.Transcription struct with all the necessary details, loads 
// a page and renders a template with the name "root_editreferences.html", 
// passing in the transpage as the data.
//
// It receives two parameters:
// - colid: An integer that represents the collection ID.
// - urn: A string that represents the Uniform Resource Name.
//
// Returns:
// - An ActionResult which can be rendered to the user as a webpage.
func (h CollectionEditHandler) GetEditReferences(colid int, urn string) actionresults.ActionResult {

	userid := h.User.GetID()
	uname := h.User.GetDisplayName()
	h.Logger.Debugf("GetEditReferences.colid:", colid, urn, userid)
	kind := "/tools/editreferences/" + strconv.Itoa(colid) + "/"
	// get all buckets
	// receive an all of buckets in database
	textRefs := h.Repository.SelectCollectionBuckets(colid, userid)
	if len(textRefs) == 0 {
		h.Logger.Info("GetEditReferences: No collection's buckets!")
		return actionresults.NewErrorJsonAction("Sorry, but your query is bad or you have not access to asked collection")
	}

	bucketName := strings.Join(strings.Split(urn, ":")[0:4], ":") + ":"
	h.Logger.Debugf("GetEditReferences.requestedbucket:", bucketName)

	// adding testing if requestedbucket exists...
	//retrieveddata, _ := BoltRetrieve(dbname, bucketName, urn)

	retrieveddata, _ := h.Repository.SelectCollectionBucketKeyValue(colid, bucketName, urn, userid)
	// h.Logger.Debugf("GetEditCatalog.retrievedcat:", retrieveddata)
	//	retrievedWork, _ := BoltRetrieveWork(dbname, requestedbucket)
	retrievedPassage := gocite.Passage{}
	json.Unmarshal([]byte(retrieveddata.JSON), &retrievedPassage)

	// ctsurn := retrievedPassage.PassageID
	text := retrievedPassage.Text.TXT
	// previous := retrievedPassage.Prev.PassageID
	// next := retrievedPassage.Next.PassageID
	imageref := []string{}
	for _, tmp := range retrievedPassage.ImageLinks {
		imageref = append(imageref, tmp.Object)
	}
	/*First := retrievedPassage.First.PassageID
	last := retrievedPassage.Last.PassageID*/

	// receive all of passages from user's bucket in sorted state
	work, err := h.retriveCollectionBucketWork(colid, bucketName, userid)
	if err != nil {
		//http.Error(w, "Internal server error", 500)
		return actionresults.NewErrorAction(fmt.Errorf("api.CollectionHandler.GetCollection: Internal server error3 %v", err))
	}

	// first := work.First.PassageID
	// last := work.Last.PassageID

	imagejs := "urn:cite2:test:googleart.positive:DuererHare1502"
	switch len(imageref) > 0 {
	case true:
		imagejs = imageref[0]
	}
	pathname := "/tools/editreferences/" + strconv.Itoa(colid) + "/"
	transcription := models.Transcription{
		ColID:         colid,
		CTSURN:        retrievedPassage.PassageID,
		Transcriber:   uname,
		Transcription: text,
		Previous:      pathname + retrievedPassage.Prev.PassageID,
		Next:          pathname + retrievedPassage.Next.PassageID,
		First:         pathname + work.First.PassageID,
		Last:          pathname + work.Last.PassageID,
		TextRef:       textRefs,
		ImageRef:      imageref,
		ImageJS:       imagejs}

	transpage, _ := h.loadPage(transcription, kind)
	h.Logger.Debugf("GetEditReferences.transcription:", transpage)
	return actionresults.NewTemplateAction("root_editreferences.html", transpage)
}


