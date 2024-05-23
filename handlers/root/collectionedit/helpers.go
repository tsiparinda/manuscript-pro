package collectionedit

import (
	"brucheion/gocite"
	"brucheion/models"
	"encoding/json"
	"fmt"
	"html/template"
	"strconv"
	"strings"
)

// retriveCollectionBucketWork is a method of the CollectionEditHandler struct.
// It retrieves an entire work from the user's database as an ordered gocite.Work object.
// This function accepts a collection ID, a URN and a user ID, and it returns the ordered
// gocite.Work and any error encountered during the operation.
func (h CollectionEditHandler) retriveCollectionBucketWork(colid int, urn string, userid int) (result gocite.Work, err error) {

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

// loadPage is a method of the CollectionEditHandler struct.
// It prepares and returns a models.TranscriptionPage instance for a given models.Transcription.
// This function accepts a Transcription object and a 'kind' string which seems to be used to
// generate HTML elements within the function. It returns a models.TranscriptionPage and any
// error encountered during the operation.
// Note: This function is manipulating HTML directly, which can be unsafe. In a production
// environment, it is recommended to use templating systems and separate the presentation layer
// from the business logic to make your application safer and more maintainable.
func (h CollectionEditHandler) loadPage(transcription models.Transcription, kind string) (models.TranscriptionPage, error) {

	userid := h.User.GetID()
	colid := transcription.ColID
	user := transcription.Transcriber
	imagejs := transcription.ImageJS
	title := transcription.CTSURN
	text := transcription.Transcription
	previous := transcription.Previous
	next := transcription.Next
	first := transcription.First
	last := transcription.Last
	catid := transcription.CatID
	catcit := transcription.CatCit
	catgroup := transcription.CatGroup
	catwork := transcription.CatWork
	catversion := transcription.CatVers
	catexpl := transcription.CatExmpl
	caton := transcription.CatOn
	catlan := transcription.CatLan

	var previouslink, nextlink string
	switch {
	case previous == "":
		previouslink = `<a href ="` + `/new/">add previous</a>`
		previous = title
	default:
		previouslink = `<a href ="` + kind + previous + `">` + previous + `</a>`
	}
	switch {
	case next == "":
		nextlink = `<a href ="` + `/new/">add next</a>`
		next = title
	default:
		nextlink = `<a href ="` + kind + next + `">` + next + `</a>`
	}
	var textrefrences []string
	for i := range transcription.TextRef {
		if transcription.TextRef[i] == "imgCollection" || transcription.TextRef[i] == "meta" {
			continue
		}
		requestedbucket := transcription.TextRef[i]
		texturn := requestedbucket + strings.Split(title, ":")[4]

		// adding testing if requestedbucket exists...
		retrieveddata, _ := h.Repository.SelectCollectionBucketKeyValue(transcription.ColID, requestedbucket, texturn, userid)
		if retrieveddata.JSON == "" {
			continue
		}

		retrievedjson := gocite.Passage{}
		json.Unmarshal([]byte(retrieveddata.JSON), &retrievedjson)

		ctsurn := retrievedjson.PassageID
		var htmllink string
		switch {
		case ctsurn == title:
			htmllink = `<option value="` + kind + ctsurn + `" selected>` + transcription.TextRef[i] + `</option>`
		case ctsurn == "":
			//ctsurn, _ = BoltRetrieveFirstKey(dbname, requestedbucket)
			ctsurn = "here was BoltRetriveFirstKey"
			htmllink = `<option value="` + kind + ctsurn + `">` + transcription.TextRef[i] + `</option>`
		default:
			htmllink = `<option value="` + kind + ctsurn + `">` + transcription.TextRef[i] + `</option>`
		}
		textrefrences = append(textrefrences, htmllink)
	}
	textref := strings.Join(textrefrences, " ")
	imageref := strings.Join(transcription.ImageRef, "#")
	beginjs := `<script type="text/javascript">
	window.onload = function() {`
	startjs := `
		var a`
	start2js := `= document.getElementById("imageLink`
	middlejs := `");
	a`
	middle2js := `.onclick = function() {
		imgUrn="`
	endjs := `"
	reloadImage();
	return false;
}`
	finaljs := `
}
</script>`
	starthtml := `<a id="imageLink`
	middlehtml := `">`
	endhtml := ` </a>`
	var jsstrings, htmlstrings []string
	jsstrings = append(jsstrings, beginjs)
	for i := range transcription.ImageRef {
		jsstring := startjs + strconv.Itoa(i) + start2js + strconv.Itoa(i) + middlejs + strconv.Itoa(i) + middle2js + transcription.ImageRef[i] + endjs
		jsstrings = append(jsstrings, jsstring)
		htmlstring := starthtml + strconv.Itoa(i) + middlehtml + transcription.ImageRef[i] + endhtml
		htmlstrings = append(htmlstrings, htmlstring)
	}
	jsstrings = append(jsstrings, finaljs)
	jsstring := strings.Join(jsstrings, "")
	htmlstring := strings.Join(htmlstrings, "")
	imagescript := template.HTML(jsstring)
	imagehtml := template.HTML(htmlstring)
	texthtml := template.HTML(textref)
	previoushtml := template.HTML(previouslink)
	nexthtml := template.HTML(nextlink)
	return models.TranscriptionPage{
		ColID:        colid,
		User:         user,
		Title:        title,
		Text:         template.HTML(text),
		Previous:     previous,
		PreviousLink: previoushtml,
		Next:         next,
		NextLink:     nexthtml,
		First:        first,
		Last:         last,
		ImageScript:  imagescript,
		ImageHTML:    imagehtml,
		TextHTML:     texthtml,
		ImageRef:     imageref,
		CatID:        catid,
		CatCit:       catcit,
		CatGroup:     catgroup,
		CatWork:      catwork,
		CatVers:      catversion,
		CatExmpl:     catexpl,
		CatOn:        caton,
		CatLan:       catlan,
		ImageJS:      imagejs}, nil
}
