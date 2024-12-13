package models

type AsignacionEvaluacion struct {
	AsignacionEvaluacionId    int    `json:"AsignacionEvaluacionId"`
	NombreProveedor           string `json:"NombreProveedor"`
	Dependencia               string `json:"Depenedencia"`
	TipoContrato              string `json:"TipoContrato"`
	NumeroContrato            string `json:"NumeroContrato"`
	VigenciaContrato          string `json:"VigenciaContrato"`
	EvaluacionId              int    `json:"EvaluacionId"`
	EstadoAsignacionEvauacion *EstadoAsignacionEvaluador
	EstadoEvaluacion          *EstadoEvaluacion
}
