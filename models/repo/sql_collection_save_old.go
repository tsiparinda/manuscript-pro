package repo

import (
	"brucheion/gocite"
	"brucheion/models"
	"encoding/json"

	"strings"
)

// functions for working with collection data
// SaveTranscription saves a transcription to the database  and fired from BlockEditor
// SaveReference saves a image reference to the database

func (repo *SqlRepository) SaveTranscription(coltrans models.Transcription, userid int) error {
	// {Bucket[]str, Data[]gocite.Work, Catalog[]BoltCatalog} ->
	//userID := repo.Context.Value("USER_SESSION_KEY").(int)

	repo.Logger.Debugf("SaveTranscription: userID= ", userid)
	colid := coltrans.ColID
	newkey := coltrans.CTSURN
	newbucket := strings.Join(strings.Split(newkey, ":")[0:4], ":") + ":"
	text := coltrans.Transcription
	linetext := text
	repo.Logger.Debugf("SaveTranscription: id_collection= ")
	text = strings.Replace(text, "\r\n", "", -1)

	//retrieveddata, _ := BoltRetrieve(dbname, newbucket, newkey)
	// receive a passage: key (urn....:x.y.z) -> value
	retrieveddata, err := repo.SelectCollectionBucketKeyValue(colid, newbucket, newkey, userid)
	if err != nil {
		//http.Error(w, "Internal server error", 500)
		repo.Logger.Debugf("SaveTranscription : Internal server error1 %v", err.Error())
		return err
	}
	retrievedjson := gocite.Passage{}
	json.Unmarshal([]byte(retrieveddata.JSON), &retrievedjson)
	retrievedjson.Text.Brucheion = text //gocite.Passage.Text.Brucheion is the text representation with newline tags
	retrievedjson.Text.TXT = linetext   //gocite.Passage.Text.TXT is the text representation with real line breaks instead of newline tags
	newnode, _ := json.Marshal(retrievedjson)

	catkey := []byte(newkey)    //
	catvalue := []byte(newnode) //

	// rr, _ := repo.DB.BeginTx(repo.Context, nil)
	tx, err := repo.DB.Begin()
	if err != nil {
		repo.Logger.Panicf("SaveTranscription Cannot create transaction: %v", err.Error())
		return err
	}

	repo.Logger.Info("SaveTranscription are starting...")
	//------------------------
	repo.Logger.Debugf("SaveTranscription are starting", newbucket)
	repo.CreateBucketIfNotExists(tx, newbucket, colid) //id_collection
	if err != nil {
		tx.Rollback()
		return err
	}

	//put bucket data (hstore)
	//repo.Logger.Debugf("SaveBoltData: id_col, newbucket, catkey, catvalue", id_col, newbucket, catkey, catvalue)
	_, err = tx.StmtContext(repo.Context, repo.Commands.SaveCiteDataDict).Exec(colid, newbucket, catkey, catvalue)
	if err != nil {
		repo.Logger.Debugf("SaveTranscription: Cannot get inserted data1", err.Error())
		tx.Rollback()
		return err
	}
	//-----------------------

	err = tx.Commit()
	if err != nil {
		repo.Logger.Panicf("SaveTranscription Transaction cannot be committed: %v", err.Error())
		err = tx.Rollback()
		if err != nil {
			repo.Logger.Panicf("SaveTranscription Transaction cannot be rolled back: %v", err.Error())
			return err
		}
	}
	return nil
}

func (repo *SqlRepository) SaveReference(coltrans models.Transcription, userid int) error {

	repo.Logger.Debugf("SaveReference: userID= ", userid)
	colid := coltrans.ColID
	newkey := coltrans.CTSURN
	newbucket := strings.Join(strings.Split(newkey, ":")[0:4], ":") + ":"

	imageref := coltrans.ImageRef

	retrieveddata, err := repo.SelectCollectionBucketKeyValue(colid, newbucket, newkey, userid)
	if err != nil {
		//http.Error(w, "Internal server error", 500)
		repo.Logger.Debugf("SaveReference : Internal server error1 %v", err.Error())
		return err
	}

	retrievedjson := gocite.Passage{}
	json.Unmarshal([]byte(retrieveddata.JSON), &retrievedjson)
	var textareas []gocite.Triple
	for i := range imageref {
		textareas = append(textareas, gocite.Triple{Subject: newkey,
			Verb:   "urn:cite2:dse:verbs.v1:appears_on",
			Object: imageref[i]})
	}
	retrievedjson.ImageLinks = textareas
	newnode, _ := json.Marshal(retrievedjson)

	catkey := []byte(newkey) //
	catvalue := newnode      //

	// rr, _ := repo.DB.BeginTx(repo.Context, nil)
	tx, err := repo.DB.Begin()
	if err != nil {
		repo.Logger.Panicf("SaveReference Cannot create transaction: %v", err.Error())
		return err
	}

	//------------------------
	repo.Logger.Debugf("SaveReference are starting", newbucket)
	repo.CreateBucketIfNotExists(tx, newbucket, colid) //id_collection
	if err != nil {
		tx.Rollback()
		return err
	}

	//put bucket data (hstore)
	//repo.Logger.Debugf("SaveBoltData: id_col, newbucket, catkey, catvalue", id_col, newbucket, catkey, catvalue)
	_, err = tx.StmtContext(repo.Context, repo.Commands.SaveCiteDataDict).Exec(colid, newbucket, catkey, catvalue)
	if err != nil {
		repo.Logger.Debugf("SaveReference: Cannot get inserted data1", err.Error())
		tx.Rollback()
		return err
	}
	//-----------------------

	err = tx.Commit()
	if err != nil {
		repo.Logger.Panicf("SaveReference Transaction cannot be committed: %v", err.Error())
		err = tx.Rollback()
		if err != nil {
			repo.Logger.Panicf("SaveReference Transaction cannot be rolled back: %v", err.Error())
			return err
		}
	}
	return nil
}
