package models

type BodyResultadoEvaluacion struct {
	AsignacionEvaluadorId int       `json:"AsignacionEvaluadorId"`
	ClasificacionId       int       `json:"ClasificacionId"`
	ResultadoEvaluacion   Resultado `json:"ResultadoEvaluacion"`
	Observaciones         string    `json:"Observaciones"`
}

type Resultado struct {
	ResultadosIndividuales []struct {
		Categoria string `json:"Categoria"`
		Titulo    string `json:"Titulo"`
		Respuesta struct {
			Pregunta      string `json:"Pregunta"`
			Cumplimiento  string `json:"Cumplimiento"`
			ValorAsignado int    `json:"ValorAsignado"`
		} `json:"Respuesta"`
	} `json:"ResultadosIndividuales"`
}
