package auth

import (
	"brucheion/models"
	"brucheion/utils"

	"github.com/vedicsociety/platform/authorization/identity"
	"github.com/vedicsociety/platform/config"
	"github.com/vedicsociety/platform/http/actionresults"
	"github.com/vedicsociety/platform/http/handling"
	"github.com/vedicsociety/platform/logging"
	sessions "github.com/vedicsociety/platform/sessions"
)

type AccountHandler struct {
	identity.User
	handling.URLGenerator
	logging.Logger
	sessions.Session
	models.Repository
	config.Configuration
}

type AccountTemplateContext struct {
	Id          int
	Username    string
	Email       string
	SaveUrl     string
	ResetPWDUrl string
}

func (handler AccountHandler) GetAccountProfile() actionresults.ActionResult {
	//creds := handler.Repository.GetUserByID(handler.Session.GetValue(SIGNIN_USER_ID).(int))
	creds := handler.Repository.GetUserByID(handler.User.GetID())
	return actionresults.NewTemplateAction("account.html", AccountTemplateContext{
		Id:       creds.Id,
		Username: creds.Username,
		Email:    creds.Email,
		SaveUrl: utils.MustGenerateUrl(handler.URLGenerator,
			AccountHandler.PostAccountSave),
		ResetPWDUrl: utils.MustGenerateUrl(handler.URLGenerator, AccountHandler.PostAccountResetPwd),
	})
}

func (handler AccountHandler) PostAccountSave(p AccountTemplateContext) actionresults.ActionResult {
	user := handler.Repository.GetUserByID(p.Id)

	if user.Id == 0 {
		handler.Logger.Panicf("PostAccountSave User not found: ", p)
		return actionresults.NewRedirectAction("/")
	}
	// user.Email = p.Email
	user.Username = p.Username
	handler.Repository.UpdateUser(&user)
	return actionresults.NewRedirectAction("/")

}

func (handler AccountHandler) PostAccountResetPwd(p AccountTemplateContext) actionresults.ActionResult {
	return actionresults.NewRedirectAction(utils.MustGenerateUrl(handler.URLGenerator, AuthenticationHandler.GetForgotPwd))
}
