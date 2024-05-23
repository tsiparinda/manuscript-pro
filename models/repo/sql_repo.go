package repo

import (
	"context"
	"database/sql"

	"github.com/vedicsociety/platform/authorization/identity"
	"github.com/vedicsociety/platform/config"
	"github.com/vedicsociety/platform/logging"
)

// SqlRepository encapsulates the methods required for communicating with a SQL database.
type SqlRepository struct {
	config.Configuration
	logging.Logger
	Commands SqlCommands
	identity identity.User
	*sql.DB
	context.Context
}

// SqlCommands aggregates all SQL command statements to be executed against the database.
// All new SQL statements should be added here.
type SqlCommands struct {
	// User related SQL commands
	AddUserToGroup,
	AddNewUser,
	AddUserAdmin,
	GetUserByID,
	GetUserByName,
	GetUserByEmail,
	GetUserByVerificationCode,
	GetCredentials,
	GetUsers,
	GetGroups,
	GetUserGroups,
	GetGroupById,
	GetGroupByName,
	UpdateUser,
	DeleteUserGroup,
	UpdateGroup,

	// Collection related SQL commands
	SelectCollectionBucketDictionary,
	SelectCollectionBuckets,
	SelectCollectionBucketKeyValue,
	SelectCollectionImageDictionary,
	SelectCollectionImageKeyValue,
	SelectCollectionBucketKeys,
	SaveCiteDataDict,
	SaveImageDataDict,
	CreateBucketIfNotExists,

	// Pagination related SQL commands
	SelectCollectionsPage,
	SelectCollectionsPageCount,
	GetAuthorPage,
	GetAuthorPageCount,

	// Other SQL commands
	SelectCollection,
	AddNewCollection,
	DeleteCollection,
	UpdateCollection,
	IsCollectionWriteble,
	SelectColUsers,
	DropCollectionsUser,
	DropCollectionUsers,
	InsertColUsers,
	SelectImageCollectionList *sql.Stmt
}
