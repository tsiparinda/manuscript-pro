package models

import "html/template"

// Transcription is the container for a transcription and context metadata
// Transcription is the container for a transcription and context metadata
type Transcription struct {
	ColID         int
	CTSURN        string
	Transcriber   string
	Transcription string
	Previous      string
	Next          string
	First         string
	Last          string
	ImageRef      []string
	TextRef       []string
	ImageJS       string
	CatID         string
	CatCit        string
	CatGroup      string
	CatWork       string
	CatVers       string
	CatExmpl      string
	CatOn         string
	CatLan        string
}

type TranscriptionPage struct {
	ColID        int
	User         string
	Title        string
	Text         template.HTML
	Previous     string
	PreviousLink template.HTML
	Next         string
	NextLink     template.HTML
	First        string
	Last         string
	ImageScript  template.HTML
	ImageHTML    template.HTML
	TextHTML     template.HTML
	ImageRef     string
	CatID        string
	CatCit       string
	CatGroup     string
	CatWork      string
	CatVers      string
	CatExmpl     string
	CatOn        string
	CatLan       string
	ImageJS      string
}
