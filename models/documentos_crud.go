package models

import "time"

type DominioTipoDocumento struct {
	Id                int     `json:"Id"`
	Nombre            string  `json:"Nombre"`
	Descripcion       string  `json:"Descripcion"`
	CodigoAbreviacion string  `json:"CodigoAbreviacion"`
	Activo            bool    `json:"Activo"`
	NumeroOrden       float64 `json:"NumeroOrden"`
}

type TipoDocumento struct {
	Id                   int                  `json:"Id"`
	Nombre               string               `json:"Nombre"`
	Descripcion          string               `json:"Descripcion"`
	CodigoAbreviacion    string               `json:"CodigoAbreviacion"`
	Activo               bool                 `json:"Activo"`
	NumeroOrden          float64              `json:"NumeroOrden"`
	Tamano               int                  `json:"Tamano"`
	Extension            string               `json:"Extension"`
	Workspace            string               `json:"Workspace"`
	TipoDocumentoNuxeo   string               `json:"TipoDocumentoNuxeo"`
	DominioTipoDocumento DominioTipoDocumento `json:"DominioTipoDocumento"`
}

type DocumentoCrud struct {
	Id                int           `json:"Id"`
	Nombre            string        `json:"Nombre"`
	Descripcion       string        `json:"Descripcion"`
	Enlace            string        `json:"Enlace"`
	TipoDocumento     TipoDocumento `json:"TipoDocumento"`
	Metadatos         string        `json:"Metadatos"`
	Activo            bool          `json:"Activo"`
	FechaCreacion     time.Time     `json:"FechaCreacion"`
	FechaModificacion time.Time     `json:"FechaModificacion"`
}
