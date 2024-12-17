package models

type DocumentoCrud struct {
	Id                int
	Nombre            string
	Descripcion       string
	Enlace            string
	TipoDocumento     *TipoDocumento
	Metadatos         string
	Activo            bool
	FechaCreacion     string
	FechaModificacion string
}

type TipoDocumento struct {
	Id                   int
	Nombre               string
	Descripcion          string
	CodigoAbreviacion    string
	Activo               bool
	NumeroOrden          float64
	Tamano               float64
	Extension            string
	Workspace            string
	TipoDocumentoNuxeo   string
	FechaCreacion        string
	FechaModificacion    string
	DominioTipoDocumento *DominioTipoDocumento
}

type DominioTipoDocumento struct {
	Id                int
	Nombre            string
	Descripcion       string
	CodigoAbreviacion string
	Activo            bool
	NumeroOrden       float64
	FechaCreacion     string
	FechaModificacion string
}
