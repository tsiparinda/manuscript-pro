package auth

import (
	"brucheion/models"

	"github.com/vedicsociety/platform/authorization/identity"
	"github.com/vedicsociety/platform/services"
)

func RegisterUserStoreService() {
	err := services.AddSingleton(func() identity.UserStore {
		return &userStore{}
	})
	if err != nil {
		panic(err)
	}
}

type userStore struct {
}

func (store *userStore) GetUserByID(id int) (identity.User, bool) {

	var repo models.Repository
	services.GetService(&repo)
	user := models.Repository.GetUserByID(repo, id)
	if user.GetID() != 0 {
		return &user, true
	}
	return &user, false
}

func (store *userStore) GetUserByName(name string) (identity.User, bool) {

	var repo models.Repository
	services.GetService(&repo)
	user := models.Repository.GetUserByName(repo, name)
	if user.GetID() != 0 {
		return &user, true
	}
	return &user, false
}
