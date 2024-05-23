// all methods in this package scan rows from a query and return different data type responce  
package repo

import (
	"brucheion/models"
	"database/sql"
	"errors"
)


func scanBoltCatalog(rows *sql.Rows) (catalog []models.BoltCatalog, err error) {
	catalog = make([]models.BoltCatalog, 0, 10)
	for rows.Next() {
		p := models.BoltCatalog{}
		err = rows.Scan(&p.URN, &p.Citation, &p.GroupName, &p.WorkTitle,
			&p.VersionLabel, &p.ExemplarLabel, &p.Online, &p.Language)
		if err == nil {
			catalog = append(catalog, p)
		} else {
			return
		}
	}
	return
}

func scanUser(rows *sql.Rows) (cred models.Credentials, err error) {
	//user = make([]identity.User, 0, 10)
	for rows.Next() {
		//var puserline = models.Credentials{} models.Credentials
		err = rows.Scan(&cred.Id, &cred.Username, &cred.Email, &cred.Password, &cred.IsVerified, &cred.VerificationCode)
		// if err == nil {
		// 	cred = models.Credentials{p.Id, p.Username, p.Email, p.Password, p.IsVerified, p.VerificationCode}
		// 	return
		// }
	}
	return
}

func scanPassage(rows *sql.Rows) (catalog []models.Passage, err error) {
	catalog = make([]models.Passage, 0, 10)
	for rows.Next() {
		p := models.Passage{}
		// !!!!!!!!!!!!!!!!!!!!!!!!!!!! for del TextRefs
		err = rows.Scan(&p.Catalog, &p.FirstPassage, &p.PassageID, &p.ImageRefs,
			&p.LastPassage, &p.NextPassage, &p.PreviousPassage, &p.TextRefs, &p.Transcriber, &p.TranscriptionLines)
		if err == nil {
			catalog = append(catalog, p)
		} else {
			return
		}
	}
	return
}

func scanDict(rows *sql.Rows) (values []models.BucketDict, err error) {
	values = make([]models.BucketDict, 0, 10)
	for rows.Next() {
		var p models.BucketDict
		err = rows.Scan(&p.Key, &p.Value)
		if err == nil {
			values = append(values, p)
		} else {
			return
		}
	}
	return
}

func scanStrings(rows *sql.Rows) (values []string, err error) {
	values = make([]string, 0, 10)
	for rows.Next() {
		var p string
		err = rows.Scan(&p)
		if err == nil {
			values = append(values, p)
		} else {
			return
		}
	}
	return
}

func scanBucketKeyValue(rows *sql.Rows) (value string, err error) {
	if rows == nil {
		err = errors.New("scanBucketKeyValue: this key not found")
	}
	for rows.Next() {
		var p string
		err = rows.Scan(&p)
		if err == nil {
			value = p
		} else {
			return
		}
	}
	return
}

func scanBucketKeys(rows *sql.Rows) (value []string, err error) {
	if rows == nil {
		err = errors.New("scanBucketKeys: this key not found")
	}
	for rows.Next() {
		var p string
		err = rows.Scan(&p)
		if err == nil {
			value = append(value, p)
		} else {
			return
		}
	}
	return
}

func scanCollectionsPage(rows *sql.Rows) (collections []models.CollectionPage, err error) {
	collections = make([]models.CollectionPage, 0, 10)
	for rows.Next() {
		p := models.CollectionPage{Author: &models.Author{}}
		err = rows.Scan(&p.Collection.Id, &p.Collection.Title, &p.Collection.IsPublic, &p.Collection.AuthorId,
			&p.Author.Id, &p.Author.Name, &p.CanEditCollection)
		if err == nil {
			collections = append(collections, p)
		} else {
			return
		}
	}
	return
}

func scanCollection(rows *sql.Rows) (p models.Collection, err error) {
	// p = models.Collection{Author: &models.Author{}}
	p = models.Collection{}
	for rows.Next() {
		err = rows.Scan(&p.Id, &p.Title, &p.IsPublic, &p.AuthorId)
		// err = rows.Scan(&p.Id, &p.Title, &p.IsPublic, &p.SampleText, &p.Author.Id,
		// 	&p.Author.Name)
	}
	return p, err
}

func scanColUsers(rows *sql.Rows) (ret []models.ColUsers, err error) {
	var p models.ColUsers
	for rows.Next() {
		err = rows.Scan(&p.Id_Col, &p.Id_User, &p.Is_Write)
		if err != nil {
			return []models.ColUsers{}, err
		}
		ret = append(ret, p)
	}
	return ret, err
}

func scanImageKeyValue(rows *sql.Rows) (value string, err error) {
	if rows == nil {
		err = errors.New("scanImageKeyValue: this key not found")
	}
	for rows.Next() {
		var p string
		err = rows.Scan(&p)
		if err == nil {
			value = p
		} else {
			return
		}
	}
	return
}
