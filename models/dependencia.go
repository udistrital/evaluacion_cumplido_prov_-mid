package models

type Dependencia struct {
	Codigo string
	Nombre string
}

type DependenciaLista struct {
	Dependencias []Dependencia
}

type DependenciasRespuesta struct {
	Dependencias struct {
		Dependencia []Dependencia `json:"dependencia"`
	} `json:"dependencias"`
}
