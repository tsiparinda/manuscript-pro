package utils

import (
	"brucheion/models"
	"bytes"
	"encoding/gob"
)

// gobDecodeImgCol decodes a byte slice from the database to an imageCollection
func GobDecodeImgCol(data []byte) (models.ImageCollection, error) {
	var p *models.ImageCollection
	buf := bytes.NewBuffer(data)
	dec := gob.NewDecoder(buf)
	err := dec.Decode(&p)
	if err != nil {
		return models.ImageCollection{}, err
	}
	return *p, nil
}
