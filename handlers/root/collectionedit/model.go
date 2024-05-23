package collectionedit

import (
	"brucheion/models"

	"github.com/vedicsociety/platform/authorization/identity"
	"github.com/vedicsociety/platform/http/handling"
	"github.com/vedicsociety/platform/logging"
)

// CollectionEditHandler is a structure that combines multiple
// interfaces required to handle edit operations on collections.
// It includes a Repository for data operations, a URLGenerator 
// for URL handling, a Logger for logging operations, and User 
// to represent the currently authenticated user.
type CollectionEditHandler struct {
	models.Repository
	handling.URLGenerator
	logging.Logger
	identity.User
}
