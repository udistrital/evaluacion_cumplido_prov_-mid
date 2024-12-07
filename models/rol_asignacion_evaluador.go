package models

type RolAsignacionEvaluador struct {
	Id                int    `orm:"column(id);pk"`
	Nombre            string `orm:"column(nombre);null"`
	Descripcion       string `orm:"column(descripcion);null"`
	CodigoAbreviacion string `orm:"column(codigo_abreviacion);null"`
	Activo            bool   `orm:"column(activo);null"`
}
