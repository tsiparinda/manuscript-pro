package repo

import (
	"brucheion/models"
	"database/sql"
	"encoding/json"

	"github.com/vedicsociety/platform/http/handling/params"
)

func (repo *SqlRepository) AddNewCollectionTX(tx *sql.Tx, authorid int, title string) (int, error) {
	colid := 0
	stmt := tx.StmtContext(repo.Context, repo.Commands.AddNewCollection)
	err := stmt.QueryRow(authorid, title).Scan(&colid)
	//err := repo.Commands.AddNewCollection.QueryRow(authorid, title).Scan(&colid)
	if err != nil {
		repo.Logger.Debugf("Cannot exec AddNewCollectionTX command: %v", err.Error())
		return colid, err
	}
	repo.Logger.Debugf("AddNewCollectionTX New collection added with ID: %v", colid)
	return colid, nil
}

func (repo *SqlRepository) CreateBucketIfNotExists(tx *sql.Tx, bucket string, id_col int) error {

	//result, err := repo.Commands.CreateBucketIfNotExists.ExecContext(repo.Context, id_col, bucket)
	result, err := tx.StmtContext(repo.Context, repo.Commands.CreateBucketIfNotExists).Exec(id_col, bucket)
	if err != nil {
		repo.Logger.Panicf("Cannot exec CreateBucketIfNotExists command: %v", err.Error())
		return err
	}
	if err == nil {
		// we can use this result's methods:
		// LastInsertId() (int64, error)
		// RowsAffected() (int64, error)
		rows, _ := result.RowsAffected()
		repo.Logger.Debugf("CreateBucketIfNotExists.Count rendering rows:", rows)
		return nil
	} else {
		repo.Logger.Debugf("CreateBucketIfNotExists Cannot get inserted data", err.Error())
		return err
	}
}

func (repo *SqlRepository) AddNewCollection(params params.InputParams, userid int) (int, error) {
	// rr, _ := repo.DB.BeginTx(repo.Context, nil)
	tx, err := repo.DB.Begin()
	if err != nil {
		repo.Logger.Panicf("Cannot create transaction: %v", err.Error())
		return 0, err
	}
	id_col, err := repo.AddNewCollectionTX(tx, userid, params.InputParam[0].Value[0])
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	err = tx.Commit()
	if err != nil {
		repo.Logger.Panicf("Transaction cannot be committed: %v", err.Error())
		err = tx.Rollback()
		if err != nil {
			repo.Logger.Panicf("Transaction cannot be rolled back: %v", err.Error())
		}
		return 0, err
	}

	repo.Logger.Debugf("AddNewCollection New collection added with ID: %v", id_col)
	return id_col, nil
}

// runs from cexupload_handler.go
// importcex
func (repo *SqlRepository) SaveBoltData(p *models.BoltData) (int, error) {
	// {Bucket[]str, Data[]gocite.Work, Catalog[]BoltCatalog} ->
	//userID := repo.Context.Value("USER_SESSION_KEY").(int)

	repo.Logger.Debugf("SaveBoltData: userID= ", p.ID_author)
	//repo.Logger.Debugf("SaveBoltData: userID= ", repo.User.GetID())

	// rr, _ := repo.DB.BeginTx(repo.Context, nil)
	tx, err := repo.DB.Begin()
	if err != nil {
		repo.Logger.Panicf("Cannot create transaction: %v", err.Error())
		return 0, err
	}
	id_col, err := repo.AddNewCollectionTX(tx, p.ID_author, p.Title)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	//id_col := 1
	repo.Logger.Debugf("SaveBoltData: id_collection= ", id_col)

	repo.Logger.Info("SaveBoltData are starting...")
	for i := range p.Bucket {
		newbucket := p.Bucket[i]
		catkey := p.Bucket[i]
		catvalue, _ := json.Marshal(p.Catalog[i])

		repo.CreateBucketIfNotExists(tx, newbucket, id_col) //id_collection
		if err != nil {
			tx.Rollback()
			return 0, err
		}
		//put bucket data (hstore)
		//repo.Logger.Debugf("SaveBoltData: id_col, newbucket, catkey, catvalue", id_col, newbucket, catkey, catvalue)
		_, err := tx.StmtContext(repo.Context, repo.Commands.SaveCiteDataDict).Exec(id_col, newbucket, catkey, catvalue)
		if err != nil {
			repo.Logger.Debugf("SaveBoltData: Cannot get inserted data1", err.Error())
			tx.Rollback()
			return 0, err
		}

		for j := range p.Data[i].Passages {
			newkey := p.Data[i].Passages[j].PassageID
			newvalue, _ := json.Marshal(p.Data[i].Passages[j])

			_, err := tx.StmtContext(repo.Context, repo.Commands.SaveCiteDataDict).Exec(id_col, newbucket, newkey, newvalue)
			if err != nil {
				repo.Logger.Debugf("SaveBoltData: Cannot get inserted data2", err.Error())
				tx.Rollback()
				return 0, err
			}
		}
	}
	err = tx.Commit()
	if err != nil {
		repo.Logger.Panicf("Transaction cannot be committed: %v", err.Error())
		err = tx.Rollback()
		if err != nil {
			repo.Logger.Panicf("Transaction cannot be rolled back: %v", err.Error())
		}
	}
	return id_col, nil
}
