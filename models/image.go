package models

// image is the container for image metadata
type Image struct {
	ColId    int    `json:"colid"`     // 49
	URN      string `json:"imagename"` // urn:cite2:iiifimages:test:1
	Name     string `json:"colname"`   // urn:cite2:iiifimages:test:
	Protocol string `json:"protocol"`  // iiif
	License  string `json:"license"`   // CC-BY-4.0
	External bool   `json:"external"`  // false/true
	Location string `json:"location"`  // https://libimages1.princeton.edu/loris/pudl0001%2F4609321%2Fs42%2F00000004.jp2/info.json
}

// imageCollection is the container for image collections along with their URN and name as strings
type ImageCollection struct {
	URN        string  `json:"urn"`
	Name       string  `json:"name"`
	Collection []Image `json:"images"`
}
