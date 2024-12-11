package models

type Clasificacion struct {
	Id                int
	Nombre            string
	CodigoAbreviacion string
	Descripcion       string
	LimiteInferior    int
	LimiteSuperior    int
	Activo            bool
}
