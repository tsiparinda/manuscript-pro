package utils

import (
	"brucheion/models"
	"bytes"
	"crypto/tls"
	"html/template"
	"os"
	"path/filepath"

	"github.com/vedicsociety/platform/config"
	"github.com/vedicsociety/platform/logging"

	gomail "gopkg.in/gomail.v2"
)

type EmailTemplateData struct {
	URL       string
	FirstName string
	Subject   string
}

type EmailContext struct {
	URL          string
	FirstName    string
	Subject      string
	Msg          string
	CallToAction string
}

// ? Email template parser
func ParseTemplateDir(dir string) (*template.Template, error) {
	var paths []string
	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			paths = append(paths, path)
		}
		return nil
	})

	if err != nil {
		return nil, err
	}

	return template.ParseFiles(paths...)
}

func SendEmailVerification(config config.Configuration, logger logging.Logger, user *models.Credentials, data *EmailTemplateData) (err error) {
	emailData := EmailContext{
		URL:          data.URL,
		FirstName:    data.FirstName,
		Subject:      data.Subject,
		Msg:          "Thanks for signing up. Click the link below to confirm your registration and you'll be on your way.",
		CallToAction: "Confirm your registration",
	}
	return SendEmail(config, logger, user, &emailData)
}

func SendEmail(config config.Configuration, logger logging.Logger, user *models.Credentials, data *EmailContext) (err error) {

	// Sender data.
	from, _ := config.GetString("sendmail:from")
	smtpPass, _ := config.GetString("sendmail:smtpPass")
	smtpUser, _ := config.GetString("sendmail:smtpUser")
	to := user.Email
	smtpHost, _ := config.GetString("sendmail:smtpHost")
	smtpPort, _ := config.GetInt("sendmail:smtpPort")
	tmpldir, _ := config.GetString("sendmail:tmplatesdir")

	var body bytes.Buffer

	template, err := ParseTemplateDir(tmpldir)
	if err != nil {
		logger.Debugf("Could not parse template!", err)
		return
	}

	template.ExecuteTemplate(&body, "email.html", &data)

	m := gomail.NewMessage()

	m.SetHeader("From", from)
	m.SetHeader("To", to)
	m.SetHeader("Subject", data.Subject)
	m.SetBody("text/html", body.String())
	//m.AddAlternative("text/plain", html2text.HTML2Text(body.String()))
	logger.Debugf("Sending email: ", smtpHost, smtpPort, smtpUser, smtpPass)
	d := gomail.NewDialer(smtpHost, smtpPort, smtpUser, smtpPass)
	d.TLSConfig = &tls.Config{InsecureSkipVerify: true}

	// Send Email
	for n := 0; n < 3; n++ {
		//  if you use gmail, you should enable "Less auth application" in gmail account settings!
		logger.Debugf("Try to call SendMail attempt number %d", n+1)
		err = d.DialAndSend(m)
		if err == nil {
			break
		}
	}
	if err != nil {
		logger.Debugf("Could not send email: ", err)
	}
	return err
}
