package models

import "time"

type AsignacionEvaluador struct {
	Id                       int
	EvaluacionId             Evaluacion
	PersonaId                int
	Cargo                    string
	PorcentajeEvaluacion     float64
	RolAsignacionEvaluadorId RolAsignacionEvaluador
	Activo                   bool
	FechaCreacion            time.Time
	FechaModificacion        time.Time
}
