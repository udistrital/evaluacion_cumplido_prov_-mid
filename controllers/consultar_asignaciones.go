package controllers

import (
	"github.com/astaxie/beego"
	"github.com/udistrital/evaluacion_cumplido_prov_mid/services"
	"github.com/udistrital/utils_oas/errorhandler"
	"github.com/udistrital/utils_oas/requestresponse"
)

// Consultar-Asignaciones-Controller operations for Consultar-Asignaciones-Controller
type ConsultarAsignacionesController struct {
	beego.Controller
}

// URLMapping ...
func (c *ConsultarAsignacionesController) URLMapping() {
	c.Mapping("ConsultarAsignaciones", c.ConsultarAsignaciones)
}

// @Title ConsultarAsignaciones por documento supervisor
// @Param Nmero de numeroDocumento path 	string	true "numeroDocumento"
// @Success 200 {object} models.CambioEstadoCumplidoResponse
// @Failure 404 {object} map[string]interface{}
// @router /:numeroDocumento [get]
func (c *ConsultarAsignacionesController) ConsultarAsignaciones() {

	defer errorhandler.HandlePanic(&c.Controller)
	response, err := services.ObtenerListaDeAsignaciones(c.Ctx.Input.Param(":numeroDocumento"))

	if err != nil {
		c.Ctx.Output.SetStatus(400)
		c.Data["json"] = requestresponse.APIResponseDTO(false, 400, nil, err.Error())
	}

	if err == nil {

		c.Ctx.Output.SetStatus(200)
		c.Data["json"] = requestresponse.APIResponseDTO(true, 200, response, "Busqueda exitosa")
	} else {
		c.Ctx.Output.SetStatus(400)
		c.Data["json"] = requestresponse.APIResponseDTO(false, 400, err)
	}
	c.ServeJSON()
}
