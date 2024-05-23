package models

//passage is a the container for passage metadata

type Passage struct {
	ColId              int         `json:"colid"`
	PassageID          string      `json:"passageid"`          // urn:cts:sktlit:skt0001.nyaya002.J1D:3.1.1
	Transcriber        string      `json:"transcriber"`        // user?
	TranscriptionLines []string    `json:"transcriptionLines"` // ["{J1D_37r4}parīkṣitāni" "pramā ṇāni prameyam idānīṃ parīkṣyate |" ...]
	PreviousPassage    string      `json:"previousPassage"`    // null or prev pass
	NextPassage        string      `json:"nextPassage"`        // urn:cts:sktlit:skt0001.nyaya002.J1D:3.1.2 or null
	FirstPassage       string      `json:"firstPassage"`       // urn:cts:sktlit:skt0001.nyaya002.J1D:3.1.1
	LastPassage        string      `json:"lastPassage"`        // urn:cts:sktlit:skt0001.nyaya002.J1D:3.2.72
	ImageRefs          []string    `json:"imageRefs"`          // [urn:cite2:nbh:J1img.negative:J1_37r urn:cite2:nbh:J1img.positive:J1_37v urn:cite2:nbh:J1img.negative:J1_37v ...]
	TextRefs           []string    `json:"textRefs"`           // !!!!!!!!!!!!!!!! for del!!!!!! list all of buckets [urn:cts:sktlit:skt0001.nyaya002.A3D:, urn:cts:sktlit:skt0001.nyaya002.M3D: ...]
	Schemes            []string    `json:"schemes"`            // list of schems of current short urn (3.1.1 3.1.2 3.1.3...)
	Catalog            BoltCatalog `json:"catalog"`            // {0 urn:cts:sktlit:skt0001.nyaya002.J1D: adhyāya.āhnika.sūtra Nyāya Nyāyabhāṣya J1D  true san}
	Text               string      `json:"text"`               // "{J1D_37r4}parīkṣitāni \r\n  pramā ṇāni prameyam idānīṃ parīkṣyate  ..."
}

// returned from transcription's BlockEditor
type PassageText struct {
	ColId     int    `json:"colid"`
	PassageID string `json:"passageid"`
	Text      string `json:"text"`
}
