package repo

import (
	"brucheion/gocite"
	"brucheion/models"
	"brucheion/utils"
	"encoding/json"
	"fmt"
	"unicode/utf8"

	"strings"
)

func (repo *SqlRepository) SavePassageTranscription(passage models.Passage, userid int) error {
	// {Bucket[]str, Data[]gocite.Work, Catalog[]BoltCatalog} ->
	//userID := repo.Contranstext.Value("USER_SESSION_KEY").(int)

	colid := passage.ColId
	passageid := passage.PassageID
	bucket := strings.Join(strings.Split(passageid, ":")[0:4], ":") + ":"
	transcriptions := passage.TranscriptionLines
	for i, l := range transcriptions {
		if strings.HasPrefix(l, "#") {
			_, size := utf8.DecodeRuneInString(l)
			if size > 0 {
				// Slice the string starting from the next rune
				transcriptions[i] = l[size:]
			} else {
				utils.RemoveSliceElementString(transcriptions, i)
			}
		}
	}

	textbyte := []byte(transcriptions[0])

	repo.Logger.Debugf("!!!!!!!!!!!!!!!SavePassageTranscription: ", textbyte)

	repo.Logger.Debugf("SavePassageTranscription: input ", userid, colid, passageid, bucket)

	transtext := strings.Join(transcriptions, "\r\n")

	origtranstext := transtext

	transtext = strings.ReplaceAll(transtext, "\r\n", "")

	// receive a passage: key (urn....:x.y.z) -> value
	retrieveddata, err := repo.SelectCollectionBucketKeyValue(colid, bucket, passageid, userid)
	if err != nil {
		//http.Error(w, "Internal server error", 500)
		repo.Logger.Debugf("SavePassageTranscription : Internal server error1 %v", err.Error())
		return err
	}
	retrievedjson := gocite.Passage{}
	json.Unmarshal([]byte(retrieveddata.JSON), &retrievedjson)
	retrievedjson.Text.Brucheion = transtext //gocite.Passage.Text.Brucheion is the transtext representation with newline tags
	retrievedjson.Text.TXT = origtranstext   //gocite.Passage.Text.TXT is the transtext representation with real line breaks instead of newline tags
	newnode, _ := json.Marshal(retrievedjson)

	catkey := []byte(passageid) //
	catvalue := []byte(newnode) //

	// rr, _ := repo.DB.BeginTx(repo.Contranstext, nil)
	tx, err := repo.DB.Begin()
	if err != nil {
		repo.Logger.Panicf("SavePassageTranscription Cannot create transaction: %v", err.Error())
		return err
	}

	repo.Logger.Info("SavePassageTranscription are starting...")
	//------------------------
	repo.Logger.Debugf("SavePassageTranscription are starting", bucket)
	repo.CreateBucketIfNotExists(tx, bucket, colid) //id_collection
	if err != nil {
		tx.Rollback()
		return err
	}

	//put bucket data (hstore)
	//repo.Logger.Debugf("SaveBoltData: id_col, newbucket, catkey, catvalue", id_col, newbucket, catkey, catvalue)
	_, err = tx.StmtContext(repo.Context, repo.Commands.SaveCiteDataDict).Exec(colid, bucket, catkey, catvalue)
	if err != nil {
		repo.Logger.Debugf("SavePassageTranscription: Cannot get inserted data1", err.Error())
		tx.Rollback()
		return err
	}
	err = tx.Commit()
	if err != nil {
		repo.Logger.Panicf("SavePassageTranscription Transaction cannot be committed: %v", err.Error())
		err = tx.Rollback()
		if err != nil {
			repo.Logger.Panicf("SavePassageTranscription Transaction cannot be rolled back: %v", err.Error())
			return err
		}
	}
	return nil
}

func (repo *SqlRepository) SavePassageText(passage models.PassageText, userid int) error {
	// {Bucket[]str, Data[]gocite.Work, Catalog[]BoltCatalog} ->
	//userID := repo.Contranstext.Value("USER_SESSION_KEY").(int)

	colid := passage.ColId
	passageid := passage.PassageID
	bucket := strings.Join(strings.Split(passageid, ":")[0:4], ":") + ":"

	transtext := passage.Text
	if len(transtext) == 0 {
		repo.Logger.Debugf("SavePassageText : Empty transcription!")
		return fmt.Errorf("Empty transcription!")
	}
	origtranstext := strings.ReplaceAll(transtext, "\n", "\r\n") //correct textarea(return chr10 only)!!!
	//origtranstext = strings.ReplaceAll(transtext, "\r\r", "\r") //correct textarea if didn't edit
	repo.Logger.Debugf("SavePassageText : ", []byte(transtext))
	transtext = strings.ReplaceAll(origtranstext, "\r\n", "")

	// receive a passage: key (urn....:x.y.z) -> value
	retrieveddata, err := repo.SelectCollectionBucketKeyValue(colid, bucket, passageid, userid)
	if err != nil {
		//http.Error(w, "Internal server error", 500)
		repo.Logger.Debugf("SavePassageText : Internal server error1 %v", err.Error())
		return err
	}
	retrievedjson := gocite.Passage{}
	json.Unmarshal([]byte(retrieveddata.JSON), &retrievedjson)
	retrievedjson.Text.Brucheion = transtext //gocite.Passage.Text.Brucheion is the transtext representation with newline tags
	retrievedjson.Text.TXT = origtranstext   //gocite.Passage.Text.TXT is the transtext representation with real line breaks instead of newline tags
	newnode, _ := json.Marshal(retrievedjson)
	//repo.Logger.Debugf("SavePassageText", []byte(origtranstext), newnode)
	catkey := []byte(passageid) //
	catvalue := []byte(newnode) //

	// rr, _ := repo.DB.BeginTx(repo.Contranstext, nil)
	tx, err := repo.DB.Begin()
	if err != nil {
		repo.Logger.Panicf("SavePassageText Cannot create transaction: %v", err.Error())
		return err
	}

	repo.Logger.Info("SavePassageText are starting...")
	//------------------------
	repo.Logger.Debugf("SavePassageText are starting", bucket)
	repo.CreateBucketIfNotExists(tx, bucket, colid) //id_collection
	if err != nil {
		tx.Rollback()
		return err
	}

	//put bucket data (hstore)
	//repo.Logger.Debugf("SaveBoltData: id_col, newbucket, catkey, catvalue", id_col, newbucket, catkey, catvalue)
	_, err = tx.StmtContext(repo.Context, repo.Commands.SaveCiteDataDict).Exec(colid, bucket, catkey, catvalue)
	if err != nil {
		repo.Logger.Debugf("SavePassageText: Cannot get inserted data1", err.Error())
		tx.Rollback()
		return err
	}
	err = tx.Commit()
	if err != nil {
		repo.Logger.Panicf("SavePassageText Transaction cannot be committed: %v", err.Error())
		err = tx.Rollback()
		if err != nil {
			repo.Logger.Panicf("SavePassageText Transaction cannot be rolled back: %v", err.Error())
			return err
		}
	}
	return nil
}
