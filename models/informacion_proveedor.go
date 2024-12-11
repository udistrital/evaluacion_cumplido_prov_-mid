package models

type Estado struct {
	Id                   int
	ClaseParametro       string
	ValorParametro       string
	DescripcionParametro string
	Abreviatura          string
}

type InformacionProveedor struct {
	Id                      int
	Tipopersona             string
	NumDocumento            string
	IdCiudadContacto        float64
	Direccion               string
	Correo                  string
	Web                     string
	NomAsesor               string
	TelAsesor               string
	Descripcion             string
	PuntajeEvaluacion       float64
	ClasificacionEvaluacion string
	Estado                  *Estado
	TipoCuentaBancaria      string
	NumCuentaBancaria       string
	IdEntidadBancaria       float64
	FechaRegistro           string
	FechaUltimaModificacion string
	NomProveedor            string
	Anexorut                string
	Anexorup                string
	RegimenContributivo     string
}
