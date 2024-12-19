package routers

import (
	"github.com/astaxie/beego"
	"github.com/udistrital/evaluacion_cumplido_prov_mid/controllers"
	"github.com/udistrital/utils_oas/errorhandler"
)

func init() {
	beego.ErrorController(&errorhandler.ErrorHandlerController{})
	ns := beego.NewNamespace("/v1",
		beego.NSNamespace("/carga-data-excel",
			beego.NSInclude(
				&controllers.CargaDataExcelController{},
			),
		),
		beego.NSNamespace("/resultado",
			beego.NSInclude(
				&controllers.ResultadoEvaluacionController{},
			),
		),
		beego.NSNamespace("/resultado-final-evaluacion",
			beego.NSInclude(
				&controllers.DocumentoEvaluacionController{},
			),
		),
		beego.NSNamespace("/informacion-evaluacion",
			beego.NSInclude(
				&controllers.InformacionEvaluacionController{},
			),
		),
		beego.NSNamespace("/evaluacion-cumplido",
			beego.NSInclude(
				&controllers.EvaluacionCumplidoController{},
			),
		),
		beego.NSNamespace("/asignaciones",
			beego.NSInclude(
				&controllers.AsignacionesController{},
			),
		),
		beego.NSNamespace("/cambio-rol-evaluador",
			beego.NSInclude(
				&controllers.CambioRolEvaluadorController{},
			),
		),
		beego.NSNamespace("/evaluacion",
			beego.NSInclude(
				&controllers.EvaluacionController{},
			),
		))

	beego.AddNamespace(ns)
}
