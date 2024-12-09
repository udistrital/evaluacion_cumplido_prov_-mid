package models

type Contrato struct {
	Vigencia       string `json:"vigencia"`
	NumeroContrato string `json:"numero_contrato"`
}

type ContratosRespuesta struct {
	Contratos struct {
		Contrato []Contrato `json:"contrato"`
	} `json:"contratos"`
}
