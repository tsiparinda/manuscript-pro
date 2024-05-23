package auth

import (
	"brucheion/utils"

	"github.com/vedicsociety/platform/authorization/identity"
	"github.com/vedicsociety/platform/http/actionresults"
	"github.com/vedicsociety/platform/http/handling"
	"github.com/vedicsociety/platform/logging"
)

type SignOutHandler struct {
	identity.User
	handling.URLGenerator
	logging.Logger
}

func (handler SignOutHandler) GetUserWidget() actionresults.ActionResult {
	//handler.Logger.Debugf("GetUserWidget user, url", handler.User, MustGenerateUrl(handler.URLGenerator, AuthenticationHandler.PostSignOut))
	return actionresults.NewTemplateAction("user_widget.html", struct {
		identity.User
		SignoutUrl string
	}{
		handler.User,
		utils.MustGenerateUrl(handler.URLGenerator,
			AuthenticationHandler.PostSignOut),
	})
}
