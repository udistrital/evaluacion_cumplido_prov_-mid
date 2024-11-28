package models

import "time"

type ItemEvaluacion struct {
	Id                string
	EvaluacionId      int
	Idendificador     *Evaluacion
	Nombre            string
	ValorInitario     float64
	Iva               float64
	FichaTecnica      string
	Unidad            int
	Cantidad          float64
	TipoNecessidad    int
	Activo            bool
	FechaCreacion     time.Time
	FechaModificacion time.Time
}
