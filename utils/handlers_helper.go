package utils

import (
	"brucheion/gocite"

	"brucheion/models"
	"encoding/json"
	"fmt"
	"sort"
	"strings"

	//"github.com/ThomasK81/gocite"
	"github.com/vedicsociety/platform/authorization/identity"
	"github.com/vedicsociety/platform/http/handling"
	"github.com/vedicsociety/platform/logging"
)

type HandlersHelper struct {
	Repository   models.Repository
	URLGenerator handling.URLGenerator
	identity.User
	logging.Logger
}

func (h HandlersHelper) FillSampleNameCollection(coll *models.CollectionPage, userid int) {
	//colls[c].CollectionURL, _ = h.URLGenerator.GenerateUrl(tools.ToolsHandler.GetSection, "CollectionOverview")

	colid := coll.Collection.Id

	h.Logger.Debugf("GetCollections colls ", coll)

	// receive an all of buckets in database
	textRefs := h.Repository.SelectCollectionBuckets(colid, userid)
	if len(textRefs) == 0 {
		h.Logger.Info("api.CollectionHandler.GetCollection: No collection's buckets!")
		return //actionresults.NewErrorAction(fmt.Errorf("api.CollectionHandler.GetCollection: No collection's buckets!"))
	}

	// set the bucket default as first accessible user's bucket
	//get the first bucket for sample
	sort.Strings(textRefs)
	urn := textRefs[0]

	h.Logger.Debugf("api.CollectionHandler.GetCollection: textRefs ", textRefs)

	// check urn
	if !gocite.IsCTSURN(urn) {
		return //actionresults.NewErrorAction(errors.New("api.CollectionHandler.GetCollection: Bad urn request"))
	}

	// cut the end of URN for receive header
	// urn:cts:sktlit:skt0001.nyaya002.M3D:3.1.1 -> urn:cts:sktlit:skt0001.nyaya002.M3D:
	bucketName := strings.Join(strings.Split(urn, ":")[0:4], ":") + ":"

	// receive all of passages from user's bucket in sorted state
	work, err := h.retriveCollectionsBucketWork(colid, bucketName, userid)
	if err != nil {
		//http.Error(w, "Internal server error", 500)
		return //actionresults.NewErrorAction(fmt.Errorf("api.CollectionHandler.GetCollection: Internal server error3 %v", err))
	}

	//h.Logger.Infof("api.PassageHandler.GetPassage:", strings.LastIndex(urn, ":")+1, len(urn), urn, work.First.PassageID)
	// correct urn if it's short (from undefined)
	if strings.LastIndex(urn, ":")+1 == len(urn) {
		urn = work.First.PassageID
	}

	// receive a passage: key (urn....:x.y.z) -> value
	pass, err := h.Repository.SelectCollectionBucketKeyValue(colid, bucketName, urn, userid)
	if err != nil {
		//http.Error(w, "Internal server error", 500)
		return //actionresults.NewErrorAction(fmt.Errorf("api.CollectionHandler.GetCollection: Internal server error1 %v", err))
	}
	// receive a header
	cat, err := h.Repository.SelectCollectionBucketKeyValue(colid, bucketName, bucketName, userid)
	if err != nil {
		//http.Error(w, "Internal server error", 500)
		return //actionresults.NewErrorAction(fmt.Errorf("api.CollectionHandler.GetCollection: Internal server error2 %v", err))
	}

	catalog := models.BoltCatalog{}
	passage := gocite.Passage{}
	json.Unmarshal([]byte(pass.JSON), &passage)
	json.Unmarshal([]byte(cat.JSON), &catalog)

	// split passage lines to passages array
	text := passage.Text.TXT
	passages := strings.Split(text, "\r\n")

	//take first line for sample of collectoin
	fline := passages[0]
	flinelen := len(fline)
	if flinelen > 100 {
		coll.SampleText = fline[0:99] + "..."
	} else {
		coll.SampleText = fline
	}
}

// SelectUserBucketWork retrieves an entire work from the users database as an (ordered) gocite.Work object
func (h HandlersHelper) retriveCollectionsBucketWork(colid int, urn string, userid int) (result gocite.Work, err error) {

	dict := h.Repository.SelectCollectionBucketDictionary(colid, urn, userid)
	result.WorkID = urn
	h.Logger.Debugf("api.CollectionHandler.retriveCollectionBucketWork: select user bucketdict: (len(dist), urn) ", len(dict), result.WorkID)

	for _, pair := range dict {
		var passage gocite.Passage
		err := json.Unmarshal([]byte(pair.Value), &passage) //unmarshal the buffer and save the gocite.Passage
		if err != nil {
			return result, fmt.Errorf("api.CollectionHandler.retriveCollectionBucketWork: Error unmarshalling Passage: %s", err)
		}

		if passage.PassageID != "" {
			result.Passages = append(result.Passages, passage)
		}
	}
	return gocite.SortPassages(result)
}
