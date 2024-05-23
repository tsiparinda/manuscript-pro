package auth

import (
	"brucheion/models"
	"brucheion/utils"

	"github.com/markbates/goth"
	"golang.org/x/crypto/bcrypt"
)

// following oauth authentication process, we arise this procedure after get user's credentials from oauth provider
func (handler AuthenticationHandler) SignInGothUser(user *goth.User) (url string) {
	if user.Email == "" {
		handler.Logger.Debugf("SignInGothUser no email address arrived from provider: %v", user)
		handler.Session.SetValue(utils.SIGNIN_MSG_KEY, "Email is not arrived from provider")
		return "/"
	}
	creds, err := handler.Repository.GetUserByEmail(user.Email)
	if err != nil {
		handler.Logger.Panicf("SignInGothUser An error arise by exec GetUserByEmail command: %v", err.Error())
		return "/"
	}

	// Generate Verification Code
	code := utils.RandomString(64)
	verification_code := code
	if creds.Id == 0 { // user not found
		allowSignup := handler.Configuration.GetBoolDefault("authorization:allowsignup", false)
		if !allowSignup {
			url = utils.GenerateMsgUrl(handler.URLGenerator, "Signup not allowed")
			return url
		}

		hash, err := bcrypt.GenerateFromPassword([]byte(utils.RandomString(64)), 13)

		creds, err = handler.AddNewUser(models.Credentials{
			Username:         user.NickName,
			Email:            user.Email,
			VerificationCode: verification_code,
			Password:         string(hash),
		})
		if err != nil {
			handler.Logger.Panicf("SignInGothUser An error arise by exec AddNewUser command: %v", err.Error())
			return "/"
		}
	}

	// add creds user roles
	if err = handler.Repository.GetUserGroups(&creds); err != nil {
		handler.Logger.Panicf("SignInGothUser An error arise by exec GetUserGroups command: %v", err.Error())
		return "/"
	}

	handler.Logger.Debugf("SignInGothUser creds.IsVerified", creds.IsVerified)
	if !(creds.IsVerified) {
		// go to verification
		url = utils.GenerateVerifiedSignup(handler.URLGenerator, verification_code)
	} else {
		handler.Session.SetValue(utils.SIGNIN_MSG_KEY, "")
		handler.Logger.Debugf("SignInGothUser creds", creds)
		handler.SignInManager.SignIn(&creds)
		url = "/"
	}
	handler.Logger.Debugf("SignInGothUser url", url)

	return url
}

func (handler AuthenticationHandler) SignInCredsUser(user models.Credentials) (url string) {

	if utils.IsNil(&user.Roles) {
		if err := handler.Repository.GetUserGroups(&user); err != nil {
			handler.Logger.Debugf("SignInCredUser NO ROLES", user)
			handler.Session.SetValue(utils.SIGNIN_MSG_KEY, err.Error())
			return "/"
		}
	}

	// user.Roles = append(user.Roles, "ToolsUsers")

	handler.Logger.Debugf("SignInCredUser user", user)
	handler.SignInManager.SignIn(user)
	url = "/"
	handler.Logger.Debugf("SignInCredUser url", url)
	handler.Session.SetValue(utils.SIGNIN_MSG_KEY, "")
	return url
}
