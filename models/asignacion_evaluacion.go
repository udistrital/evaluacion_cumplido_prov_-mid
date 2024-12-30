package models

type AsignacionEvaluacion struct {
	AsignacionEvaluacionId    int
	NombreProveedor           string
	RolEvaluador              string
	Dependencia               string
	TipoContrato              string
	NumeroContrato            string
	VigenciaContrato          string
	EvaluacionId              int
	EstadoAsignacionEvaluador *EstadoAsignacionEvaluador
	EstadoEvaluacion          *EstadoEvaluacion
}
