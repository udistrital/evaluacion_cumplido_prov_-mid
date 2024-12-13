package models

type AsignacionEvaluacion struct {
	AsignacionEvaluacionId    int
	NombreProveedor           string
	Dependencia               string
	TipoContrato              string
	NumeroContrato            string
	VigenciaContrato          string
	EvaluacionId              int
	EstadoAsignacionEvauacion *EstadoAsignacionEvaluador
	EstadoEvaluacion          *EstadoEvaluacion
}
