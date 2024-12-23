package models

import "time"

type Evaluacion struct {
	Id                 int
	ContratoSuscritoId int
	VigenciaContrato   int
	Activo             bool
	DocumentoId        int
	FechaCreacion      time.Time
	FechaModificacion  time.Time
}
