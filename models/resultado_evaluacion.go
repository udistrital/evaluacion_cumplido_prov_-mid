package models

import "time"

type ResultadoEvaluacion struct {
	Id                    int
	AsignacionEvaluadorId AsignacionEvaluador
	ClasificacionId       Clasificacion
	ResultadoEvaluacion   string
	Observaciones         string
	Activo                bool
	FechaCreacion         time.Time
	FechaModificacion     time.Time
}

type ResultadoFinalEvaluacion struct {
	Evaluadores []struct {
		Nombre string
		Cargo  string
		Items  string
		Rol    string
	}
	Resultados []struct {
		Categoria     string
		Titulo        string
		Pregunta      string
		Cumplimiento  string
		ValorAsignado int
	}
}
