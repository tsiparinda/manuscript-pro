package repo

import (
	"brucheion/models"
	"fmt"
)

func (repo *SqlRepository) SelectCollectionBucketDictionary(colid int, urn string, userid int) (result []models.BucketDict) {

	rows, err := repo.Commands.SelectCollectionBucketDictionary.QueryContext(repo.Context, colid, urn, userid)
	if err == nil {
		if result, err = scanDict(rows); err != nil {
			repo.Logger.Panicf("Cannot scan data: %v", err.Error())
			return
		}
	} else {
		repo.Logger.Panicf("Cannot exec SelectUserBucketDict command: %v", err)
	}
	return
}

func (repo *SqlRepository) SelectCollectionBuckets(colid int, userid int) (values []string) {

	rows, err := repo.Commands.SelectCollectionBuckets.QueryContext(repo.Context, colid, userid)
	if err == nil {
		if values, err = scanStrings(rows); err != nil {
			repo.Logger.Panicf("scanStrings: Cannot scan data: %v", err.Error())
			return
		}
	} else {
		repo.Logger.Panicf("Cannot exec SelectCollectionBuckets command: %v", err)
	}
	return
}

func (repo *SqlRepository) SelectCollectionBucketKeyValue(colid int, urn string, key string, userid int) (value models.BoltJSON, err error) {

	rows, err := repo.Commands.SelectCollectionBucketKeyValue.QueryContext(repo.Context, colid, urn, key, userid)
	if err == nil {
		if value.JSON, err = scanBucketKeyValue(rows); err != nil {
			err = fmt.Errorf("scanBucketKeyValue: Cannot scan data: %v", err.Error())
			return
		}
	} else {
		repo.Logger.Panicf("Cannot exec SelectUserBucketKeyValue command: %v", err)
	}
	return
}

func (repo *SqlRepository) SelectCollectionBucketKeys(colid int, bucket string, userid int) (value []string, err error) {

	rows, err := repo.Commands.SelectCollectionBucketKeys.QueryContext(repo.Context, colid, bucket, userid)
	if err == nil {
		if value, err = scanBucketKeys(rows); err != nil {
			err = fmt.Errorf("SelectCollectionBucketKeys scanBucketKeys: Cannot scan data: %v", err.Error())
			return
		}
	} else {
		repo.Logger.Panicf("Cannot exec SelectUserBucketKeys command: %v", err)
	}
	return
}
