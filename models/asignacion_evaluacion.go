package models

type AsignacionEvaluacion struct {
	NombreProveedor  string `json:"NombreProveedor"`
	Dependencia      string `json:"Depenedencia"`
	TipoContrato     string `json:"TipoContrato"`
	NumeroContrato   string `json:"NumeroContrato"`
	VigenciaContrato string `json:"VigenciaContrato"`
	EvaluacionId     int    `json:"EvaluacionId"`
	Estado           bool   `json:"Estado"`
}
