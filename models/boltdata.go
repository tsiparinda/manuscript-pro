package models

import (
	"brucheion/gocite"
)

// *** CITE Data Containers ***
//These are probably going to be retired altogether.
// BoltCatalog contains all metadata of a CITE URN and is
// used in handleCEXLoad and page functions
type BoltCatalog struct {
	ColId         int    `json:"colid"`          // 0
	URN           string `json:"urn"`            // urn:cts:sktlit:skt0001.nyaya002.J1D:
	Citation      string `json:"citationScheme"` // adhyāya.āhnika.sūtra
	GroupName     string `json:"groupName"`      // Nyāya
	WorkTitle     string `json:"workTitle"`      // Nyāyabhāṣya
	VersionLabel  string `json:"versionLabel"`   // J1D
	ExemplarLabel string `json:"exemplarLabel"`  //
	Online        string `json:"online"`         // true
	Language      string `json:"language"`       // san
}

// BoltData is the container for CITE data imported from CEX files and is used in handleCEXLoad
type BoltData struct {
	Bucket    []string // workurn
	Data      []gocite.Work
	Catalog   []BoltCatalog
	ID_author int
	Title     string
}

// BoltWork is the container for BultURNs and their associated keys and is used in handleCEXLoad
type BoltWork struct {
	Key  []string // cts-node urn
	Data []BoltURN
}

// BoltURN is the container for a textpassage along with its URN, its image reference,
// and some information on preceding and anteceding works.
// Used for loading and saving CEX files, for pages, and for nodes
type BoltURN struct {
	URN      string   `json:"urn"`
	Text     string   `json:"text"`
	LineText []string `json:"linetext"`
	Previous string   `json:"previous"`
	Next     string   `json:"next"`
	First    string   `json:"first"`
	Last     string   `json:"last"`
	Index    int      `json:"sequence"`
	ImageRef []string `json:"imageref"`
}

// BoltJSON is a string representation of a JSON used in BoltRetrieve
type BoltJSON struct {
	JSON string
}

// BucketDict is the container for key-value pair from one of buckets
type BucketDict struct {
	Key   string
	Value string
}
