package models

import (
	"strings"
)

type Credentials struct {
	Id               int
	Username         string
	Email            string
	Description      string
	Password         string
	VerificationCode string
	IsVerified       bool
	// Roles            []string
	Roles []Group
}

func (u Credentials) GetID() int {
	//	fmt.Printf("user id %v \n", u.Id)
	return u.Id
}

func (u Credentials) GetDisplayName() (name string) {
	if strings.TrimSpace(u.Username) == "" {
		name = strings.Split(u.Email, "@")[0]
	} else {
		name = u.Username
	}
	return
}

func (u Credentials) GetEmail() string {
	return strings.ToLower(u.Email)
}

func (u Credentials) InRole(role string) bool {
	for _, r := range u.Roles {
		// fmt.Printf("InRole  r %v\n", r)
		if strings.EqualFold(r.Name, role) {
			return true
		}
	}
	return false
}

func (u Credentials) IsAuthenticated() bool {
	return u.IsVerified
}
