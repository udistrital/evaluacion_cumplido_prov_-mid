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
		beego.NSNamespace("/resultado_final-evaluacion",
			beego.NSInclude(
				&controllers.DocumentoEvaluacionController{},
			),
		),
		beego.NSNamespace("/informacion_evaluacion",
			beego.NSInclude(
				&controllers.InformacionEvaluacionController{},
			),
		),
		beego.NSNamespace("/evaluacion_cumplido",
			beego.NSInclude(
				&controllers.EvaluacionCumplidoController{},
			),
		),
		beego.NSNamespace("/asignaciones",
			beego.NSInclude(
				&controllers.AsignacionesController{},
			),
		),
		beego.NSNamespace("/cambio_rol_evaluador",
			beego.NSInclude(
				&controllers.CambioRolEvaluadorController{},
			),
		),
		beego.NSNamespace("/evaluacion",
			beego.NSInclude(
				&controllers.EvaluacionController{},
			),
		),
		beego.NSNamespace("/firma_electronica",
			beego.NSInclude(
				&controllers.FirmaElectronica{},
			),
		))

	beego.AddNamespace(ns)
}
