package models

type PeticionCambioEstadoEvaluacion struct {
	EvaluacionId      *Evaluacion `json:"EvaluacionId"`
	AbreviacionEstado string      `json:"AbreviacionEstado"`
}

type PeticionCambioEstadoAsignacion struct {
	AsignacionId      *AsignacionEvaluador `json:"AsignacionId"`
	AbreviacionEstado string               `json:"AbreviacionEstado"`
}
