package newauth

import (
	"brucheion/auth"
	"brucheion/models"
	"brucheion/utils"
	"fmt"
	"strings"
	"unicode"

	"github.com/vedicsociety/platform/authorization/identity"
	"github.com/vedicsociety/platform/config"
	"github.com/vedicsociety/platform/http/actionresults"
	"github.com/vedicsociety/platform/http/handling"
	"github.com/vedicsociety/platform/logging"
	"github.com/vedicsociety/platform/sessions"
	"golang.org/x/crypto/bcrypt"
)

type NewAuthenticationHandler struct {
	identity.User
	identity.SignInManager
	identity.UserStore
	sessions.Session
	handling.URLGenerator
	logging.Logger
	config.Configuration
	models.Repository
}

type Links struct {
	SignUpUrl      string
	SignInUrl      string
	ForgotPwdUrl   string
	GoogleUrl      string
	SignUpAllowed  bool
	ShowGoogleAuth bool
	ShowGithubAuth bool
	UseRecaptcha   bool
	RecaptchaKey   string
}

func (handler NewAuthenticationHandler) getLinks() Links {
	signUpUrl := utils.MustGenerateUrl(handler.URLGenerator, NewAuthenticationHandler.GetSignUp)
	signInUrl := utils.MustGenerateUrl(handler.URLGenerator, NewAuthenticationHandler.GetSignIn)
	forgotPwdUrl := utils.MustGenerateUrl(handler.URLGenerator, NewAuthenticationHandler.GetForgotPwd)
	googleUrl := "../oauth/google"
	allowSignup := handler.Configuration.GetBoolDefault("authorization:allowsignup", false)
	showGoogle := handler.Configuration.GetBoolDefault("authorization:showgoogle", false)
	showGithub := handler.Configuration.GetBoolDefault("authorization:showgithub", false)
	useRecaptcha := handler.Configuration.GetBoolDefault("authorization:userecaptcha", false)
	recaptchaKey, _ := handler.Configuration.GetString("authorization:recaptchakey")

	return Links{
		SignUpUrl:      signUpUrl,
		SignInUrl:      signInUrl,
		ForgotPwdUrl:   forgotPwdUrl,
		GoogleUrl:      googleUrl,
		SignUpAllowed:  allowSignup,
		ShowGoogleAuth: showGoogle,
		ShowGithubAuth: showGithub,
		UseRecaptcha:   useRecaptcha,
		RecaptchaKey:   recaptchaKey,
	}
}

// region Signin
func (handler NewAuthenticationHandler) GetSignIn() actionresults.ActionResult {
	return actionresults.NewTemplateAction("newauth_signin.html",
		handler.getLinks())
}

type SigninRequest struct {
	Email             string `json:"email"`
	Password          string `json:"password"`
	RecaptchaResponse string `json:"g-recaptcha-response"`
}

type SigninContext struct {
	SigninRequest
	Links
	EmailErrorMsg    string
	PasswordErrorMsg string
}

func (handler NewAuthenticationHandler) PostSignIn(params SigninRequest) actionresults.ActionResult {
	template := "widget_signin.html"
	context := SigninContext{
		params,
		handler.getLinks(),
		"",
		"",
	}
	if params.Email == "" {
		return actionresults.NewTemplateAction(template, context)
	}

	useRecaptcha := handler.Configuration.GetBoolDefault("authorization:userecaptcha", false)
	if useRecaptcha {
		recaptchaSecret, _ := handler.Configuration.GetString("authorization:recaptchasecret")
		if err := auth.CheckRecaptcha(recaptchaSecret, params.RecaptchaResponse, "signin"); err != nil {
			context.EmailErrorMsg = "Recaptcha Failed"
			return actionresults.NewTemplateAction(template, context)
		}
	}

	creds, err := handler.Repository.GetUserByEmail(params.Email)
	if err != nil {
		panic(err)
	}

	if creds.Id == 0 || !creds.IsVerified {
		context.EmailErrorMsg = "User doesn't exists"
		return actionresults.NewTemplateAction(template, context)
	}

	if !IsPasswordOK(creds, params.Password) {
		context.PasswordErrorMsg = "Wrong credentials"
		return actionresults.NewTemplateAction(template, context)
	}
	err = handler.SignInManager.SignIn(creds)
	if err != nil {
		panic(err)
	}

	rcontext := RedirectContext{
		Msg:       "Sign in successful, you will now be redirected to the home page.",
		Url:       "/",
		AfterTime: 2000,
	}

	return actionresults.NewTemplateAction("widget_redirect.html", rcontext)
}

func IsPasswordOK(u models.Credentials, p string) (f bool) {
	// Compare the stored hashed password, with the hashed version of the password that was received
	err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(p))
	if err != nil {
		f = false
	} else {
		f = true
	}
	return
}

//endregion

//region ForgotPwd

func (handler NewAuthenticationHandler) GetForgotPwd() actionresults.ActionResult {
	return actionresults.NewTemplateAction("newauth_forgotpwd.html",
		handler.getLinks())
}

type ForgotPwdRequest struct {
	Email             string `json:"email"`
	RecaptchaResponse string `json:"g-recaptcha-response"`
}

type ForgotPwdContext struct {
	ForgotPwdRequest
	Links
	ErrorMsg   string
	SuccessMsg string
}

func (handler NewAuthenticationHandler) PostForgotPwd(params ForgotPwdRequest) actionresults.ActionResult {
	template := "widget_forgotpwd.html"
	context := ForgotPwdContext{
		params,
		handler.getLinks(),
		"",
		"",
	}
	if params.Email == "" {
		return actionresults.NewTemplateAction(template, context)
	}
	useRecaptcha := handler.Configuration.GetBoolDefault("authorization:userecaptcha", false)
	if useRecaptcha {
		recaptchaSecret, _ := handler.Configuration.GetString("authorization:recaptchasecret")
		if err := auth.CheckRecaptcha(recaptchaSecret, params.RecaptchaResponse, "fpwd"); err != nil {
			context.ErrorMsg = "Recaptcha Failed"
			return actionresults.NewTemplateAction(template, context)
		}
	}

	creds, err := handler.Repository.GetUserByEmail(params.Email)
	if err != nil {
		panic(err)
	}
	if creds.Id == 0 {
		context.ErrorMsg = "User doesn't exits"
		return actionresults.NewTemplateAction(template, context)
	}
	code := utils.RandomString(64)
	creds.VerificationCode = code
	err = handler.Repository.UpdateUser(&creds)
	if err != nil {
		panic(err)
	}
	hostname, _ := handler.Configuration.GetString("system:hostname")
	emailData := utils.EmailContext{
		URL:          hostname + utils.MustGenerateUrl(handler.URLGenerator, NewAuthenticationHandler.GetVerifiedSignUp, "reset", code),
		FirstName:    creds.Username,
		Subject:      "Brucheion Password Reset",
		Msg:          "We have received a request to reset your password, please click the link below to reset your password:",
		CallToAction: "Reset Password",
	}
	err = utils.SendEmail(handler.Configuration, handler.Logger, &creds, &emailData)
	if err != nil {
		panic(err)
	} else {
		context.SuccessMsg = "Kindly check your email for your reset link."
		return actionresults.NewTemplateAction(template, context)
	}

}

//endregion

//region SignUp

func (handler NewAuthenticationHandler) GetSignUp() actionresults.ActionResult {
	allowSignup := handler.Configuration.GetBoolDefault("authorization:allowsignup", false)
	if !allowSignup {
		rcontext := MessageContext{
			Message: "Signups not allowed",
		}
		return actionresults.NewTemplateAction("message.html", rcontext)
	}

	return actionresults.NewTemplateAction("newauth_signup.html",
		handler.getLinks())
}

type SignUpRequest struct {
	Email             string `json:"email"`
	RecaptchaResponse string `json:"g-recaptcha-response"`
}

type SignUpContext struct {
	SignUpRequest
	Links
	ErrorMsg   string
	SuccessMsg string
}

func (handler NewAuthenticationHandler) PostSignUp(params SignUpRequest) actionresults.ActionResult {
	allowSignup := handler.Configuration.GetBoolDefault("authorization:allowsignup", false)
	if !allowSignup {
		rcontext := MessageContext{
			Message: "Signups not allowed",
		}
		return actionresults.NewTemplateAction("message.html", rcontext)
	}

	template := "widget_signup.html"
	context := SignUpContext{
		params,
		handler.getLinks(),
		"",
		"",
	}
	if params.Email == "" {
		return actionresults.NewTemplateAction("widget_signup.html", context)
	}
	useRecaptcha := handler.Configuration.GetBoolDefault("authorization:userecaptcha", false)
	if useRecaptcha {
		recaptchaSecret, _ := handler.Configuration.GetString("authorization:recaptchasecret")
		if err := auth.CheckRecaptcha(recaptchaSecret, params.RecaptchaResponse, "signup"); err != nil {
			context.ErrorMsg = "Recaptcha Failed"
			return actionresults.NewTemplateAction(template, context)
		}
	}

	creds, err := handler.Repository.GetUserByEmail(params.Email)
	if err != nil {
		panic(err)
	}

	if creds.Id > 0 && creds.IsVerified {
		context.ErrorMsg = "Account already exists"
		return actionresults.NewTemplateAction(template, context)
	}
	code := utils.RandomString(64)
	if creds.Id > 0 && !creds.IsVerified {
		code = creds.VerificationCode
	}
	if creds.Id == 0 {
		creds, err = handler.Repository.AddNewUser(models.Credentials{
			Email:            params.Email,
			Username:         strings.Split(params.Email, "@")[0],
			Password:         "NONE",
			IsVerified:       false,
			VerificationCode: code,
		})
		if err != nil {
			panic(err)
		}
	}
	hostname, _ := handler.Configuration.GetString("system:hostname")
	emailData := utils.EmailContext{
		URL:          hostname + utils.MustGenerateUrl(handler.URLGenerator, NewAuthenticationHandler.GetVerifiedSignUp, "signup", code),
		FirstName:    params.Email,
		Subject:      "Your account verification code",
		Msg:          "Thanks for signing up. Click the link below to confirm your registration and you'll be on your way.",
		CallToAction: "Confirm your registration",
	}
	err = utils.SendEmail(handler.Configuration, handler.Logger, &creds, &emailData)
	if err != nil {
		panic(err)
	} else {
		context.SuccessMsg = "Verification email sent"
	}
	return actionresults.NewTemplateAction("widget_signup.html", context)
}

//endregion

// region VerifiedSignup

func (handler NewAuthenticationHandler) GetVerifiedSignUp(purpose string, code string) actionresults.ActionResult {
	context := VerifiedSignUpContext{}
	context.Code = code
	context.Purpose = purpose
	if purpose == "signup" {
		context.CallToAction = "Sign up"
		context.ForSignup = true
	} else if purpose == "reset" {
		context.CallToAction = "Reset Password"
		context.ForReset = true
	} else {
		context.CallToAction = purpose
	}

	if context.ForSignup {
		allowSignup := handler.Configuration.GetBoolDefault("authorization:allowsignup", false)
		if !allowSignup {
			rcontext := MessageContext{
				Message: "Signups not allowed",
			}
			return actionresults.NewTemplateAction("message.html", rcontext)
		}
	}

	creds, err := handler.Repository.GetUserByVerificationCode(code)
	if err != nil {
		panic(err)
	}
	if creds.Id == 0 {
		return actionresults.NewErrorAction(fmt.Errorf("Invalid Verification code"))
	}

	context.Email = creds.Email
	return actionresults.NewTemplateAction("newauth_verifiedsignup.html", context)
}

type VerifiedSignUpRequest struct {
	Name      string `json:"name"`
	Code      string `json:"code"`
	Password  string `json:"password"`
	Password2 string `json:"password2"`
	Purpose   string `json:"purpose"`
}

type VerifiedSignUpContext struct {
	VerifiedSignUpRequest
	Links
	Email            string
	PasswordErrorMsg string
	SuccessMsg       string
	ForReset         bool
	ForSignup        bool
	CallToAction     string
}

type RedirectContext struct {
	Msg       string
	Url       string
	AfterTime int
}

type MessageContext struct {
	Message string
}

func (handler NewAuthenticationHandler) PostVerifiedSignUp(params VerifiedSignUpRequest) actionresults.ActionResult {
	if params.Code == "" {
		return actionresults.NewErrorAction(fmt.Errorf("Invalid Verification code"))
	}
	creds, err := handler.Repository.GetUserByVerificationCode(params.Code)
	if err != nil {
		panic(err)
	}
	if creds.Id == 0 {
		return actionresults.NewErrorAction(fmt.Errorf("Invalid Verification code"))
	}
	context := VerifiedSignUpContext{}
	context.Links = handler.getLinks()
	context.VerifiedSignUpRequest = params
	context.Email = creds.Email
	if params.Purpose == "signup" {
		context.CallToAction = "Sign up"
		context.ForSignup = true
	} else if params.Purpose == "reset" {
		context.CallToAction = "Reset Password"
		context.ForReset = true
	} else {
		context.CallToAction = context.Purpose
	}
	if context.Password == "" {
		context.Name = creds.Username
		return actionresults.NewTemplateAction("widget_verifiedsignup.html", context)
	}
	if params.Password != params.Password2 {
		context.PasswordErrorMsg = "Passwords do not match."
		return actionresults.NewTemplateAction("widget_verifiedsignup.html", context)
	}
	if isPasswordWeak(params.Password) {
		context.PasswordErrorMsg = "Password should be > 8 characters and contain 1 alphabet. 1 number."
		return actionresults.NewTemplateAction("widget_verifiedsignup.html", context)
	}
	creds.Username = params.Name
	creds.VerificationCode = "verified"
	hash, err := bcrypt.GenerateFromPassword([]byte(params.Password), 13)
	if err != nil {
		panic(err)
	}
	creds.Password = string(hash)
	creds.IsVerified = true

	err = handler.Repository.UpdateUser(&creds)
	if err != nil {
		panic(err)
	}
	err = handler.SignInManager.SignIn(creds)
	if err != nil {
		panic(err)
	}

	rcontext := RedirectContext{
		Msg:       "Thank you for registering with Brucheion, you will now be redirected to the home page.",
		Url:       "/",
		AfterTime: 2000,
	}
	return actionresults.NewTemplateAction("widget_redirect.html", rcontext)
}

func isPasswordWeak(password string) bool {
	if len(password) < 8 {
		return true
	}
	hasAplha := false
	hasNum := false
	for _, r := range password {
		if unicode.IsLetter(r) {
			hasAplha = true
		}
		if unicode.IsNumber(r) {
			hasNum = true
		}
	}

	return !(hasAplha && hasNum)
}

// endregion
type Query struct {
	Message string
}

func (q Query) Get(key string) string {
	var value string
	switch key {
	case "message":
		value = q.Message
	default:
		value = ""
	}
	return value
}

func (handler NewAuthenticationHandler) GetMessage(query Query) actionresults.ActionResult {
	rcontext := MessageContext{
		Message: query.Message,
	}
	return actionresults.NewTemplateAction("message.html", rcontext)

}
