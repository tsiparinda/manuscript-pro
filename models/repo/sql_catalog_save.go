package repo

func (repo *SqlRepository) SaveCiteDataDict(id_col int, bucket string, catkey string, catvalue []byte) error {
	_, err := repo.Commands.SaveCiteDataDict.Exec(id_col, bucket, catkey, catvalue)
	if err != nil {
		repo.Logger.Panicf("SaveCatalog Cannot exec SaveCatalog command: %v", err.Error())
		return err
	}

	return nil
}
