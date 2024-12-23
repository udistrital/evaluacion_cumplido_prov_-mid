package models

import "time"

type TipoDocumentoCrud struct {
	Id                   int
	ClaseParametro       string
	ValorParametro       string
	DescripcionParametro string
	Abreviatura          string
}

type Perfil struct {
	Id                   int
	ClaseParametro       string
	ValorParametro       string
	DescripcionParametro string
	Abreviatura          string
}

type Persona struct {
	TipoDocumento                     TipoDocumentoCrud
	Id                                string
	DigitoVerificacion                int
	PrimerApellido                    string
	SegundoApellido                   string
	PrimerNombre                      string
	SegundoNombre                     string
	Cargo                             string
	IdPaisNacimiento                  int
	Perfil                            Perfil
	Profesion                         string
	Especialidad                      string
	MontoCapitalAutorizado            int
	Genero                            string
	GrupoEtnico                       string
	ComunidadLgbt                     bool
	CabezaFamilia                     bool
	PersonasACargo                    bool
	NumeroPersonasACargo              int
	EstadoCivil                       string
	Discapacitado                     bool
	TipoDiscapacidad                  string
	DeclaranteRenta                   bool
	MedicinaPrepagada                 bool
	ValorUvtPrepagada                 int
	CuentaAhorroAfc                   bool
	NumCuentaBancariaAfc              string
	IdEntidadBancariaAfc              int
	InteresViviendaAfc                int
	DependienteHijoMenorEdad          bool
	DependienteHijoMenos23Estudiando  bool
	DependienteHijoMas23Discapacitado bool
	DependienteConyuge                bool
	DependientePadreOHermano          bool
	IdNucleoBasico                    int
	IdArl                             int
	IdEps                             int
	IdFondoPension                    int
	IdCajaCompensacion                int
	FechaExpedicionDocumento          time.Time
	IdCiudadExpedicionDocumento       int
}
