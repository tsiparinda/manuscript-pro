/*
Adapted from package gothic to use with Platform framework
Package gothic wraps common behaviour when using Goth. This makes it quick, and easy, to get up
and running with Goth. Of course, if you want complete control over how things flow, in regard
to the authentication process, feel free and use Goth directly.
See https://github.com/markbates/goth/blob/master/examples/main.go to see this in action.
*/
package auth

import (
	utils "brucheion/utils"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"net/url"

	"github.com/gorilla/sessions"
	"github.com/markbates/goth"
)

// SessionName is the key used to access the session store.
const SessionName = "_gothic_session"

// Store can/should be set by applications using gothic. The default is a cookie store.
var Store sessions.Store
var defaultStore sessions.Store

var keySet = false

type key int

// ProviderParamKey can be used as a key in context when passing in a provider
const ProviderParamKey key = iota

func init() {
	//key := []byte(os.Getenv("SESSION_SECRET"))
	key := []byte(utils.RandomString(64))
	keySet = len(key) != 0

	cookieStore := sessions.NewCookieStore(key)
	cookieStore.Options.HttpOnly = true
	Store = cookieStore
	defaultStore = Store
}

/*
BeginAuthHandler is a convenience handler for starting the authentication process.
It expects to be able to get the name of the provider from the query parameters
as either "provider" or ":provider".
BeginAuthHandler will redirect the user to the appropriate authentication end-point
for the requested provider.
See https://github.com/markbates/goth/examples/main.go to see this in action.
*/
func (handler AuthenticationHandler) BeginAuthHandler(provider string) string {

	url, err := handler.GetAuthURL(provider)
	handler.Logger.Debugf("BeginAuthHandler url", url)
	if err != nil {
		handler.Logger.Debugf("BeginAuthHandler error GetAuthURL ", err.Error())
		// res.WriteHeader(http.StatusBadRequest)
		// fmt.Fprintln(res, err)
		return "/"
	}
	//actionresults.NewRedirectAction("https://google.com")
	// http.Redirect(res, req, url, http.StatusTemporaryRedirect)
	return url
}

// SetState sets the state string associated with the given request.
// If no state string is associated with the request, one will be generated.
// This state is sent to the provider and can be retrieved during the
// callback.
func (handler AuthenticationHandler) SetState() string {
	state := handler.Session.GetValueDefault(utils.OAUTH_STATE, "").(string)
	if len(state) > 0 {
		return state
	}

	// If a state query param is not passed in, generate a random
	// base64-encoded nonce so that the state on the auth URL
	// is unguessable, preventing CSRF attacks, as described in
	//
	// https://auth0.com/docs/protocols/oauth2/oauth-state#keep-reading
	nonceBytes := make([]byte, 64)
	_, err := io.ReadFull(rand.Reader, nonceBytes)
	if err != nil {
		panic("gothic: source of randomness unavailable: " + err.Error())
	}
	return base64.URLEncoding.EncodeToString(nonceBytes)
}

// GetState gets the state returned by the provider during the callback.
// This is used to prevent CSRF attacks, see
// http://tools.ietf.org/html/rfc6749#section-10.12
func (handler AuthenticationHandler) GetState() string {
	params := handler.Session.GetValueDefault(utils.OAUTH_STATE, "").(string)
	return params
}

/*
GetAuthURL starts the authentication process with the requested provided.
It will return a URL that should be used to send users to.
It expects to be able to get the name of the provider from the query parameters
as either "provider" or ":provider".
I would recommend using the BeginAuthHandler instead of doing all of these steps
yourself, but that's entirely up to you.
*/
func (handler AuthenticationHandler) GetAuthURL(providerName string) (string, error) {
	if !keySet && defaultStore == Store {
		fmt.Println("GetAuthURL goth/gothic: no SESSION_SECRET environment variable is set. The default cookie store is not available and any calls will fail. Ignore this warning if you are using a different store.")
	}

	provider, err := goth.GetProvider(providerName)
	if err != nil {
		return "", err
	}
	handler.Logger.Debugf("oauth.GetAuthURL Provider: ", provider)
	sess, err := provider.BeginAuth(handler.SetState())
	if err != nil {
		return "", err
	}

	url, err := sess.GetAuthURL()
	if err != nil {
		return "", err
	}

	handler.Session.SetValue(utils.SIGNIN_PROVIDER, providerName)
	handler.Session.SetValue(providerName, sess.Marshal())
	handler.Logger.Debugf("oauth.GetAuthURL providerName, sess.Marshal: ", providerName, sess.Marshal())
	return url, err
}

// validateState ensures that the state token param from the original
// AuthURL matches the one included in the current (callback) request.
func (handler AuthenticationHandler) validateState(sess goth.Session) error {
	rawAuthURL, err := sess.GetAuthURL()
	if err != nil {
		handler.Logger.Debugf("validateState rawAuthURL err", err.Error())
		return err
	}
	handler.Logger.Debugf("validateState rawAuthURL", rawAuthURL)

	authURL, err := url.Parse(rawAuthURL)
	if err != nil {
		handler.Logger.Debugf("validateState authURL err", err.Error())
		return err
	}
	handler.Logger.Debugf("validateState authURL", authURL)

	reqState := handler.GetState()
	handler.Logger.Debugf("validateState reqState", reqState)

	originalState := authURL.Query().Get("state")
	handler.Logger.Debugf("validateState originalState", originalState)
	if originalState != "" && (originalState != reqState) {
		return errors.New("state token mismatch")
	}
	return nil
}

/*
CompleteUserAuth does what it says on the tin. It completes the authentication
process and fetches all the basic information about the user from the provider.
It expects to be able to get the name of the provider from the query parameters
as either "provider" or ":provider".
See https://github.com/markbates/goth/examples/main.go to see this in action.
*/
func (handler AuthenticationHandler) CompleteUserAuth(providerName string) (goth.User, error) {
	if !keySet && defaultStore == Store {
		handler.Logger.Debugf("goth/gothic: no SESSION_SECRET environment variable is set. The default cookie store is not available and any calls will fail. Ignore this warning if you are using a different store.")
	}
	handler.Logger.Debugf("CompleteUserAuth providerName", providerName)
	provider, err := goth.GetProvider(providerName)
	if err != nil {
		return goth.User{}, err
	}

	handler.Logger.Debugf("CompleteUserAuth provider", provider)

	value := handler.Session.GetValueDefault(providerName, "").(string)
	handler.Logger.Debugf("CompleteUserAuth value_intf", value)

	if value == "" {
		err := errors.New("CompleteUserAuth could not find matching session for this request")
		handler.Logger.Debugf("CompleteUserAuth err", err.Error())
		return goth.User{}, err
	}

	handler.Logger.Debugf("CompleteUserAuth value", value)
	//	handler.Logger.Debugf("value", value)

	//	defer handler.Logout()
	sess, err := provider.UnmarshalSession(value)
	if err != nil {
		return goth.User{}, err
	}
	handler.Logger.Debugf("CompleteUserAuth sess: ", sess)

	err = handler.validateState(sess)
	if err != nil {
		handler.Logger.Debugf("CompleteUserAuth validateState: ", err.Error())
		return goth.User{}, err
	}
	handler.Logger.Debugf("CompleteUserAuth validateState passed: ")

	user, err := provider.FetchUser(sess)
	handler.Logger.Debugf("CompleteUserAuth user %v\n err %v", user, err.Error())
	if err == nil {
		// user can be found with existing session data
		return user, err
	}

	code := handler.Session.GetValueDefault(utils.OAUTH_CODE, "").(string)
	handler.Logger.Debugf("CompleteUserAuth code_intf", code)

	if code == "" {
		err := errors.New("CompleteUserAuth could not find code for this session")
		handler.Logger.Debugf("CompleteUserAuth err", err.Error())
		return goth.User{}, err
	}

	var query Query = Query{State: "",
		Code: code,
	}

	// query explanation see in comments
	handler.Logger.Debugf("CompleteUserAuth query ", query)
	// get new token and retry fetch
	_, err = sess.Authorize(provider, query)
	if err != nil {
		handler.Logger.Debugf("Authorize", err.Error())
		return goth.User{}, err
	}

	handler.Session.SetValue(providerName, sess.Marshal())

	gu, err := provider.FetchUser(sess)
	return gu, err
}
