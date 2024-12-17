package models

import "time"

type DependenciaSic struct {
	ESFIDESPACIO    string    `json:"ESFIDESPACIO"`
	ESFCODIGODEP    string    `json:"ESFCODIGODEP"`
	ESFDEPENCARGADA string    `json:"ESFDEPENCARGADA"`
	ESFESTADO       string    `json:"ESFESTADO"`
	Ciudad          string    `json:"Ciudad"`
	Localidad       string    `json:"Localidad"`
	Barrio          string    `json:"Barrio"`
	Direccion       string    `json:"Direccion"`
	EstadoRegistro  bool      `json:"EstadoRegistro"`
	FechaRegistro   time.Time `json:"FechaRegistro"`
	Id              int       `json:"Id"`
}
