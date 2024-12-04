package models

type UnidadMedida struct {
	Id          int    `json:"Id"`
	Unidad      string `json:"Unidad"`
	Tipo        string `json:"Tipo"`
	Descripcion string `json:"Descripcion"`
	Estado      bool   `json:"Estado"`
}
