package models

type UnidadMedida struct {
	Id          int    `json:"id"`
	Unidad      string `json:"unidad"`
	Tipo        string `json:"tipo"`
	Descripcion string `json:"descripcion"`
	Estado      bool   `json:"estado"`
}
