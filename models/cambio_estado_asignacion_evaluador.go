package models

import "time"

type CambioEstadoAsignacionEvaluador struct {
	Id                          int
	EstadoAsignacionEvaluadorId EstadoAsignacionEvaluador
	AsignacionEvaluadorId       AsignacionEvaluador
	Activo                      bool
	FechaCreacion               time.Time
	FechaModificacion           int16
}
