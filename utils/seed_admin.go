package utils

import (
	"brucheion/models"

	"github.com/vedicsociety/platform/config"
	"github.com/vedicsociety/platform/logging"
	"golang.org/x/crypto/bcrypt"
)

func SeedAdmin(cfg config.Configuration, log logging.Logger, repo models.Repository) {
	// init oauth providers
	adminname, _ := cfg.GetString("system:adminName")
	adminemail, _ := cfg.GetString("system:adminEmail")
	adminpassword, _ := cfg.GetString("system:adminPassword")
	hash, err := bcrypt.GenerateFromPassword([]byte(adminpassword), 13)
	if err != nil {
		log.Debugf("SeedAdmin: An error arise after call bcrypt.GenerateFromPassword", err.Error())
	}
	creds := models.Credentials{}

	creds.Password = string(hash)
	creds.IsVerified = true
	creds.Email = adminemail
	creds.Username = adminname
	creds.Id = 1

	// check if admin exists (id=1)
	admin := repo.GetUserByID(1)
	if admin.Id == 1 {
		log.Debug("SeedAdmin: User with ID=1 already exists")
		return
	}
	// add new user
	newUser, err := repo.AddUserAdmin(creds)
	if err != nil {
		log.Debugf("SeedAdmin: An error arise after call AddNewUser: ", err.Error())
	} else {
		log.Infof("SeedAdmin: User %v was added to DB", newUser.Username)
	}
}
