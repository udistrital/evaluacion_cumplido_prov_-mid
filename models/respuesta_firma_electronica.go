package models

type CuerpoFimaElectronica struct {
	Id            int           `json:"Id"`
	Nombre        string        `json:"Nombre"`
	Descripcion   string        `json:"Descripcion"`
	Enlace        string        `json:"Enlace"`
	TipoDocumento TipoDocumento `json:"TipoDocumento"`
	Metadatos     interface{}   `json:"Metadatos"`
	Activo        bool          `json:"Activo"`
}

type RespuestaFirmaElectronica struct {
	Status string                `json:"Status"`
	Res    CuerpoFimaElectronica `json:"res"`
}
