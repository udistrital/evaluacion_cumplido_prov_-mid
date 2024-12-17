package controllers

import (
	"github.com/astaxie/beego"
	"github.com/udistrital/evaluacion_cumplido_prov_mid/services"
	"github.com/udistrital/utils_oas/errorhandler"
	"github.com/udistrital/utils_oas/requestresponse"
)

// EvaluacionCumplidoController handles operations for uploading evaluations
// to completed contracts.
type EvaluacionCumplidoController struct {
	beego.Controller
}

// URLMapping defines the mappings for the controller.
func (c *EvaluacionCumplidoController) URLMapping() {
	c.Mapping("SubirEvaluacionCumplido", c.SubirEvaluacionCumplido)
}

// SubirEvaluacionCumplido handles the uploading of evaluations to completed contracts.
// @Title SubirEvaluacionCumplido
// @Description Upload evaluations to completed contracts by evaluacion_id
// @Param	evaluacion_id	path	string	true	"ID of the evaluation to upload"
// @Success 200 {object} []models.CumplidoProveedor
// @Failure 400 :evaluacion_id is invalid
// @Failure 404 No valid completed contracts found
// @router /:evaluacion_id [get]
func (c *EvaluacionCumplidoController) SubirEvaluacionCumplido() {

	defer errorhandler.HandlePanic(&c.Controller)

	evaluacionID := c.Ctx.Input.Param(":evaluacion_id")

	cumplidos, err := services.SubirEvaluacionCumplido(evaluacionID)
	if err == nil {
		c.Ctx.Output.SetStatus(200)
		c.Data["json"] = requestresponse.APIResponseDTO(true, 200, cumplidos, "Archivo subido con exito al cumplido")
	} else if err != nil && len(cumplidos) > 0 {
		c.Ctx.Output.SetStatus(200)
		c.Data["json"] = requestresponse.APIResponseDTO(true, 200, cumplidos, err.Error())
	} else {
		c.Ctx.Output.SetStatus(400)
		c.Data["json"] = requestresponse.APIResponseDTO(false, 400, nil, err.Error())
	}
	c.ServeJSON()
}
