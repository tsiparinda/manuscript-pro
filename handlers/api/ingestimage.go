// Package api provides a set of HTTP APIs for handling various server responses.
package api

import (
	"brucheion/models"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/vedicsociety/platform/authorization/identity"
	"github.com/vedicsociety/platform/http/actionresults"
	"github.com/vedicsociety/platform/http/handling"
	"github.com/vedicsociety/platform/logging"
)

// IngestImageHandler is a struct that contains repositories, user identities, 
// loggers and URL generators for managing the ingestion of images.
type IngestImageHandler struct {
	Repository models.Repository
	identity.User
	logging.Logger
	handling.URLGenerator
}

// JSONlist is a struct that contains a slice of strings. It's used for handling JSON requests.
type JSONlist struct {
	Item []string `json:"item"`
}

// GetImageInfo is a function that retrieves the metadata of a specific image in the user database.
// It fetches the image metadata based on provided collection ID, collection Images, and image URN.
// It returns an ActionResult containing the image metadata or an error if one occurred.
// from old Brucheion, image.go, getImageInfo
// prints the metadata of a specific image in the user database.
func (h IngestImageHandler) GetImageInfo(colid int, colImages string, imageURN string) actionresults.ActionResult {

	userid := h.User.GetID()

	newImage := models.Image{}

	val, err := h.Repository.LoadCollectionImageKeyValue(colid, colImages, userid)
	if err != nil {
		h.Logger.Debugf("api.GetCollectionImages: : Internal server error1 %v", err.Error())
		return actionresults.NewErrorAction(fmt.Errorf("api.GetCollectionImages: failed to get bucket"))
	}
	h.Logger.Debugf("api.GetCollectionImages: : ", val)
	collection := models.ImageCollection{}

	json.Unmarshal([]byte(val.JSON), &collection.Collection)

	for _, v := range collection.Collection {
		if v.URN == imageURN {
			newImage = v
		}
	}
	h.Logger.Debugf("api.GetCollectionImages: newImage ", newImage)
	resp := models.JSONResponse{
		Status:  "success",
		Message: "",
		Data:    newImage,
	}
	return actionresults.NewJsonAction(resp)
}

// PostAddImageToCITE is a function that adds image metadata to the specified collection in the bucket imgCollection in a user database.
// It extracts the reference from the the http.Request and passes it to addtoCITECollection.
// It returns an ActionResult containing the HTTP status code indicating success or failure of the operation.
// from old Brucheion, CITEcollection.go, addCITE and addImageToCITECollection
// adds image metadata to the specified collection in the bucket imgCollection in a user database
// It extracts the reference from the the http.Request and passes it to addtoCITECollection
// Examples:
// localhost:7000/addtoCITE?name="urn:cite2:iiifimages:test:"&urn="urn:cite2:iiifimages:test:1"&external="true"&protocol="iiif"&location="https://libimages1.princeton.edu/loris/pudl0001%2F4609321%2Fs42%2F00000004.jp2/info.json"
// localhost:7000/addtoCITE?name="urn:cite2:staticimages:test:"&urn="urn:cite2:staticimages:test:1"&external="true"&protocol="static"&location="https://upload.wikimedia.org/wikipedia/commons/8/81/Rembrandt_The_Three_Crosses_1653.jpg"
// localhost:7000/addtoCITE?name="urn:cite2:dzi:test:"&urn="urn:cite2:nyaya:M3img.positive:m3_098"&external="false"&protocol="localDZ"&location="urn:cite2:nyaya:M3img.positive:m3_098"
// localhost:7000/addtoCITE?name="urn:cite2:mixedimages:test:"&urn="urn:cite2:iiifimages:test:1"&external="true"&protocol="iiif"&location="https://libimages1.princeton.edu/loris/pudl0001%2F4609321%2Fs42%2F00000004.jp2/info.json"
// localhost:7000/addtoCITE?name="urn:cite2:mixedimages:test:"&urn="urn:cite2:staticimages:test:1"&external="true"&protocol="static"&location="https://upload.wikimedia.org/wikipedia/commons/8/81/Rembrandt_The_Three_Crosses_1653.jpg"
// localhost:7000/addtoCITE?name="urn:cite2:mixedimages:test:"&urn="urn:cite2:nyaya:M3img.positive:m3_098"&external="false"&protocol="localDZ"&location="urn:cite2:nyaya:M3img.positive:m3_098"
func (h IngestImageHandler) PostAddImageToCITE(image models.Image) actionresults.ActionResult {
	h.Logger.Debugf("PostAddImageToCITE input: ", image)
	userid := h.User.GetID()

	err := h.Repository.SaveImageData(userid, &image)

	if err != nil {
		h.Logger.Debugf("PostAddImageToCITE: Cannot save metadata", err.Error())
		return &actionresults.ErrorActionResult{}
	}

	return &StatusCodeResult{http.StatusOK}
}
