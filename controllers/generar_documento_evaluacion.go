package controllers

import (
	"strconv"

	"github.com/astaxie/beego"
	"github.com/udistrital/evaluacion_cumplido_prov_mid/services"
	"github.com/udistrital/utils_oas/errorhandler"
	"github.com/udistrital/utils_oas/requestresponse"
)

type DocumentoEvaluacionController struct {
	beego.Controller
}

func (c *DocumentoEvaluacionController) URLMapping() {
	c.Mapping("GenerarDocumentoEvaluacion", c.GenerarDocumentoEvaluacion)
}

// @Title GenerarDocumentoEvaluacion
// @Description Crea el pdf de la evaluación
// @Param id path string true "Id de la evaluación"
// @Success 200 {object} models.ResultadoFinalEvaluacion
// @Failure 404 {object} map[string]interface{}
// @router /:evaluacion_id [get]
func (c *DocumentoEvaluacionController) GenerarDocumentoEvaluacion() {
	defer errorhandler.HandlePanic(&c.Controller)

	evaluacion_id := c.Ctx.Input.Param(":evaluacion_id")

	id_evaluacion, _ := strconv.Atoi(evaluacion_id)

	err := services.GenerarDocumentoEvaluacion(id_evaluacion)
	if err == nil {
		c.Ctx.Output.SetStatus(200)
		c.Data["json"] = requestresponse.APIResponseDTO(true, 200, "Documento generado correctamente")
	} else {
		c.Ctx.Output.SetStatus(404)
		c.Data["json"] = requestresponse.APIResponseDTO(false, 404, nil, err.Error())
	}
	c.ServeJSON()
}
