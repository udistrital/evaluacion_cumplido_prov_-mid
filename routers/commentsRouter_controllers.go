package routers

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context/param"
)

func init() {

	beego.GlobalControllerRouter["github.com/udistrital/evaluacion_cumplido_prov_mid/controllers:CargaDataExcelController"] = append(beego.GlobalControllerRouter["github.com/udistrital/evaluacion_cumplido_prov_mid/controllers:CargaDataExcelController"],
		beego.ControllerComments{
			Method:           "UploadExcel",
			Router:           "/upload",
			AllowHTTPMethods: []string{"post"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["github.com/udistrital/evaluacion_cumplido_prov_mid/controllers:ResultadoEvaluacionController"] = append(beego.GlobalControllerRouter["github.com/udistrital/evaluacion_cumplido_prov_mid/controllers:ResultadoEvaluacionController"],
		beego.ControllerComments{
			Method:           "GuardarResultadoEvaluacion",
			Router:           "/resultado-evaluacion",
			AllowHTTPMethods: []string{"post"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})
}
