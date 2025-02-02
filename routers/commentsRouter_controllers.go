package routers

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context/param"
)

func init() {

    beego.GlobalControllerRouter["github.com/udistrital/evaluacion_cumplido_prov_mid/controllers:AsignacionesController"] = append(beego.GlobalControllerRouter["github.com/udistrital/evaluacion_cumplido_prov_mid/controllers:AsignacionesController"],
        beego.ControllerComments{
            Method: "CambiarEstadoAsignacionEvaluacion",
            Router: "/cambiar-estado",
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/evaluacion_cumplido_prov_mid/controllers:AsignacionesController"] = append(beego.GlobalControllerRouter["github.com/udistrital/evaluacion_cumplido_prov_mid/controllers:AsignacionesController"],
        beego.ControllerComments{
            Method: "ConsultarAsignaciones",
            Router: "/consultar/:numeroDocumento",
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/evaluacion_cumplido_prov_mid/controllers:CambioRolEvaluadorController"] = append(beego.GlobalControllerRouter["github.com/udistrital/evaluacion_cumplido_prov_mid/controllers:CambioRolEvaluadorController"],
        beego.ControllerComments{
            Method: "CambiarRolAsignacionEvaluador",
            Router: "/:idEvaluacion",
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/evaluacion_cumplido_prov_mid/controllers:CargaDataExcelController"] = append(beego.GlobalControllerRouter["github.com/udistrital/evaluacion_cumplido_prov_mid/controllers:CargaDataExcelController"],
        beego.ControllerComments{
            Method: "UploadExcel",
            Router: "/upload",
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/evaluacion_cumplido_prov_mid/controllers:DocumentoEvaluacionController"] = append(beego.GlobalControllerRouter["github.com/udistrital/evaluacion_cumplido_prov_mid/controllers:DocumentoEvaluacionController"],
        beego.ControllerComments{
            Method: "GenerarDocumentoEvaluacion",
            Router: "/:evaluacion_id",
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/evaluacion_cumplido_prov_mid/controllers:EvaluacionController"] = append(beego.GlobalControllerRouter["github.com/udistrital/evaluacion_cumplido_prov_mid/controllers:EvaluacionController"],
        beego.ControllerComments{
            Method: "CambiarEstadoEvaluacion",
            Router: "/cambiar-estado/",
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/evaluacion_cumplido_prov_mid/controllers:EvaluacionCumplidoController"] = append(beego.GlobalControllerRouter["github.com/udistrital/evaluacion_cumplido_prov_mid/controllers:EvaluacionCumplidoController"],
        beego.ControllerComments{
            Method: "SubirEvaluacionCumplido",
            Router: "/:evaluacion_id",
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/evaluacion_cumplido_prov_mid/controllers:FirmaElectronica"] = append(beego.GlobalControllerRouter["github.com/udistrital/evaluacion_cumplido_prov_mid/controllers:FirmaElectronica"],
        beego.ControllerComments{
            Method: "FirmarEvaluacion",
            Router: "/firmar_evaluacion/",
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/evaluacion_cumplido_prov_mid/controllers:InformacionEvaluacionController"] = append(beego.GlobalControllerRouter["github.com/udistrital/evaluacion_cumplido_prov_mid/controllers:InformacionEvaluacionController"],
        beego.ControllerComments{
            Method: "ObtenerInformacionEvaluacion",
            Router: "/:asignacion_id",
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/evaluacion_cumplido_prov_mid/controllers:NotificacionesEvaluacionController"] = append(beego.GlobalControllerRouter["github.com/udistrital/evaluacion_cumplido_prov_mid/controllers:NotificacionesEvaluacionController"],
        beego.ControllerComments{
            Method: "NotificacionAsignacionEvaluacion",
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

}
