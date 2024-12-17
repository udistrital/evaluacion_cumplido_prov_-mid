package models

type DocumentosSoporteSimplificado struct {
	SoporteCumplidoId int
	Documento         DocumentoSimplificado
	Archivo           FileGestorDocumental
}

type DocumentoSimplificado struct {
	Id                             int
	Nombre                         string
	TipoDocumento                  string
	CodigoAbreviacionTipoDocumento string
	Descripcion                    string
	Observaciones                  string
	FechaCreacion                  string
}
