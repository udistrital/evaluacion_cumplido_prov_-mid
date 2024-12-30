package models

import "time"

type ItemEvaluacion struct {
	EvaluacionId      Evaluacion
	Identificador     string
	Nombre            string
	ValorUnitario     float64
	Iva               float64
	FichaTecnica      string
	Unidad            int
	Cantidad          float64
	TipoNecesidad     int
	Activo            bool
	FechaCreacion     time.Time
	FechaModificacion time.Time
}
