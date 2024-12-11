package models

type InformacionDocumentoEvaluacion struct {
	EmpresaProveedor         string
	Documento                string
	ObjetoContrato           string
	Dependencia              string
	ResultadoFinalEvaluacion ResultadoFinalEvaluacion
}
