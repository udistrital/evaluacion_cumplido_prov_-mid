package controllers

import (
	"github.com/astaxie/beego"
	"github.com/udistrital/evaluacion_cumplido_prov_mid/services"
	"github.com/udistrital/utils_oas/errorhandler"
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

	services.ConsultarAsignaciones(c.Ctx.Input.Param(":numeroDocumento"))
}
