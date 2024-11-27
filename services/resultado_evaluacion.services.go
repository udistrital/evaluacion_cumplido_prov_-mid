package services

import (
	"encoding/json"
	"fmt"

	"github.com/astaxie/beego"
	"github.com/udistrital/evaluacion_cumplido_prov_mid/helpers"
	"github.com/udistrital/evaluacion_cumplido_prov_mid/models"
)

func GuardarResultadoEvaluacion(resultado models.BodyResultadoEvaluacion) (outputError error) {
	defer func() {
		if err := recover(); err != nil {
			outputError = fmt.Errorf("%v", err)
			panic(outputError.Error())
		}
	}()

	resultado_json, err := json.Marshal(resultado.ResultadoEvaluacion)
	if err != nil {
		outputError = fmt.Errorf("Error al convertir a Json el resultado de la evaluación")
		return outputError
	}

	resultado_string := string(resultado_json)
	resultado_map := make(map[string]interface{})

	resultado_map["AsignacionEvaluadorId"] = map[string]interface{}{"Id": resultado.AsignacionEvaluadorId}
	resultado_map["ClasificacionId"] = map[string]interface{}{"Id": resultado.ClasificacionId}
	resultado_map["ResultadoEvaluacion"] = resultado_string
	resultado_map["Observaciones"] = resultado.Observaciones

	var respuesta_peticion map[string]interface{}

	//fmt.Println("URL: ", beego.AppConfig.String("UrlEvaluacionesCumplidosProveedoresCrud")+"/resultado_evaluacion")
	if err := helpers.SendJson(beego.AppConfig.String("UrlEvaluacionesCumplidosProveedoresCrud")+"/resultado_evaluacion", "POST", &respuesta_peticion, resultado_map); err != nil {
		outputError = fmt.Errorf("Error al guardar el resultado de la evaluación")
		return outputError
	}

	return nil

}
