package routers

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context/param"
)

func init() {

    beego.GlobalControllerRouter["github.com/udistrital/evaluacion_cumplido_prov_mid/controllers:DocumentoEvaluacionController"] = append(beego.GlobalControllerRouter["github.com/udistrital/evaluacion_cumplido_prov_mid/controllers:DocumentoEvaluacionController"],
        beego.ControllerComments{
            Method: "GenerarDocumentoEvaluacion",
            Router: "/:evaluacion_id",
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/evaluacion_cumplido_prov_mid/controllers:ResultadoEvaluacionController"] = append(beego.GlobalControllerRouter["github.com/udistrital/evaluacion_cumplido_prov_mid/controllers:ResultadoEvaluacionController"],
        beego.ControllerComments{
            Method: "GuardarResultadoEvaluacion",
            Router: "/resultado-evaluacion",
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/evaluacion_cumplido_prov_mid/controllers:TestController"] = append(beego.GlobalControllerRouter["github.com/udistrital/evaluacion_cumplido_prov_mid/controllers:TestController"],
        beego.ControllerComments{
            Method: "Ping",
            Router: "/ping",
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

}
