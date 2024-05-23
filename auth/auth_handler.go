package auth

import (
	"brucheion/models"
	"brucheion/utils"
	"errors"
	"strings"

	"net/mail"

	"github.com/vedicsociety/platform/authorization/identity"
	"github.com/vedicsociety/platform/config"
	"github.com/vedicsociety/platform/http/actionresults"
	"github.com/vedicsociety/platform/http/handling"
	"github.com/vedicsociety/platform/logging"
	"github.com/vedicsociety/platform/sessions"
	"golang.org/x/crypto/bcrypt"
)

func validateEmail(email string) bool {
	_, err := mail.ParseAddress(email)
	return err == nil
}

type AuthenticationHandler struct {
	identity.User
	identity.SignInManager
	identity.UserStore
	sessions.Session
	handling.URLGenerator
	logging.Logger
	config.Configuration
	models.Repository
}

type AuthenticationTemplateContext struct {
	ErrMessage   string
	GitHubUrl    string
	GoogleUrl    string
	RecaptchaKey string
	ForgotUrl    string
	AllowSignup  bool
	ShowGoogle   bool
	ShowGithub   bool
}

func (handler AuthenticationHandler) GetSignIn() actionresults.ActionResult {
	message := handler.Session.GetValueDefault(utils.SIGNIN_MSG_KEY, "").(string)
	allowSignup := handler.Configuration.GetBoolDefault("authorization:allowsignup", false)
	showGoogle := handler.Configuration.GetBoolDefault("authorization:showgoogle", false)
	showGithub := handler.Configuration.GetBoolDefault("authorization:showgithub", false)
	recaptchaKey, _ := handler.Configuration.GetString("authorization:recaptchakey")
	forgoturl := utils.MustGenerateUrl(handler.URLGenerator, AuthenticationHandler.GetForgotPwd)
	//handler.Logger.Debugf("forgoturl: ", forgoturl  allowSignup)
	return actionresults.NewTemplateAction("auth_sign.html",
		AuthenticationTemplateContext{
			ErrMessage:   message,
			GitHubUrl:    "/oauth/github",
			GoogleUrl:    "/oauth/google",
			RecaptchaKey: recaptchaKey,
			ForgotUrl:    forgoturl,
			AllowSignup:  allowSignup,
			ShowGoogle:   showGoogle,
			ShowGithub:   showGithub,
		})
}

type userinfo struct {
	Login          string
	Name           string
	Email          string
	Verified_email string
	Picture        string
}

func (handler AuthenticationHandler) PostSignIn(in LoginRequest) actionresults.ActionResult {
	handler.Logger.Debugf("PostSignIn in", in)
	handler.Session.SetValue(utils.SIGNIN_PROVIDER, "local")
	var url string = "/"

	creds, err := handler.GetUserByEmail(strings.ToLower(in.Email))
	if err == nil {
		if creds.Id != 0 {
			handler.Logger.Debugf("PostSignIn User found:", creds)
			// check isvalidated and ispasswordok
			if IsPasswordOK(creds, in.Password) {
				if IsVerified(creds) {
					url = handler.SignInCredsUser(creds)
					handler.Session.SetValue(utils.SIGNIN_MSG_KEY, "")
					//handler.SignInManager.SignIn(creds)
					return actionresults.NewRedirectAction(url)
				} else {
					handler.Logger.Debugf("PostSignIn user is not verified")
				}
			} else {
				handler.Logger.Debugf("PostSignIn check password false")
				handler.Session.SetValue(utils.SIGNIN_MSG_KEY, "Wrong password")
			}
		} else {

			handler.Session.SetValue(utils.SIGNIN_MSG_KEY, "The specified username/password combination was not found")
		}
	} else {
		handler.Session.SetValue(utils.SIGNIN_MSG_KEY, err.Error())
	}

	if handler.Session.GetValueDefault(utils.SIGNIN_MSG_KEY, "").(string) == "" {
		handler.Session.SetValue(utils.SIGNIN_MSG_KEY, "Access Denied")
	}
	return actionresults.NewRedirectAction("/signin")
}

func (handler AuthenticationHandler) PostSignUp(in LoginRequest) actionresults.ActionResult {
	handler.Logger.Debugf("PostSignUp in", in)
	handler.Session.SetValue(utils.SIGNIN_PROVIDER, "local")
	var url string = "/signin"

	recaptchaSecret, _ := handler.Configuration.GetString("authorization:recaptchasecret")

	if err := CheckRecaptcha(recaptchaSecret, in.RecaptchaResponse, "signup"); err != nil {
		handler.Logger.Debugf("PostSignUp Bad reCaptcha", err.Error())
		handler.Session.SetValue(utils.SIGNIN_MSG_KEY, "Bad reCaptcha")

		return actionresults.NewRedirectAction(url)
	}

	if !validateEmail(in.Email) {
		handler.Logger.Debugf("PostSignUp Bad Email")
		handler.Session.SetValue(utils.SIGNIN_MSG_KEY, "Bad Email")
		return actionresults.NewRedirectAction(url)
	}
	if len(in.Password) == 0 {
		handler.Logger.Debugf("PostSignUp Empty password")
		handler.Session.SetValue(utils.SIGNIN_MSG_KEY, "Empty password")
		return actionresults.NewRedirectAction(url)
	}
	// Generate Verification Code
	code := utils.RandomString(64)
	verification_code := code
	// check if user exists
	creds, err := handler.GetUserByEmail(in.Email)
	if err == nil {
		if creds.Id == 0 {
			handler.Logger.Debugf("PostSignUp User not found:", creds)
			handler.Logger.Debugf("PostSignUp  verification_code: ", verification_code)

			// bcrypt password
			hash, err := bcrypt.GenerateFromPassword([]byte(in.Password), 13)
			if err != nil {
				handler.Session.SetValue(utils.SIGNIN_MSG_KEY, err.Error())
				handler.Logger.Debugf("PostSignUp An error arise after call bcrypt.GenerateFromPassword", err.Error())
			} else {
				creds.Password = string(hash)
				// Add User to Database
				creds.VerificationCode = verification_code
				creds.Email = in.Email
				creds.Username = strings.ToLower(in.Email)
				newUser, err := handler.AddNewUser(creds)
				if err != nil {
					handler.Session.SetValue(utils.SIGNIN_MSG_KEY, err.Error())
					handler.Logger.Debugf("PostSignUp An error arise after call AddNewUser: ", err.Error())
				} else {
					handler.Logger.Debugf("PostSignIn Up newUser: ", newUser)

					// ? Send Email
					// get user's name
					var firstName = newUser.GetDisplayName()
					if strings.Contains(strings.TrimSpace(firstName), " ") {
						firstName = strings.Split(firstName, " ")[1]
					}

					hostname, _ := handler.Configuration.GetString("system:hostname")
					emailData := utils.EmailTemplateData{
						URL:       hostname + utils.MustGenerateUrl(handler.URLGenerator, AuthenticationHandler.GetSignUpVerification, verification_code),
						FirstName: firstName,
						Subject:   "Your account verification code",
					}
					handler.Logger.Debugf("PostSignUp verification_code: ", emailData)
					err = utils.SendEmailVerification(handler.Configuration, handler.Logger, &newUser, &emailData)
					if err != nil {
						handler.Session.SetValue(utils.SIGNIN_MSG_KEY, err.Error())
						handler.Logger.Debugf("PostSignUp An error arise after call SendMail: ", err.Error())
					} else {
						handler.Session.SetValue(utils.SIGNIN_MSG_KEY, "We sent an email with a verification code to "+newUser.GetDisplayName())
					}
				}
			}
		} else {
			//check isverify
			if IsVerified(creds) {
				err := errors.New("A user with such credentials is already present")
				handler.Logger.Debugf("PostSignUp ", err.Error())
				handler.Session.SetValue(utils.SIGNIN_MSG_KEY, err.Error())
			} else {
				//sending email
				// get user's name
				var firstName = creds.GetDisplayName()
				if strings.Contains(strings.TrimSpace(firstName), " ") {
					firstName = strings.Split(firstName, " ")[1]
				}

				hostname, _ := handler.Configuration.GetString("system:hostname")
				emailData := utils.EmailTemplateData{
					URL:       hostname + utils.MustGenerateUrl(handler.URLGenerator, AuthenticationHandler.GetSignUpVerification, creds.VerificationCode),
					FirstName: firstName,
					Subject:   "Your account verification code",
				}
				handler.Logger.Debugf("PostSignUp verification_code: ", emailData)
				err = utils.SendEmailVerification(handler.Configuration, handler.Logger, &creds, &emailData)
				if err != nil {
					handler.Session.SetValue(utils.SIGNIN_MSG_KEY, err.Error())
					handler.Logger.Debugf("PostSignUp An error arise after call SendMail: ", err.Error())
				} else {
					handler.Session.SetValue(utils.SIGNIN_MSG_KEY, "We sent an email with a verification code to "+creds.GetDisplayName())
				}
			}
		}
	} else {
		handler.Logger.Debugf("PostSignUp An error arise after call GetUserByEmail: ", err.Error())
	}

	if handler.Session.GetValueDefault(utils.SIGNIN_MSG_KEY, "").(string) == "" {
		handler.Session.SetValue(utils.SIGNIN_MSG_KEY, "Access Denied")

	}
	return actionresults.NewRedirectAction(url)
}

// handler for endering user's email verification click
func (handler AuthenticationHandler) GetSignUpVerification(code string) actionresults.ActionResult {

	verification_code := code // utils.Encode(code)
	handler.Logger.Debugf("GetSignUpVerification code", verification_code)

	// get user by code
	user, err := handler.GetUserByVerificationCode(verification_code)
	if err != nil {
		handler.Logger.Debugf("GetSignUpVerification Error to get GetUserByVerificationCode: ", err.Error())
		handler.Session.SetValue(utils.SIGNIN_MSG_KEY, err.Error())
	} else {
		if user.Id == 0 {
			err := errors.New("Unknown Brucheion user, please sign up first.")
			handler.Session.SetValue(utils.SIGNIN_MSG_KEY, err.Error())
			handler.Logger.Debugf("", err.Error())
		} else {
			// if user already verified?
			if IsVerified(user) {
				url := handler.SignInCredsUser(user)
				return actionresults.NewRedirectAction(url)
			} else {
				if user.Password == "" { // forgot password
					handler.Session.SetValue(utils.SIGNIN_USER_ID, user.GetID())
					handler.Session.SetValue(utils.SIGNIN_USER_EMAIL, user.Email)
					handler.Session.SetValue(utils.SIGNIN_USER_VERIFICATION_CODE, verification_code)
					return actionresults.NewTemplateAction("auth_resetpwd.html",
						struct {
							Email    string
							Username string
							Password string
							Message  string
							// Id       int
						}{
							Email:    user.Email,
							Username: user.GetDisplayName(),
							Password: "",
							Message:  "",
						})
				} else { // verified new user
					handler.Session.SetValue(utils.SIGNIN_USER_ID, user.GetID())
					handler.Session.SetValue(utils.SIGNIN_USER_EMAIL, user.Email)
					handler.Session.SetValue(utils.SIGNIN_USER_VERIFICATION_CODE, verification_code)
					return actionresults.NewTemplateAction("auth_create.html",
						struct {
							Email    string
							Username string
							Message  string
							// Id       int
						}{
							Email:    user.Email,
							Username: user.GetDisplayName(), // user.username
							Message:  "",
							// Id:       user.GetID(),
						})
				}
			}
		}
	}
	handler.Logger.Debugf("GetSignUpVerification fail!")
	if handler.Session.GetValueDefault(utils.SIGNIN_MSG_KEY, "").(string) == "" {
		handler.Session.SetValue(utils.SIGNIN_MSG_KEY, "Something goes wrong...")
	}
	//return actionresults.NewRedirectAction("/signin")
	return actionresults.NewRedirectAction(utils.MustGenerateUrl(handler.URLGenerator,
		AuthenticationHandler.GetSignIn))
}

// receive user's data after clarification
func (handler AuthenticationHandler) PostSignUpVerification(creds models.Credentials) actionresults.ActionResult {
	vcode := handler.Session.GetValueDefault(utils.SIGNIN_USER_VERIFICATION_CODE, "").(string)
	handler.Logger.Debugf("PostSignUpVerification creds", creds, vcode)
	user, err := handler.GetUserByVerificationCode(vcode)
	if err != nil {
		handler.Logger.Debugf("PostSignUpVerification Error to get GetUserByVerificationCode: ", err.Error())
		handler.Session.SetValue(utils.SIGNIN_MSG_KEY, err.Error())
	}
	handler.Logger.Debugf("PostSignUpVerification user: ", user)
	// update user as verified
	creds.IsVerified = true
	creds.VerificationCode = ""
	creds.Id = user.Id       //handler.Session.GetValue(SIGNIN_USER_ID).(int)
	creds.Email = user.Email //strings.ToLower(handler.Session.GetValue(SIGNIN_USER_EMAIL).(string))
	creds.Username = strings.TrimSpace(creds.Username)
	creds.Password = user.Password

	handler.Logger.Debugf("PostSignUpVerification params:", creds)
	err = handler.UpdateUser(&creds)
	if err != nil {
		handler.Logger.Debugf("PostSignUpVerification Error to run UpdateUser: ", err.Error())
		handler.Session.SetValue(utils.SIGNIN_MSG_KEY, err.Error())
	} else {
		url := handler.SignInCredsUser(creds)
		return actionresults.NewRedirectAction(url)
	}
	handler.Logger.Debugf("PostSignUpVerification fail!")
	if handler.Session.GetValueDefault(utils.SIGNIN_MSG_KEY, "").(string) == "" {
		handler.Session.SetValue(utils.SIGNIN_MSG_KEY, "Something goes wrong...")
	}
	return actionresults.NewRedirectAction("/")
}

func (handler AuthenticationHandler) GetSignOut() actionresults.ActionResult {
	handler.Logger.Debugf("email session: ", handler.Session.GetValue(utils.SIGNIN_USER_EMAIL))
	handler.Session.Logout()
	handler.Logger.Debugf("email session: ", handler.Session.GetValue(utils.SIGNIN_USER_EMAIL))
	return actionresults.NewRedirectAction("/")
}

func (handler AuthenticationHandler) GetForgotPwd() actionresults.ActionResult {
	//message := handler.Session.GetValueDefault(SIGNIN_MSG_KEY, "").(string)
	return actionresults.NewTemplateAction("auth_forgot_pwd.html",
		struct {
			Email   string
			Message string
		}{
			Email:   "",
			Message: "",
		})
}

func (handler AuthenticationHandler) PostForgotPwd(in LoginRequest) actionresults.ActionResult {
	handler.Logger.Debugf("PostForgotPwd in:", in)
	var url string = "/"
	if in.Email == "" {
		handler.Logger.Debugf("PostForgotPwd Email not specified ")

	} else {
		err := utils.ResetPassword(handler.Repository, handler.Session, handler.URLGenerator, handler.Logger, handler.Configuration, in.Email)
		if err != nil {
			handler.Logger.Debugf("PostForgotPwd An error arise after call ResetPassword: ", err.Error())
		}
	}
	if handler.Session.GetValueDefault(utils.SIGNIN_MSG_KEY, "").(string) == "" {
		handler.Session.SetValue(utils.SIGNIN_MSG_KEY, "Access Denied")
	}
	return actionresults.NewRedirectAction(url)
}

func (handler AuthenticationHandler) PostResetPwd(in LoginRequest) actionresults.ActionResult {
	handler.Logger.Debugf("PostResetPwd in:", in)

	creds := handler.Repository.GetUserByID(handler.Session.GetValueDefault(utils.SIGNIN_USER_ID, 0).(int))
	if creds.Email != strings.ToLower(handler.Session.GetValueDefault(utils.SIGNIN_USER_EMAIL, "").(string)) {
		handler.Logger.Debugf("PostResetPwd Emails doesn't similar: ", strings.ToLower(handler.Session.GetValueDefault(utils.SIGNIN_USER_EMAIL, "").(string)))
		handler.Session.SetValue(utils.SIGNIN_MSG_KEY, "Error compare e-mails")
		return actionresults.NewRedirectAction("/")
	}
	// update user as verified
	creds.IsVerified = true
	creds.VerificationCode = ""
	// bcrypt password
	hash, err := bcrypt.GenerateFromPassword([]byte(in.Password), 13)
	if err != nil {
		handler.Logger.Debugf("PostResetPwd Error to generate password's hash: ", creds.Password)
		handler.Session.SetValue(utils.SIGNIN_MSG_KEY, "Error generate password's hash")
		return actionresults.NewRedirectAction("/")
	}

	creds.Password = string(hash)

	handler.Logger.Debugf("PostResetPwd params:", creds)
	if err := handler.UpdateUser(&creds); err != nil {
		handler.Logger.Debugf("PostResetPwd Error to run UpdateUser: ", err.Error())
		handler.Session.SetValue(utils.SIGNIN_MSG_KEY, err.Error())
	} else {
		url := handler.SignInCredsUser(creds)
		return actionresults.NewRedirectAction(url)
	}
	handler.Logger.Debugf("PostResetPwd fail!")
	if handler.Session.GetValue(utils.SIGNIN_MSG_KEY).(string) == "" {
		handler.Session.SetValue(utils.SIGNIN_MSG_KEY, "Something goes wrong...")
	}
	return actionresults.NewRedirectAction("/")
}

func IsVerified(u models.Credentials) bool {
	return u.IsVerified
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

func (handler AuthenticationHandler) PostSignOut(creds models.Credentials) actionresults.ActionResult {
	handler.Logger.Debugf("PostSignOut session: ", handler.Session.GetValue(utils.SIGNIN_USER_EMAIL))
	handler.Session.Logout()
	return actionresults.NewRedirectAction("/")
}
