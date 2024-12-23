package controllers

import (
	"encoding/json"

	"github.com/astaxie/beego"
	"github.com/udistrital/evaluacion_cumplido_prov_mid/models"
	"github.com/udistrital/evaluacion_cumplido_prov_mid/services"
	"github.com/udistrital/utils_oas/errorhandler"
	"github.com/udistrital/utils_oas/requestresponse"
)

type FirmaElectronica struct {
	beego.Controller
}

func (c *FirmaElectronica) URLMapping() {
	c.Mapping("FirmarEvaluacion", c.FirmarEvaluacion)
}

// @Title FirmarDocumento
// @router /firmar_evaluacion/ [post]
func (c *FirmaElectronica) FirmarEvaluacion() {
	defer errorhandler.HandlePanic(&c.Controller)

	var v models.PeticionFirmaElectronica
	json.Unmarshal(c.Ctx.Input.RequestBody, &v)
	response, err := services.EjecutarProcesoDefirma(v)

	if err == nil {

		c.Ctx.Output.SetStatus(200)
		c.Data["json"] = requestresponse.APIResponseDTO(true, 200, response["Data"], response["Message"])
	} else {
		c.Ctx.Output.SetStatus(400)
		c.Data["json"] = requestresponse.APIResponseDTO(false, 400, nil, err.Error())
	}
	c.ServeJSON()
}
