package models


type Bucket struct {
	Bucket    string    `json:"bucket"` // "urn:cts:sktlit:skt0001.nyaya002.J1D:"
	Citations []string `json:"citations"`
}

type Collection struct {
	Id       int      `json:"id"`
	Title    string   `json:"title"`
	IsPublic bool     `json:"is_public"`
	AuthorId int      `json:"author_id"`
	Buckets  []Bucket `json:"buckets"`
}

type CollectionPage struct {
	Collection           Collection `json:"collection"`
	SampleText           string     `json:"sample"`
	CollectionURL        string     `json:"collectionurl"`
	EditCollectionURL    string     `json:"editcollectionurl"`
	SharingCollectionURL string     `json:"sharingcollectionurl"`
	DropCollectionURL    string     `json:"dropcollectionurl"`
	CanEditCollection    bool       `json:"caneditcollection"`
	CanSharingCollection bool       `json:"cansharingcollection"`
	CanDropCollection    bool       `json:"candropcollection"`
	Author               *Author    `json:"author"`
}

type ColUsers struct {
	Id_Col   int  `json:"col_id"`
	Id_User  int  `json:"user_id"`
	Is_Write bool `json:"is_write"`
}

type ShareCollection struct {
	Collection `json:"collection"`
	ColUsers   []ColUsers `json:"colusers"`
}
