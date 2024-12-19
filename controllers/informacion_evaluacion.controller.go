package controllers

import (
	"github.com/astaxie/beego"
	"github.com/udistrital/evaluacion_cumplido_prov_mid/services"
	"github.com/udistrital/utils_oas/errorhandler"
	"github.com/udistrital/utils_oas/requestresponse"
)

// InformacionEvaluacionController operations for InformacionEvaluacionController
type InformacionEvaluacionController struct {
	beego.Controller
}

// URLMapping ...
func (c *InformacionEvaluacionController) URLMapping() {
	c.Mapping("ObtenerInformacionEvaluacion", c.ObtenerInformacionEvaluacion)

}

// ObtenerInformacionEvaluacion ...
// @Title ObtenerInformacionEvaluacion
// @Description get InformacionEvaluacionController by asignacion_id
// @Param	id		path 	string	true		"The key for staticblock"
// @Success 200 {object} models.InformacionEvaluacion
// @Failure 400 :id is empty
// @router /:asignacion_id [get]
func (c *InformacionEvaluacionController) ObtenerInformacionEvaluacion() {

	defer errorhandler.HandlePanic(&c.Controller)

	asignacion_id := c.Ctx.Input.Param(":asignacion_id")

	data, err := services.ObtenerInformacionEvaluacion(asignacion_id)
	if err == nil {
		c.Ctx.Output.SetStatus(200)
		c.Data["json"] = requestresponse.APIResponseDTO(true, 200, data)
	} else {
		c.Ctx.Output.SetStatus(400)
		c.Data["json"] = requestresponse.APIResponseDTO(false, 400, nil, err.Error())
	}
	c.ServeJSON()
}
