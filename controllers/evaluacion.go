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

type EvaluacionController struct {
	beego.Controller
}

func (c *EvaluacionController) URLMapping() {
	c.Mapping("CambiarEstadoEvaluacion", c.CambiarEstadoEvaluacion)
}

// @Title CambiarEstadoEvaluacion por documento supervisor
// @Param Nmero de numeroDocumento path 	string	true "numeroDocumento"
// @Success 200 {object} models.CambioEstadoCumplidoResponse
// @Failure 404 {object} map[string]interface{}
// @router /cambiar-estado/ [post]
func (c *EvaluacionController) CambiarEstadoEvaluacion() {
	defer errorhandler.HandlePanic(&c.Controller)

	var v models.PeticionCambioEstadoEvaluacion
	json.Unmarshal(c.Ctx.Input.RequestBody, &v)

	fmt.Println(v)
	response, err := services.CambioEstadoEvaluacion(v.EvaluacionId.Id, v.AbreviacionEstado)

	if err == nil {

		c.Ctx.Output.SetStatus(200)
		c.Data["json"] = requestresponse.APIResponseDTO(true, 200, nil, response)
	} else {
		c.Ctx.Output.SetStatus(400)
		c.Data["json"] = requestresponse.APIResponseDTO(false, 400, nil, err.Error())
	}
	c.ServeJSON()
}
