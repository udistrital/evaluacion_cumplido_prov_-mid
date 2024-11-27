package controllers

import (
	"encoding/json"

	"github.com/astaxie/beego"
	"github.com/udistrital/evaluacion_cumplido_prov_mid/models"
	"github.com/udistrital/evaluacion_cumplido_prov_mid/services"
	"github.com/udistrital/utils_oas/errorhandler"
	"github.com/udistrital/utils_oas/requestresponse"
)

type ResultadoEvaluacionController struct {
	beego.Controller
}

func (c *ResultadoEvaluacionController) URLMapping() {
	c.Mapping("GuardarResultadoEvaluacion", c.GuardarResultadoEvaluacion)
}

// @Title GuardarResultadoEvaluacion
// @Description Guarda el resultado de una evaluación de cumplimiento de proveedor
// @Param body body BodyResultadoEvaluacion true "Estructura con el resultado de la evaluación"
// @Success 200 {object} models.CambioEstadoCumplidoResponse
// @Failure 404 {object} map[string]interface{}
// @router /resultado-evaluacion [post]
func (c *ResultadoEvaluacionController) GuardarResultadoEvaluacion() {
	defer errorhandler.HandlePanic(&c.Controller)

	var v models.BodyResultadoEvaluacion

	json.Unmarshal(c.Ctx.Input.RequestBody, &v)

	response, err := services.GuardarResultadoEvaluacion(v)
	if err == nil {
		c.Ctx.Output.SetStatus(200)
		c.Data["json"] = requestresponse.APIResponseDTO(true, 200, response)
	} else {
		c.Ctx.Output.SetStatus(404)
		c.Data["json"] = requestresponse.APIResponseDTO(false, 404, err)
	}
	c.ServeJSON()
}
