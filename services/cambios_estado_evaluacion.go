package services

import (
	"fmt"

	"github.com/astaxie/beego"
	"github.com/udistrital/evaluacion_cumplido_prov_mid/helpers"
	"github.com/udistrital/evaluacion_cumplido_prov_mid/models"
)

func CambioEstadoEvaluacion(id_evaluacion int, codigo_estado string) (mapResponse map[string]interface{}, outputError error) {
	defer func() {
		if err := recover(); err != nil {
			outputError = fmt.Errorf("%v", err)
			panic(outputError)
		}
	}()

	estado_actual, err := ConsultarEstadoActualEvaluacion(id_evaluacion)

	if err != nil {
		outputError = fmt.Errorf("error al consultar el estado de la evaluacion")
		return nil, outputError
	}

	estado_a_asignar, err := ConsultarEstadoEvaluacion(codigo_estado)

	if err != nil {
		outputError = fmt.Errorf("el estado a asignar no existe")
		return nil, outputError
	}

	estados_asignables := map[string][]string{
		"GNT": {"EPR"},
		"EPR": {"PRE"},
		"PRE": {"AEV"},
	}

	lista_asignable, existe := verificarSecuenciaEvaluacion(estados_asignables, estado_actual.EstadoEvaluacionId.CodigoAbreviacion)

	if !existe {
		outputError = fmt.Errorf("el estado %s no se puede asignar, no se encuentra en la lista de asignables", estado_actual.EstadoEvaluacionId.CodigoAbreviacion)
		return nil, outputError
	}

	existe_en_slice := verificarSliceEvaluacion(lista_asignable, estado_a_asignar.CodigoAbreviacion)

	if existe_en_slice {

		err = DesabilitarEstadoEvaluacion(estado_actual)
		if err != nil {
			outputError = fmt.Errorf("error al desabilitar el estado")
			return nil, outputError
		}

		err = AgregarEstadoEvaluacion(codigo_estado, id_evaluacion)
		if err != nil {
			outputError = fmt.Errorf("error al agregar el estado")
			return nil, outputError
		}

		if codigo_estado == "EPR" {

			lista_asignaciones, err := ConsultarAsignacionesPorIdEvaluacion(id_evaluacion)

			if err != nil {
				outputError = fmt.Errorf("error al consultar asignaciones")
				return nil, outputError
			}

			for _, asignacion := range *lista_asignaciones {

				_, err = CambioEstadoAsignacionEvaluacion(asignacion.Id, "EA")

			}

		}

		mapResponse := make(map[string]interface{})
		mapResponse["Message"] = fmt.Sprintf("Se cambio  el estado de %s a %s", codigo_estado, estado_a_asignar.CodigoAbreviacion)
		return mapResponse, nil

	} else {
		outputError = fmt.Errorf("el estado %s no se puede asignar, no se encuentra en la lista de asignables", estado_actual.EstadoEvaluacionId.CodigoAbreviacion)
		return nil, outputError
	}

}

func ConsultarEstadoActualEvaluacion(id_evaluacion int) (estado_asignacion *models.CambioEstadoEvaluacion, outputError error) {

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

func ConsultarEstadoEvaluacion(codigo_abreviacion_estado_evalauacion string) (estado__evaluacion *models.EstadoEvaluacion, outputError error) {

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

func DesabilitarEstadoEvaluacion(estado_evalacion *models.CambioEstadoEvaluacion) (outputError error) {
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

func AgregarEstadoEvaluacion(codigo_abreviacion string, id_evaluacion int) (outputError error) {
	defer func() {
		if err := recover(); err != nil {
			outputError = fmt.Errorf("%v", err)
			panic(outputError)
		}
	}()
	var respuestaPeticion map[string]interface{}

	estado_evaluacion_consulta, err := ConsultarEstadoEvaluacion(codigo_abreviacion)

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

func verificarSecuenciaEvaluacion(estados map[string][]string, abreviacion string) (lista_estado []string, existe bool) {

	lista_estado, existe = estados[abreviacion]
	return lista_estado, existe
}

func verificarSliceEvaluacion(estados []string, abreviacion string) (existe bool) {

	for _, estado := range estados {

		if estado == abreviacion {
			return true
		}
	}
	return false
}
