// Package api provides a set of HTTP APIs to handle passage-related operations
// in a RESTful style.
package api

import (
	"brucheion/gocite"
	"brucheion/models"
	"encoding/json"
	"errors"
	"fmt"
	"sort"
	"strings"

	//"github.com/ThomasK81/gocite"
	"github.com/vedicsociety/platform/authorization/identity"
	"github.com/vedicsociety/platform/http/actionresults"
	"github.com/vedicsociety/platform/logging"
)

// PassageHandler is a struct that has necessary dependencies for handling
// passage-related requests. It implements various methods to fetch and manipulate
// passages in a collection.
type PassageHandler struct {
	Repository models.Repository
	identity.User
	logging.Logger
}

// GetCollectionUserRights checks and returns the user's rights for a given collection.
// It returns a JSONActionResult containing a JSON object with the user's rights.
func (h PassageHandler) GetCollectionUserRights(colid int) actionresults.ActionResult {

	userid := h.User.GetID()

	isEdit, err := h.Repository.IsCollectionWriteble(colid, userid)
	if err != nil {
		h.Logger.Info("api.PassageHandler.GetCollectionUserRights: No collection's buckets!")
		return actionresults.NewErrorJsonAction("Sorry, but your query is bad or you have not access to asked collection")
	}

	resp := models.JSONResponse{
		Status:  "success",
		Message: "",
		Data: struct {
			CanEditTranscription bool `json:"canEditTranscription"`
			CanEditReference     bool `json:"canEditReference"`
			CanEditMetadata      bool `json:"canEditMetadata"`
		}{
			CanEditTranscription: isEdit,
			CanEditReference:     isEdit,
			CanEditMetadata:      isEdit,
		},
	}

	return actionresults.NewJsonAction(resp)
}

// GetPassage retrieves a passage identified by its URN from a collection identified by colid.
// It returns a JSONActionResult containing a JSON object with the passage data.
func (h PassageHandler) GetPassage(colid int, urn string) actionresults.ActionResult {
	h.Logger.Debugf("api.PassageHandler.GetPassage: colid, urn, userid, user ", colid, urn)
	//(EXTRA int=10, string=urn:cts:sktlit:skt0001.nyaya002.J1D:3.1.2)
	// or (EXTRA int=10, string=undefined)

	// get userid
	userid := h.User.GetID()

	h.Logger.Debugf("api.PassageHandler.GetPassage: colid, urn, userid, user ", colid, urn, userid)

	// receive an all of buckets in database
	buckets := h.Repository.SelectCollectionBuckets(colid, userid)
	if len(buckets) == 0 {
		h.Logger.Info("api.PassageHandler.GetPassage: No collection's buckets!")
		// /return actionresults.NewErrorJsonAction("Sorry, but your query is bad or you have not access to asked collection")
		// render empty collection
		p := models.Passage{}
		resp := models.JSONResponse{
			Status:  "success",
			Message: "",
			Data:    p,
		}

		return actionresults.NewJsonAction(resp)
	}

	// set the bucket default as first accessible user's bucket

	sort.Strings(buckets)
	if urn == "default" || urn == "undefined" {
		//urn = "urn:cts:sktlit:skt0001.nyaya002.M3D:5.1.1"
		urn = buckets[0]
	}
	h.Logger.Debugf("api.PassageHandler.GetPassage: buckets ", buckets, urn)

	p, err := h.retrivePassage(colid, urn, userid)
	if err != nil {
		return actionresults.NewErrorAction(fmt.Errorf("api.PassageHandler.GetPassage: Internal server error1 %v", err))
	}
	p.TextRefs = buckets // !!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!

	// generate responce in json format
	resp := models.JSONResponse{
		Status:  "success",
		Message: "",
		Data:    p,
	}

	return actionresults.NewJsonAction(resp)
}

// GetPassages retrieves all passages from a collection identified by colid.
// It returns a JSONActionResult containing a JSON object with an array of passage data.
func (h PassageHandler) GetPassages(colid int) actionresults.ActionResult {
	h.Logger.Debugf("api.PassageHandler.GetPassage: colid, urn, userid, user ", colid)
	//(EXTRA int=10, string=urn:cts:sktlit:skt0001.nyaya002.J1D:3.1.2)
	// or (EXTRA int=10, string=undefined)

	// get userid
	userid := h.User.GetID()

	h.Logger.Debugf("api.PassageHandler.GetPassage: colid, urn, userid, user ", colid, userid)

	// receive an all of buckets in database
	buckets := h.Repository.SelectCollectionBuckets(colid, userid)
	if len(buckets) == 0 {
		h.Logger.Info("api.PassageHandler.GetPassage: No collection's buckets!")
		return actionresults.NewErrorJsonAction("Sorry, but your query is bad or you have not access to asked collection")
	}
	var passages []models.Passage
	// set the bucket default as first accessible user's bucket
	for _, urn := range buckets {
		p, err := h.retrivePassage(colid, urn, userid)
		if err != nil {
			return actionresults.NewErrorAction(fmt.Errorf("api.PassageHandler.GetPassage: Internal server error2 ", err, urn))
		}
		// p.TextRefs=buckets
		passages = append(passages, p)
	}

	// generate responce in json format
	resp := models.JSONResponse{
		Status:  "success",
		Message: "",
		Data:    passages,
	}

	return actionresults.NewJsonAction(resp)
}

// GetS3CollectionImages retrieves all images from an S3 collection identified by colid.
// It returns a JSONActionResult containing a JSON object with an array of image data.
func (h PassageHandler) GetS3CollectionImages(colid int) actionresults.ActionResult {

	userid := h.User.GetID()

	dict := h.Repository.LoadCollectionImageDictionary(colid, userid)

	image := []models.Image{}
	//imagecol := []models.ImageCollection{}
	h.Logger.Debugf("api.GetCollectionImages: select user bucketdict: (len(dist), urn) ", dict)

	for _, k := range dict {
		newimagecol := []models.Image{}
		err := json.Unmarshal([]byte(k.Value), &newimagecol)
		if err != nil {
			h.Logger.Debugf("api.GetCollectionImages: val ", err.Error())
		}
		h.Logger.Debugf("api.GetCollectionImages: val ", k.Key, k.Value, newimagecol)
		for _, im := range newimagecol {
			image = append(image, im)
		}
	}

	resp := models.JSONResponse{
		Status:  "success",
		Message: "",
		Data:    image,
	}
	return actionresults.NewJsonAction(resp)
}

// GetCollectionImages retrieves all images from a collection identified by colid.
// It returns a JSONActionResult containing a JSON object with an array of image data.
// from old Brucheion, image.go, requestImgCollection
// prints a list of the collections (as a keys from hstore) in the image collection
func (h PassageHandler) GetCollectionImages(colid int) actionresults.ActionResult {

	userid := h.User.GetID()

	dict := h.Repository.LoadCollectionImageDictionary(colid, userid)

	image := []models.Image{}
	//imagecol := []models.ImageCollection{}
	h.Logger.Debugf("api.GetCollectionImages: select user bucketdict: (len(dist), urn) ", dict)

	for _, k := range dict {
		newimagecol := []models.Image{}
		err := json.Unmarshal([]byte(k.Value), &newimagecol)
		if err != nil {
			h.Logger.Debugf("api.GetCollectionImages: val ", err.Error())
		}
		h.Logger.Debugf("api.GetCollectionImages: val ", k.Key, k.Value, newimagecol)
		for _, im := range newimagecol {
			image = append(image, im)
		}
	}

	resp := models.JSONResponse{
		Status:  "success",
		Message: "",
		Data:    image,
	}
	return actionresults.NewJsonAction(resp)
}

// retrivePassageBucketWork is a helper function that retrieves a work by its URN from a collection identified by colid.
// It returns a sorted gocite.Work object and an error if any occurred during the operation.
func (h PassageHandler) retrivePassageBucketWork(colid int, urn string, userid int) (result gocite.Work, err error) {

	// take all passages for this urn
	// returns models.BucketDict{key, value}, one key in short form, all another as urn:cts:sktlit:skt0001.nyaya002.J2D:3.2.16
	dict := h.Repository.SelectCollectionBucketDictionary(colid, urn, userid)
	result.WorkID = urn
	h.Logger.Debugf("api.PassageHandler.retriveCollectionBucketWork: select user bucketdict: (len(dist), urn) ", len(dict), result.WorkID)

	for _, pair := range dict {
		var passage gocite.Passage
		err := json.Unmarshal([]byte(pair.Value), &passage) //unmarshal the buffer and save the gocite.Passage
		if err != nil {
			return result, fmt.Errorf("api.PassageHandler.retriveCollectionBucketWork: Error unmarshalling Passage: %s", err)
		}

		if passage.PassageID != "" {
			result.Passages = append(result.Passages, passage)
		}
	}
	return gocite.SortPassages(result)
}

// retrivePassage retrieves a passage and its related metadata by its URN from a collection identified by colid.
// It returns a models.Passage object containing the passage data and an error if any occurred during the operation.
func (h PassageHandler) retrivePassage(colid int, urn string, userid int) (p models.Passage, e error) {

	// check urn
	if !gocite.IsCTSURN(urn) {
		e = errors.New("api.PassageHandler.retrivePassage: Bad urn request")
		return p, e
	}

	// cut the end of URN for receive header
	// urn:cts:sktlit:skt0001.nyaya002.M3D:3.1.1 -> urn:cts:sktlit:skt0001.nyaya002.M3D:
	urn_short := strings.Join(strings.Split(urn, ":")[0:4], ":") + ":"

	// receive all of passages from user's bucket in sorted state
	work, err := h.retrivePassageBucketWork(colid, urn_short, userid) /// a lot of work!!! return all data only for take First.PassageID
	// h.Logger.Debugf("api.PassageHandler.retrivePassage work: ", work.WorkID, work.First, work.Last, work.Ordered, work.Passages[0])

	if err != nil {
		//http.Error(w, "Internal server error", 500)
		e = err
		return p, e
	}

	var schemes []string
	// get the list of all schemas in urn_short
	for _, s := range work.Passages {
		lastIndex := strings.LastIndex(s.PassageID, ":")
		schemes = append(schemes, s.PassageID[lastIndex+1:])
	}

	//h.Logger.Infof("api.PassageHandler.retrivePassage:", strings.LastIndex(urn, ":")+1, len(urn), urn, work.First.PassageID)
	// correct urn if it's short (from undefined moved to short form and now we take first passage in a long form with :3.1.1)
	if strings.LastIndex(urn, ":")+1 == len(urn) {
		urn = work.First.PassageID
	}

	// receive a passage: key (urn....:x.y.z) -> value
	d, err := h.Repository.SelectCollectionBucketKeyValue(colid, urn_short, urn, userid)
	if err != nil {
		//http.Error(w, "Internal server error", 500)
		e = err
		return p, e
		// return actionresults.NewErrorAction(fmt.Errorf("api.PassageHandler.retrivePassage: Internal server error1 %v", err))
	}

	// receive a header
	c, err := h.Repository.SelectCollectionBucketKeyValue(colid, urn_short, urn_short, userid)
	if err != nil {
		e = err
		return p, e
		//http.Error(w, "Internal server error", 500)
		// return actionresults.NewErrorAction(fmt.Errorf("api.PassageHandler.retrivePassage: Internal server error2 %v", err))
	}

	catalog := models.BoltCatalog{}
	catalog.ColId = colid
	passage := gocite.Passage{}
	json.Unmarshal([]byte(d.JSON), &passage)
	json.Unmarshal([]byte(c.JSON), &catalog)

	// split passage lines to passages array
	text := passage.Text.TXT
	//h.Logger.Debugf("api.PassageHandler.retrivePassage: passage.Text.TXT:", text)

	var passages []string
	for _, str := range strings.Split(text, "\r\n") {
		if str != "" {
			passages = append(passages, str)
		}
	}

	//h.Logger.Debugf("work:", work.First.PassageID, work.Last.PassageID)
	var imageRefs []string
	for _, tmp := range passage.ImageLinks {
		imageRefs = append(imageRefs, tmp.Object)
	}

	p = models.Passage{
		ColId:              colid,
		PassageID:          passage.PassageID,      // current long urn
		Transcriber:        "user",                 // user name SHOULD BU AUTHORID
		TranscriptionLines: passages,               // for current urn
		PreviousPassage:    passage.Prev.PassageID, // for current urn
		NextPassage:        passage.Next.PassageID, // for current urn
		FirstPassage:       work.First.PassageID,   // for current urn first number x.y.z
		LastPassage:        work.Last.PassageID,    // for current urn last number x.y.z
		ImageRefs:          imageRefs,              // for current urn
		Schemes:            schemes,
		Catalog:            catalog, // header of current urn
		Text:               text,    // transcription's full text for blockeditor
		// TextRefs:           buckets,                // array all of users urns
	}
	// h.Logger.Debugf("api.PassageHandler.retrivePassage: passage:", p.PassageID, p.Transcriber, p.FirstPassage,
	// 	p.LastPassage, p.NextPassage, p.PreviousPassage, p.ImageRefs, p.Catalog, p.TranscriptionLines[0], p.Text)
	return p, nil
}
