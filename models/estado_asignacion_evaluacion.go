package models

type EstadoAsignacionEvaluador struct {
	Id                int    `json:"Id"`
	Nombre            string `json:"Nombre"`
	CodigoAbreviacion string `json:"CodigoAbreviacion"`
	Descripcion       string `json:"Descripcion"`
	Activo            bool   `json:"Activo"`
}
