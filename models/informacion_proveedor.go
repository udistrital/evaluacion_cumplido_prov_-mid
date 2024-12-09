package models

type Estado struct {
	Id                   int    `json:"Id"`
	ClaseParametro       string `json:"ClaseParametro"`
	ValorParametro       string `json:"ValorParametro"`
	DescripcionParametro string `json:"DescripcionParametro"`
	Abreviatura          string `json:"Abreviatura"`
}

type Proveedor struct {
	Id                      int    `json:"Id"`
	Tipopersona             string `json:"Tipopersona"`
	NumDocumento            string `json:"NumDocumento"`
	IdCiudadContacto        int    `json:"IdCiudadContacto"`
	Direccion               string `json:"Direccion"`
	Correo                  string `json:"Correo"`
	Web                     string `json:"Web"`
	NomAsesor               string `json:"NomAsesor"`
	TelAsesor               string `json:"TelAsesor"`
	Descripcion             string `json:"Descripcion"`
	PuntajeEvaluacion       int    `json:"PuntajeEvaluacion"`
	ClasificacionEvaluacion string `json:"ClasificacionEvaluacion"`
	Estado                  Estado `json:"Estado"`
	TipoCuentaBancaria      string `json:"TipoCuentaBancaria"`
	NumCuentaBancaria       string `json:"NumCuentaBancaria"`
	IdEntidadBancaria       int    `json:"IdEntidadBancaria"`
	FechaRegistro           string `json:"FechaRegistro"`
	FechaUltimaModificacion string `json:"FechaUltimaModificacion"`
	NomProveedor            string `json:"NomProveedor"`
	Anexorut                string `json:"Anexorut"`
	Anexorup                string `json:"Anexorup"`
	RegimenContributivo     string `json:"RegimenContributivo"`
}
