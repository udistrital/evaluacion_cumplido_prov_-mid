package models

import (
	"time"
)

type CambioEstadoEvaluacion struct {
	Id                 int               `orm:"column(id);pk;auto"`
	EvaluacionId       *Evaluacion       `orm:"column(evaluacion_id);rel(fk)"`
	EstadoEvaluacionId *EstadoEvaluacion `orm:"column(estado_evaluacion_id);rel(fk)"`
	Activo             bool              `orm:"column(activo);default(true)"`
	FechaCreacion      time.Time         `orm:"auto_now_add;column(fecha_creacion);type(timestamp without time zone);null"`
	FechaModificacion  time.Time         `orm:"auto_now;column(fecha_modificacion);type(timestamp without time zone);null"`
}
