package utils

import (
	"brucheion/models"
	"strings"

	"github.com/vedicsociety/platform/config"
	"github.com/vedicsociety/platform/http/handling"
	"github.com/vedicsociety/platform/logging"
	"github.com/vedicsociety/platform/sessions"
)

func ResetPassword(repo models.Repository, session sessions.Session, urlgenerator handling.URLGenerator, logger logging.Logger, config config.Configuration, email string) error {

	// Generate Verification Code
	code := RandomString(64)
	verification_code := code
	// check if user exists
	creds, err := repo.GetUserByEmail(strings.ToLower(strings.TrimSpace(email)))
	if err == nil {
		if creds.Id == 0 {
			logger.Debugf("PostForgotPwd User not found:", creds)
			session.SetValue(SIGNIN_MSG_KEY, "User with this email address not found")
		} else {
			// update user: isverified=false, password="", verificationcode=new
			creds.IsVerified = false
			creds.VerificationCode = verification_code
			creds.Password = ""

			//handler.Logger.Debugf("PostForgotPwd params:", creds)
			err := repo.UpdateUser(&creds)
			if err != nil {
				logger.Debugf("PostForgotPwd Error to run UpdateUser: ", err.Error())
				session.SetValue(SIGNIN_MSG_KEY, err.Error())
			} else {

				//sending email
				// get user's name
				var firstName = creds.GetDisplayName()
				if strings.Contains(strings.TrimSpace(firstName), " ") {
					firstName = strings.Split(firstName, " ")[1]
				}

				hostname, _ := config.GetString("system:hostname")
				emailData := EmailTemplateData{

					URL:       hostname + "/signupverification/" + verification_code,
					FirstName: firstName,
					Subject:   "Your account verification code",
				}
				logger.Debugf("PostForgotPwd emailData: ", emailData)
				err = SendEmailVerification(config, logger, &creds, &emailData)
				if err != nil {
					config.SetValue(SIGNIN_MSG_KEY, err.Error())
					logger.Debugf("PostForgotPwd An error arise after call SendMail: ", err.Error())
					return err
				} else {
					session.SetValue(SIGNIN_MSG_KEY, "We sent an email with a verification code to "+creds.GetDisplayName())
				}
			}
		}
	} else {
		return err
	}
	return nil
}
