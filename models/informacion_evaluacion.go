package models

type InformacionEvaluacion struct {
	NombreEvaluador        string      `json:"NombreEvaluador"`
	Cargo                  string      `json:"Cargo"`
	CodigoAbreviacionRol   string      `json:"CodigoAbreviacionRol"`
	DependenciaEvaluadora  string      `json:"DependenciaEvaluadora"`
	FechaEvaluacion        string      `json:"FechaEvaluacion"`
	EmpresaProveedor       string      `json:"EmpresaProveedor"`
	ObjetoContrato         string      `json:"ObjetoContrato"`
	PuntajeTotalEvaluacion int         `json:"PuntajeTotalEvaluacion"`
	Clasificacion          string      `json:"Clasificacion"`
	ItemsEvaluados         []Item      `json:"ItemsEvaluados"`
	Evaluadores            []Evaluador `json:"Evaluadores"`
	ResultadoEvaluacion    Resultado   `json:"ResultadoEvaluacion"`
}

type Evaluador struct {
	Documento         string `json:"Documento"`
	Cargo             string `json:"Cargo"`
	Rol               string `json:"Rol"`
	ItemsEvaluados    string `json:"ItemEvaluado"`
	PuntajeEvaluacion int    `json:"PuntajeEvaluacion"`
	EstadoEvaluacion  string `json:"EstadoEvaluacion"`
	Observaciones     string `json:"Observaciones"`
}
