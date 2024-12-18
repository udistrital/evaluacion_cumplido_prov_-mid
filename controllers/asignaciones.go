package controllers

import (
	"encoding/json"
	"fmt"

	"github.com/astaxie/beego"
	"github.com/udistrital/evaluacion_cumplido_prov_mid/models"
	"github.com/udistrital/evaluacion_cumplido_prov_mid/services"
	"github.com/udistrital/utils_oas/errorhandler"
	"github.com/udistrital/utils_oas/requestresponse"
)

// Consultar-Asignaciones-Controller operations for Consultar-Asignaciones-Controller
type AsignacionesController struct {
	beego.Controller
}

// URLMapping ...
func (c *AsignacionesController) URLMapping() {
	c.Mapping("ConsultarAsignaciones", c.ConsultarAsignaciones)
	c.Mapping("CambiarEstadoAsignacionEvaluacion", c.CambiarEstadoAsignacionEvaluacion)
}

// @Title ConsultarAsignaciones por documento supervisor
// @Param Nmero de numeroDocumento path 	string	true "numeroDocumento"
// @Success 200 {object} models.CambioEstadoCumplidoResponse
// @Failure 404 {object} map[string]interface{}
// @router /consultar/:numeroDocumento [get]
func (c *AsignacionesController) ConsultarAsignaciones() {
	defer errorhandler.HandlePanic(&c.Controller)
	response, err := services.ObtenerListaDeAsignaciones(c.Ctx.Input.Param(":numeroDocumento"))

	if err == nil {

		c.Ctx.Output.SetStatus(200)
		c.Data["json"] = requestresponse.APIResponseDTO(true, 200, nil, response)
	} else {
		c.Ctx.Output.SetStatus(400)
		c.Data["json"] = requestresponse.APIResponseDTO(false, 400, nil, err.Error())
	}
	c.ServeJSON()
}

// @Title CambiarEstadoAsignacion por documento supervisor
// @Param Nmero de numeroDocumento path 	string	true "numeroDocumento"
// @Success 200 {object} models.CambioEstadoCumplidoResponse
// @Failure 404 {object} map[string]interface{}
// @router /cambiar-estado/ [post]
func (c *AsignacionesController) CambiarEstadoAsignacionEvaluacion() {
	defer errorhandler.HandlePanic(&c.Controller)

	var v models.PeticionCambioEstadoAsignacion
	json.Unmarshal(c.Ctx.Input.RequestBody, &v)

	fmt.Println(v)
	response, err := services.CambioEstadoAsignacionEvaluacion(v.AsignacionId.Id, v.AbreviacionEstado)

	if err == nil {

		c.Ctx.Output.SetStatus(200)
		c.Data["json"] = requestresponse.APIResponseDTO(true, 200, nil, response)
	} else {
		c.Ctx.Output.SetStatus(400)
		c.Data["json"] = requestresponse.APIResponseDTO(false, 400, nil, err.Error())
	}
	c.ServeJSON()
}
