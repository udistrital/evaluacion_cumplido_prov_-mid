package models

import "time"

type AsignacionEvaluadorItem struct {
	Id                       int
	AsignacionEvaluadorId    AsignacionEvaluador
	ItemId                   Item
	Activo                   bool
	RolAsignacionEvaluadorId RolAsignacionEvaluador
	FechaCreacion            time.Time
	FechaModificacion        time.Time
}
