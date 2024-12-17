package services

import (
	"fmt"

	"github.com/astaxie/beego"
	"github.com/udistrital/evaluacion_cumplido_prov_mid/helpers"
	"github.com/udistrital/evaluacion_cumplido_prov_mid/models"
)

func consultarEstadoActualEvaluacion(id_evaluacion int) (estado_asignacion *models.CambioEstadoEvaluacion, outputError error) {

	defer func() {
		if err := recover(); err != nil {
			outputError = fmt.Errorf("%v", err)
			panic(outputError)
		}
	}()

	var respuestaPeticion map[string]interface{}
	var cambio_estado_evaluacion []models.CambioEstadoEvaluacion

	query := fmt.Sprintf("/cambio_estado_evaluacion/?query=EvaluacionId.Id:%d,Activo:true", id_evaluacion)
	if response, err := helpers.GetJsonWSO2Test(beego.AppConfig.String("urlEvaluacionCumplidosCrud")+query, &respuestaPeticion); err == nil && response == 200 {

		helpers.LimpiezaRespuestaRefactor(respuestaPeticion, &cambio_estado_evaluacion)
		if len(cambio_estado_evaluacion) > 0 && cambio_estado_evaluacion[0].Id != 0 {
			estado_asignacion = &cambio_estado_evaluacion[0]
			fmt.Println(estado_asignacion)

		} else {
			return estado_asignacion, nil
		}

	} else {
		outputError = fmt.Errorf("error al consultar asignaciones")
		return nil, outputError
	}

	return estado_asignacion, nil
}

func consultarEstadoEvaluacion(codigo_abreviacion_estado_evalauacion string) (estado__evaluacion *models.EstadoEvaluacion, outputError error) {

	defer func() {
		if err := recover(); err != nil {
			outputError = fmt.Errorf("%v", err)
			panic(outputError.Error())
		}
	}()

	var respuesta_peticion = make(map[string]interface{})
	var lista_estados_evaluacion []models.EstadoEvaluacion
	query := fmt.Sprintf("/estado_evaluacion?query=CodigoAbreviacion:%s", codigo_abreviacion_estado_evalauacion)

	if response, err := helpers.GetJsonWSO2Test(beego.AppConfig.String("urlEvaluacionCumplidosCrud")+query, &respuesta_peticion); err == nil && response == 200 {

		helpers.LimpiezaRespuestaRefactor(respuesta_peticion, &lista_estados_evaluacion)

		if lista_estados_evaluacion[0].Id != 0 {
			estado__evaluacion = &lista_estados_evaluacion[0]

		} else {
			outputError = fmt.Errorf("error al consultar el estado de la asignacion")
			return nil, outputError
		}

	}
	return estado__evaluacion, nil

}

func desabilitarEstadoEvaluacion(estado_evalacion *models.CambioEstadoEvaluacion) (outputError error) {
	defer func() {
		if err := recover(); err != nil {
			outputError = fmt.Errorf("%v", err)
			panic(outputError)
		}
	}()
	var respuestaPeticion map[string]interface{}
	if estado_evalacion != nil {
		estado_evalacion.Activo = false

		query := fmt.Sprintf("/cambio_estado_evaluacion/%d", estado_evalacion.Id)
		if response := helpers.SendJson(beego.AppConfig.String("urlEvaluacionCumplidosCrud")+query, "PUT", &respuestaPeticion, estado_evalacion); response == nil {
		} else {
			outputError = fmt.Errorf("error al desabilitar el estado")
			return outputError
		}

	} else {
		outputError = fmt.Errorf("error al desabilitar el estado")
		return outputError
	}

	return nil

}

func agregarEstadoEvaluacion(codigo_abreviacion string, id_evaluacion int) (outputError error) {
	defer func() {
		if err := recover(); err != nil {
			outputError = fmt.Errorf("%v", err)
			panic(outputError)
		}
	}()
	var respuestaPeticion map[string]interface{}

	estado_evaluacion_consulta, err := consultarEstadoEvaluacion(codigo_abreviacion)

	if err != nil {
		outputError = fmt.Errorf("error al consultar el estado de la asignacion")
		return outputError
	}
	evaluacion := models.Evaluacion{Id: id_evaluacion}
	estado_evaluacion := models.EstadoEvaluacion{Id: estado_evaluacion_consulta.Id}

	cambio_estado_evaluacion := models.CambioEstadoEvaluacion{
		EvaluacionId:       &evaluacion,
		EstadoEvaluacionId: &estado_evaluacion,
		Activo:             true,
	}

	fmt.Println(cambio_estado_evaluacion.Activo)
	if response := helpers.SendJson(beego.AppConfig.String("urlEvaluacionCumplidosCrud")+"/cambio_estado_evaluacion", "POST", &respuestaPeticion, cambio_estado_evaluacion); response == nil {

		fmt.Println(response)
	} else {
		outputError = fmt.Errorf("error al desabilitar el estado")
		return outputError
	}
	return nil
}
