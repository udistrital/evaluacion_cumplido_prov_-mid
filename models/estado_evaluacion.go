package models

type EstadoEvaluacion struct {
	Id                int    `orm:"column(id);pk;auto"`
	Nombre            string `orm:"column(nombre)"`
	CodigoAbreviacion string `orm:"column(codigo_abreviacion)"`
	Descripcion       string `orm:"column(descripcion);null"`
	Activo            bool   `orm:"column(activo);default(true)"`
}
