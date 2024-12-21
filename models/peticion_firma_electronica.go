package models

type PeticionFirmaElectronica struct {
	PersonaId    string `json:"PersonaId"`
	AsignacionId int    `json:"AsignacionId"`
}

type Firmante struct {
	Cargo          string `json:"cargo"`
	Identificacion string `json:"identificacion"`
	Nombre         string `json:"nombre"`
	TipoId         string `json:"tipoId"`
}
type Metadatos map[string]interface{}

type PeticionFirmaElectronicaCrud struct {
	Descripcion     string     `json:"descripcion"`
	File            string     `json:"file"`
	Firmantes       []Firmante `json:"firmantes"`
	EtapaFirma      int        `json:"etapa_firma"`
	IdTipoDocumento int        `json:"IdTipoDocumento"`
	Metadatos       Metadatos  `json:"metadatos"`
	Nombre          string     `json:"nombre"`
	Representantes  []Firmante `json:"representantes"`
}
