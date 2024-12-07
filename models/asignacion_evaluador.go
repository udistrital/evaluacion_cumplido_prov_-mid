package models

import (
	"time"
)

type AsignacionEvaluador struct {
	Id                       int                     `orm:"column(id);pk;auto"`
	EvaluacionId             *Evaluacion             `orm:"column(evaluacion_id);rel(fk)"`
	PersonaId                int                     `orm:"column(persona_id)"`
	Cargo                    string                  `orm:"column(cargo)"`
	PorcentajeEvaluacion     float64                 `orm:"column(porcentaje_evaluacion)"`
	RolAsignacionEvaluadorId *RolAsignacionEvaluador `orm:"column(rol_asignacion_evaluador_id);rel(fk)"`
	Activo                   bool                    `orm:"column(activo);default(true)"`
	FechaCreacion            time.Time               `orm:"auto_now_add;column(fecha_creacion);type(timestamp without time zone);null"`
	FechaModificacion        time.Time               `orm:"auto_now;column(fecha_modificacion);type(timestamp without time zone);null"`
}
