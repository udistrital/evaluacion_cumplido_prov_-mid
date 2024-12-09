package models

type Dependencia struct {
	Codigo string `json:"codigo"`
	Nombre string `json:"nombre"`
}

type DependenciaLista struct {
	Dependencias []Dependencia `json:"dependencia"`
}

type DependenciasRespuesta struct {
	Dependencias struct {
		Dependencia []Dependencia `json:"dependencia"`
	} `json:"dependencias"`
}
