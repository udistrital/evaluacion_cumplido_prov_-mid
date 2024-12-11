package controllers

import (
	"github.com/astaxie/beego"
	"github.com/udistrital/evaluacion_cumplido_prov_mid/services"
	"github.com/udistrital/utils_oas/errorhandler"
	"github.com/udistrital/utils_oas/requestresponse"
)

type CambioRolEvaluadorController struct {
	beego.Controller
}

func (c *CambioRolEvaluadorController) URLMapping() {

}

// @Title CambiarRolAsignacionEvaluador
// @Param Nmero de idEvaluacion path 	string	true "idEvaluacion"
// @Success 200 {object} models.CambioEstadoCumplidoResponse
// @Failure 404 {object} map[string]interface{}
// @router /:idEvaluacion [post]
func (c *CambioRolEvaluadorController) CambiarRolAsignacionEvaluador() {

	defer errorhandler.HandlePanic(&c.Controller)
	var evaluacionId = c.Ctx.Input.Param(":idEvaluacion")
	response, err := services.CambiarRolAsignacionEvaluador(evaluacionId)

	if err != nil {
		c.Ctx.Output.SetStatus(400)
		c.Data["json"] = requestresponse.APIResponseDTO(false, 400, nil, err.Error())
	}

	if err == nil {

		c.Ctx.Output.SetStatus(200)
		c.Data["json"] = requestresponse.APIResponseDTO(true, 200, response, "Cambio de rol exitoso")
	} else {
		c.Ctx.Output.SetStatus(400)
		c.Data["json"] = requestresponse.APIResponseDTO(false, 400, err)
	}
	c.ServeJSON()

}
