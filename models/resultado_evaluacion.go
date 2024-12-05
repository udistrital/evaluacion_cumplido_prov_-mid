package models

type ResultadoEvaluacion struct {
	Id                    int
	AsignacionEvaluadorId int
	ClasificacionId       int
	ResultadoEvaluacion   string
	Activo                bool
	FechaCreacion         string
	FechaModificacion     string
}
