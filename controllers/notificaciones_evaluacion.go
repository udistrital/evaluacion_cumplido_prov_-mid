package controllers

import (
	"strconv"

	"github.com/astaxie/beego"
	"github.com/udistrital/evaluacion_cumplido_prov_mid/services"
	"github.com/udistrital/utils_oas/errorhandler"
	"github.com/udistrital/utils_oas/requestresponse"
)

type NotificacionesEvaluacionController struct {
	beego.Controller
}

func (c *NotificacionesEvaluacionController) URLMapping() {
	c.Mapping("NotificacionAsignacionEvaluacion", c.NotificacionAsignacionEvaluacion)
}

// @Title NotificacionAsignacionEvaluacion
// @Description Notifica a los evaluadores que se les ha asignado una evaluación
// @Param id path string true "Id de la evaluación"
// @Success 200 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @router /:evaluacion_id [get]
func (c *NotificacionesEvaluacionController) NotificacionAsignacionEvaluacion() {
	defer errorhandler.HandlePanic(&c.Controller)

	evaluacion_id := c.Ctx.Input.Param(":evaluacion_id")

	id_evaluacion, _ := strconv.Atoi(evaluacion_id)

	data, err := services.EnviarNotificacionesAsignacionEvaluacion(id_evaluacion)
	if err == nil && len(data) == 0 {
		c.Ctx.Output.SetStatus(200)
		c.Data["json"] = requestresponse.APIResponseDTO(true, 200, nil, "Notificaciones enviadas correctamente")
	} else if err == nil && len(data) > 0 {
		c.Ctx.Output.SetStatus(200)
		c.Data["json"] = requestresponse.APIResponseDTO(false, 200, data, "Notificaciones enviadas con errores")
	} else {
		c.Ctx.Output.SetStatus(404)
		c.Data["json"] = requestresponse.APIResponseDTO(false, 404, nil, err.Error())
	}
	c.ServeJSON()
}
