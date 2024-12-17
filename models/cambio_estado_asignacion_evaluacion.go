package models

import "time"

type CambioEstadoASignacionEnvaluacion struct {
	Id                        int                        `json:"Id"`
	EstadoAsignacionEvaluador *EstadoAsignacionEvaluador `json:"EstadoAsignacionEvaluadorId"`
	AsignacionEvaluadorId     AsignacionEvaluador        `json:"AsignacionEvaluadorId"`
	Activo                    bool                       `json:"Activo"`
	FechaCreacion             time.Time                  `json:"FechaCreacion"`
	FechaModificacion         time.Time                  `json:"FechaModificacion"`
}
