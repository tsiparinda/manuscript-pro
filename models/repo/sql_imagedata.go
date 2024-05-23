package repo

import (
	"brucheion/models"
	"encoding/json"
	"fmt"
)

// functions for working with image data

func (repo *SqlRepository) SaveImageData(userid int, image *models.Image) error {
	repo.Logger.Debugf("SaveImageData: image", image)

	collection := models.ImageCollection{}
	//newImage := models.Image{}
	// get value from imagedata by image.colname as key
	val, err := repo.LoadCollectionImageKeyValue(image.ColId, image.Name, userid)
	if err != nil {
		return err
	}
	repo.Logger.Debugf("SaveImageData: val", val)
	json.Unmarshal([]byte(val.JSON), &collection.Collection)
	repo.Logger.Debugf("SaveImageData: collection", collection)
	found := false
	// check if image is already in collection
	for coli, colv := range collection.Collection {
		if colv.URN == image.URN {
			found = true
			collection.Collection[coli] = *image
		}
	}
	// if not, add it to the collection
	if !found {
		collection.Collection = append(collection.Collection, *image)
	}
	repo.Logger.Debugf("SaveImageData: collection", collection)
	value, _ := json.Marshal(collection.Collection)
	repo.Logger.Debugf("SaveImageData: image.Name, value", image.Name, value)
	// save image collection to imagedata
	err = repo.SaveImageDataDict(image.ColId, image.Name, value)
	if err != nil {
		repo.Logger.Debugf("SaveImageData: err", err.Error())
	}

	return nil
}

func (repo *SqlRepository) LoadCollectionImageKeyValue(colid int, key string, userid int) (value models.BoltJSON, err error) {
	repo.Logger.Debugf("LoadCollectionImageKeyValue input: ", colid,key,userid)
	rows, err := repo.Commands.SelectCollectionImageKeyValue.QueryContext(repo.Context, colid, key, userid)
	repo.Logger.Debugf("LoadCollectionImageKeyValue : ", rows)
	if err == nil {
		if value.JSON, err = scanImageKeyValue(rows); err != nil {
			err = fmt.Errorf("scanImageKeyValue: Cannot scan data: %v", err.Error())
			return
		}
	} else {
		repo.Logger.Panicf("Cannot exec LoadCollectionImageKeyValue command: %v", err)
	}
	return
}

func (repo *SqlRepository) SaveImageDataDict(id_col int, catkey string, catvalue []byte) error {
	repo.Logger.Debugf("SaveImageDataDict: ", id_col, catkey, catvalue)
	_, err := repo.Commands.SaveImageDataDict.Exec(id_col, catkey, catvalue)
	if err != nil {
		repo.Logger.Panicf("SaveImageDataDict Cannot exec SaveImageDataDict command: %v", err.Error())
		return err
	}

	return nil
}

func (repo *SqlRepository) LoadCollectionImageDictionary(colid int,  userid int) (result []models.BucketDict) {

	rows, err := repo.Commands.SelectCollectionImageDictionary.QueryContext(repo.Context, colid, userid)
	if err == nil {
		if result, err = scanDict(rows); err != nil {
			repo.Logger.Panicf("LoadCollectionImageDictionary Cannot scan data: %v", err.Error())
			return
		}
	} else {
		repo.Logger.Panicf("Cannot exec LoadCollectionImageDictionary command: %v", err)
	}
	return
}

func (repo *SqlRepository) LoadImageCollectionList() (result []string) {

	rows, err := repo.Commands.SelectImageCollectionList.QueryContext(repo.Context)
	if err == nil {
		if result, err = scanStrings(rows); err != nil {
			repo.Logger.Panicf("LoadImageCollectionList Cannot scan data: %v", err.Error())
			return
		}
	} else {
		repo.Logger.Panicf("Cannot exec LoadImageCollectionList command: %v", err)
	}
	return
}