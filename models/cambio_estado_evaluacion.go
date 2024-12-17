package models

import "time"

type CambioEstadoEvaluacion struct {
	Id                 int
	EvaluacionId       *Evaluacion
	EstadoEvaluacionId *EstadoEvaluacion
	Activo             bool
	FechaCreacion      time.Time
	FechaModificacion  time.Time
}
