// Package api provides a set of HTTP APIs for handling various server responses.
package api

import (
	"brucheion/models"
	"strings"

	//"github.com/ThomasK81/gocite"
	"github.com/vedicsociety/platform/authorization/identity"
	"github.com/vedicsociety/platform/http/actionresults"
	"github.com/vedicsociety/platform/logging"
)


// ImagesHandler is a struct that encapsulates methods for handling operations related to images.
type ImagesHandler struct {
	Repository models.Repository
	identity.User
	logging.Logger
}

// GetLocalImages is a method that retrieves a list of local images.
// It generates a list of URNs for all images in the local image collection list.
// The URNs are transformed into a specific format and returned in a JSON response.
func (h ImagesHandler) GetLocalImages() actionresults.ActionResult {

	//userid := h.User.GetID()

	dict := h.Repository.LoadImageCollectionList()
	urndict := make([]string, 0)
	for _, f := range dict {
		modf := "urn:cite2:" + strings.TrimSuffix(strings.ReplaceAll(f, "/", ":"), ".dzi")
		parts := strings.Split(modf, ":")
		parts[3] = parts[3] + "." + parts[4] + ":" + parts[5]// Concatenate parts[3] and parts[4] with a dot
		parts = parts[:4]                    // Remove the now redundant parts[4]
		modf = strings.Join(parts, ":")
		urndict = append(urndict, modf)
	}

	h.Logger.Debugf("api.GetCollectionImages: select user bucketdict: (len(dist), urn) ", urndict)

	resp := models.JSONResponse{
		Status:  "success",
		Message: "",
		Data:    urndict,
	}
	return actionresults.NewJsonAction(resp)
}
