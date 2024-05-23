package models

import (
	"database/sql"

	"github.com/vedicsociety/platform/http/handling/params"
)

// Repository defines the interface for performing CRUD operations on database entities.
// It provides methods for handling users, groups, collections, transcriptions, and images.
// all new functions must be added to this interface
type Repository interface {
	AddUserToGroup(creds *Credentials, group int) error
	GetGroupById(id int) Group
	GetGroupByName(string) Group
	GetUserByID(userid int) Credentials
	GetUserByName(username string) Credentials
	GetUserByEmail(string) (Credentials, error)
	GetUserGroups(*Credentials) error
	AddNewUser(creds Credentials) (Credentials, error)
	AddUserAdmin(creds Credentials) (Credentials, error)
	GetUserByVerificationCode(code string) (Credentials, error)
	UpdateUser(creds *Credentials) error
	UpdateGroup(*Group) error
	UpdateUserGroups(*Credentials) error
	GetCredentials() []Credentials
	GetUsers(int) []User
	GetGroups() []Group

	SelectCollectionBucketDictionary(colid int, urn string, userid int) []BucketDict
	SelectCollectionBuckets(colid int, userid int) []string
	SelectCollectionBucketKeyValue(colid int, urn string, key string, userid int) (BoltJSON, error)
	SelectCollectionBucketKeys(colid int, urn string, userid int) (value []string, err error)

	LoadCollectionImageDictionary(colid int, userid int) (result []BucketDict)
	LoadCollectionImageKeyValue(colid int, key string, userid int) (value BoltJSON, err error)
	SaveImageData(userid int, image *Image) error

	SaveBoltData(*BoltData) (int, error)
	SaveCiteDataDict(id_col int, bucket string, catkey string, catvalue []byte) error
	CreateBucketIfNotExists(*sql.Tx, string, int) error
	AddNewCollection(params.InputParams, int) (int, error)
	LoadCollectionsPageAuthor(authorid, categoryId, page, pageSize int) (collections []CollectionPage, totalAvailable int)
	DropCollection(int, int) error

	LoadTranscriptionForEdit(int, int) Collection
	SaveTranscription(Transcription, int) error
	SavePassageTranscription(Passage, int) error
	SavePassageText(PassageText, int) error
	SaveReference(Transcription, int) error

	LoadCollection(colid int, userid int) Collection
	SaveCollection(Collection, int) error
	LoadCollectionForShare(int, int) ShareCollection
	IsCollectionWriteble(int, int) (bool, error)

	LoadColUsers(int, int) []ColUsers
	DropCollectionsUser(int, int) error
	DropCollectionUsers(int, int) error
	AddColUsers(ColUsers, int) error

	LoadImageCollectionList() (result []string)
}
