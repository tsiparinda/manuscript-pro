package auth

import (
	"brucheion/utils"

	"github.com/markbates/goth"

	//"github.com/markbates/goth/gothic"

	"github.com/vedicsociety/platform/authorization/identity"
	"github.com/vedicsociety/platform/http/actionresults"
)

// type = immitator this code
// params := req.URL.Query()
//
//	if params.Encode() == "" && req.Method == "POST" {
//		req.ParseForm()
//		params = req.Form
//	}
type Query struct {
	State string
	Code  string
}

func (q Query) Get(key string) string {
	var value string
	switch key {
	case "state":
		value = q.State
	case "code":
		value = q.Code
	default:
		value = ""
	}
	return value
}

func (handler AuthenticationHandler) GetOauth(provider string) actionresults.ActionResult {

	var oauth_state, url string
	oauth_state = utils.RandomString(64)
	handler.Session.SetValue(utils.OAUTH_STATE, oauth_state)
	handler.Logger.Debugf("GetOauth provider ", provider)
	// try to get the user without re-authenticating
	if gothUser, err := handler.CompleteUserAuth(provider); err == nil {
		handler.Logger.Debugf("GetOauth go away by CompleteUserAuth, gothUser", gothUser)
		handler.SignInManager.SignIn(identity.NewBasicUser(2, gothUser.Email, "ToolsUser"))
		return actionresults.NewRedirectAction("/")
	} else {
		handler.Logger.Debugf("GetOauth from CompleteUserAuth, err", err.Error())
		// add provider to header in the req
		url = handler.BeginAuthHandler(provider)
		handler.Logger.Debugf("GetOauth provider, url ", provider, url)
		return actionresults.NewRedirectAction(url)
	}
}

func (handler AuthenticationHandler) GetOauthCallback(query Query) actionresults.ActionResult {
	handler.Logger.Debugf("GetOauthCallback query: ", query)
	handler.Session.SetValue(utils.OAUTH_CODE, query.Code)
	if !keySet && defaultStore == Store {
		handler.Logger.Debugf("GetOauthCallback goth/gothic: no SESSION_SECRET environment variable is set. The default cookie store is not available and any calls will fail. Ignore this warning if you are using a different store.")
		return actionresults.NewRedirectAction(utils.MustGenerateUrl(handler.URLGenerator,
			AuthenticationHandler.GetSignIn))
	}

	providerName := handler.Session.GetValueDefault(utils.SIGNIN_PROVIDER, "").(string)
	if providerName == "" {
		handler.Logger.Debugf("GetOauthCallback providername not found")
		return actionresults.NewRedirectAction(utils.MustGenerateUrl(handler.URLGenerator,
			AuthenticationHandler.GetSignIn))
	}
	handler.Logger.Debugf("GetOauthCallback providername: ", providerName)

	provider, err := goth.GetProvider(providerName)
	handler.Logger.Debugf("GetOauthCallback provider ", provider)
	if err != nil {
		handler.Logger.Debugf("GetOauthCallback getprovider error", err.Error())
		return actionresults.NewRedirectAction(utils.MustGenerateUrl(handler.URLGenerator,
			AuthenticationHandler.GetSignIn))
	}

	value := handler.Session.GetValueDefault(providerName, "").(string)
	handler.Logger.Debugf("GetOauthCallback value :", value)
	// defer handler.Session.Logout()
	//defer handler.Logout()

	sess, err := provider.UnmarshalSession(value)
	if err != nil {
		handler.Logger.Debugf("GetOauthCallback UnmarshalSession error", err.Error())
		return actionresults.NewRedirectAction(utils.MustGenerateUrl(handler.URLGenerator,
			AuthenticationHandler.GetSignIn))
	}
	handler.Logger.Debugf("GetOauthCallback sess ", sess)

	err = handler.validateState(sess)
	if err != nil {
		handler.Logger.Debugf("GetOauthCallback validateState error", err.Error())
		return actionresults.NewRedirectAction(utils.MustGenerateUrl(handler.URLGenerator,
			AuthenticationHandler.GetSignIn))
	}

	user, err := provider.FetchUser(sess)
	handler.Logger.Debugf("GetOauthCallback user %v\n err %v", user, err.Error())
	if err == nil {
		url := handler.SignInGothUser(&user)
		return actionresults.NewRedirectAction(url)
	}

	// query explanation see in comments
	handler.Logger.Debugf("GetOauthCallback query ", query)
	// get new token and retry fetch
	_, err = sess.Authorize(provider, query)
	if err != nil {
		handler.Logger.Debugf("GetOauthCallback Authorize", err.Error())
		return actionresults.NewRedirectAction(utils.MustGenerateUrl(handler.URLGenerator,
			AuthenticationHandler.GetSignIn))
	}

	handler.Session.SetValue(providerName, sess.Marshal())
	handler.Logger.Debugf("GetOauthCallback sess", sess)

	gu, err := provider.FetchUser(sess)
	if err == nil {
		handler.Logger.Debugf("GetOauthCallback qu:", gu)
		url := handler.SignInGothUser(&gu)
		return actionresults.NewRedirectAction(url)
	}
	handler.Logger.Debugf("GetOauthCallback qu:", gu, err.Error())
	return actionresults.NewRedirectAction(utils.MustGenerateUrl(handler.URLGenerator,
		AuthenticationHandler.GetSignIn))
}
