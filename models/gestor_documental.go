package models

import "time"

type Thumbnail struct {
	Name            string  `json:"name"`
	MimeType        string  `json:"mime-type"`
	Encoding        *string `json:"encoding"`
	DigestAlgorithm string  `json:"digestAlgorithm"`
	Digest          string  `json:"digest"`
	Length          string  `json:"length"`
	Data            string  `json:"data"`
}

type FileContent struct {
	Name            string  `json:"name"`
	MimeType        string  `json:"mime-type"`
	Encoding        *string `json:"encoding"`
	DigestAlgorithm string  `json:"digestAlgorithm"`
	Digest          string  `json:"digest"`
	Length          string  `json:"length"`
	Data            string  `json:"data"`
}

type DocumentoEnlace struct {
	UIDUID               *string       `json:"uid:uid"`
	UIDMajorVersion      int           `json:"uid:major_version"`
	UIDMinorVersion      int           `json:"uid:minor_version"`
	Thumbnail            Thumbnail     `json:"thumb:thumbnail"`
	FileContent          FileContent   `json:"file:content"`
	CommonIconExpanded   *string       `json:"common:icon-expanded"`
	CommonIcon           string        `json:"common:icon"`
	FilesFiles           []interface{} `json:"files:files"`
	Description          string        `json:"dc:description"`
	Language             *string       `json:"dc:language"`
	Coverage             *string       `json:"dc:coverage"`
	Valid                *string       `json:"dc:valid"`
	Creator              string        `json:"dc:creator"`
	Modified             time.Time     `json:"dc:modified"`
	LastContributor      string        `json:"dc:lastContributor"`
	Rights               *string       `json:"dc:rights"`
	Expired              *string       `json:"dc:expired"`
	Format               *string       `json:"dc:format"`
	Created              time.Time     `json:"dc:created"`
	Title                string        `json:"dc:title"`
	Issued               *string       `json:"dc:issued"`
	Nature               *string       `json:"dc:nature"`
	Subjects             []interface{} `json:"dc:subjects"`
	Contributors         []string      `json:"dc:contributors"`
	Source               *string       `json:"dc:source"`
	Publisher            *string       `json:"dc:publisher"`
	RelatedTextResources []interface{} `json:"relatedtext:relatedtextresources"`
	Tags                 []interface{} `json:"nxtag:tags"`
	File                 string        `json:"file"`
}
