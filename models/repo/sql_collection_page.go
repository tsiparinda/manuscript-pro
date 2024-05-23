package repo

import (
	"brucheion/models"
	"strings"
)

// functions for working with collections

func (repo *SqlRepository) LoadCollectionsPage(userid, authorid, page, pageSize int) (collections []models.CollectionPage, totalAvailable int) {
	rows, err := repo.Commands.SelectCollectionsPage.QueryContext(repo.Context, userid, authorid, pageSize, (pageSize*page)-pageSize)
	if err == nil {
		if collections, err = scanCollectionsPage(rows); err != nil {
			repo.Logger.Panicf("Cannot scan data: %v", err.Error())
			return
		}
	} else {
		repo.Logger.Panicf("Cannot exec LoadCollectionsPage command: %v", err)
		return
	}
	row := repo.Commands.SelectCollectionsPageCount.QueryRowContext(repo.Context)
	if row.Err() == nil {
		if err := row.Scan(&totalAvailable); err != nil {
			repo.Logger.Panicf("Cannot scan data: %v", err.Error())
		}
	} else {
		repo.Logger.Panicf("Cannot exec GetCollectionPageCount command: %v", row.Err().Error())
	}
	return
}

func (repo *SqlRepository) LoadCollectionsPageAuthor(userid, authorId, page, pageSize int) (collections []models.CollectionPage, totalAvailable int) {

	repo.Logger.Debugf("GetCollections userid ", userid)

	repo.Logger.Debugf("GetCollectionPage section ", authorId, page, pageSize)
	return repo.LoadCollectionsPage(userid, authorId, page, pageSize)

}

func (repo *SqlRepository) DropCollection(colid int, userid int) error {
	_, err := repo.Commands.DeleteCollection.Exec(colid, userid)
	if err != nil {
		repo.Logger.Panicf("DropCollection Cannot exec DropCollection command: %v", err.Error())
		return err
	}

	return nil
}

func (repo *SqlRepository) SaveCollection(col models.Collection, userid int) error {
	_, err := repo.Commands.UpdateCollection.Exec(col.Id, userid, col.Title, col.IsPublic)
	if err != nil {
		repo.Logger.Panicf("UpdateCollection Cannot exec DropCollection command: %v", err.Error())
		return err
	}

	return nil
}

func (repo *SqlRepository) LoadCollection(colid int, userid int) (Collection models.Collection) {
	repo.Logger.Debugf("LoadCollection", colid, userid)
	rows, err := repo.Commands.SelectCollection.QueryContext(repo.Context, colid, userid)
	if err == nil {
		if Collection, err = scanCollection(rows); err != nil {
			repo.Logger.Panicf("LoadCollection Cannot scanCollection: %v", err.Error())
			return models.Collection{}
		}
	} else {
		repo.Logger.Panicf("LoadCollection Cannot exec SelectCollection command: %v", err)
		return models.Collection{}
	}
	// fill the buckets array
	buckets := repo.SelectCollectionBuckets(colid, userid)
	if len(buckets) != 0 {
		for _, bucket := range buckets {
			citations, _ := repo.SelectCollectionBucketKeys(colid, bucket, userid)
			var citations_short []string
			if len(citations) != 0 {
				for _, citation := range citations {
					citation = strings.Join(strings.Split(citation, ":")[4:5], ":")
					if len([]rune(citation)) > 0 {
						citations_short = append(citations_short, citation)
					}
				}
				Collection.Buckets = append(Collection.Buckets, models.Bucket{Bucket: bucket, Citations: citations_short})
			}
		}
	} // fill the passages array
	return Collection
}

func (repo *SqlRepository) LoadColUsers(colid int, userid int) (ColUsers []models.ColUsers) {
	rows, err := repo.Commands.SelectColUsers.QueryContext(repo.Context, colid, userid)
	if err == nil {
		if ColUsers, err = scanColUsers(rows); err != nil {
			repo.Logger.Panicf("LoadColUsers Cannot scan data: %v", err.Error())
			return []models.ColUsers{}
		}
	} else {
		repo.Logger.Panicf("LoadColUsers Cannot exec SelectColUsers command: %v", err)
		return []models.ColUsers{}
	}
	return ColUsers
}

func (repo *SqlRepository) LoadCollectionForShare(colid int, userid int) (ShareCollection models.ShareCollection) {
	// we need check if user can share collection
	ShareCollection.Collection = repo.LoadCollection(colid, userid)
	repo.Logger.Debugf("LoadCollectionForShare ShareCollection.Collection: ", ShareCollection.Collection)
	if ShareCollection.Collection.AuthorId != userid {
		return models.ShareCollection{}
	}
	ShareCollection.ColUsers = repo.LoadColUsers(colid, userid)
	return ShareCollection
}

func (repo *SqlRepository) LoadTranscriptionForEdit(colid int, userid int) (EditCollection models.Collection) {
	rows, err := repo.Commands.SelectCollection.QueryContext(repo.Context, colid, userid)
	if err == nil {
		if EditCollection, err = scanCollection(rows); err != nil {
			repo.Logger.Panicf("LoadCollectionForEdit Cannot scanCollection: %v", err.Error())
			return models.Collection{}
		}
	} else {
		repo.Logger.Panicf("LoadCollectionForEdit Cannot exec SelectCollection command: %v", err)
		return models.Collection{}
	}

	return EditCollection
}

func (repo *SqlRepository) DropCollectionUsers(colid int, userid int) error {
	_, err := repo.Commands.DropCollectionUsers.Exec(colid, userid)
	if err != nil {
		repo.Logger.Panicf("DropCollectionUsers Cannot exec DeleteColUsers command: %v", err.Error())
		return err
	}

	return nil
}

func (repo *SqlRepository) DropCollectionsUser(colid int, userid int) error {
	_, err := repo.Commands.DropCollectionsUser.Exec(colid, userid)
	if err != nil {
		repo.Logger.Panicf("DropCollectionsUser Cannot exec DeleteColUsers command: %v", err.Error())
		return err
	}

	return nil
}

func (repo *SqlRepository) AddColUsers(colusers models.ColUsers, userid int) error {
	_, err := repo.Commands.InsertColUsers.Exec(colusers.Id_Col, colusers.Id_User, colusers.Is_Write)
	if err != nil {
		repo.Logger.Panicf("AddColUsers Cannot exec InsertColUsers command: %v", err.Error())
		return err
	}

	return nil
}

func (repo *SqlRepository) IsCollectionWriteble(colid, userid int) (bool, error) {
	var f bool
	row := repo.Commands.IsCollectionWriteble.QueryRowContext(repo.Context, colid, userid)
	if row.Err() == nil {
		if err := row.Scan(&f); err != nil {
			repo.Logger.Panicf("IsCollectionWriteble Cannot scan data: %v", err.Error())
			return false, err
		}
	} else {
		repo.Logger.Panicf("Cannot exec IsCollectionWriteble command: %v", row.Err().Error())
		return false, row.Err()
	}
	return f, nil
}
