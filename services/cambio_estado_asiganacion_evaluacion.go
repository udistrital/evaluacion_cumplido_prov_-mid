package services

import (
	"fmt"

	"github.com/astaxie/beego"
	"github.com/udistrital/evaluacion_cumplido_prov_mid/helpers"
	"github.com/udistrital/evaluacion_cumplido_prov_mid/models"
)

func CambioEstadoAsignacionEvaluacion(id_asiganacion int, codigo_estado string) (mapResponse map[string]interface{}, outputError error) {

	estados_asignables := map[string][]string{
		"EA": {"ER"},
		"ER": {"EAP"},
	}

	defer func() {
		if err := recover(); err != nil {
			outputError = fmt.Errorf("%v", err)
			panic(outputError.Error())
		}
	}()

	if codigo_estado == "" || id_asiganacion == 0 {
		outputError = fmt.Errorf("error en los datos de entrada")
		return nil, outputError
	}
	//consultar estado por codigo de abreviacion
	cambio_estado_asignacion_evaluador, err := consultarEstadoAsignacionEvaluacion(codigo_estado)
	//si es nulo es por que no exite el estado
	if cambio_estado_asignacion_evaluador == nil || err != nil {
		outputError = fmt.Errorf("error al consultar el estado de la asignacion, el estado es no existe")
		return nil, outputError
	}

	//Consultar el estado actual de la asignacion
	estado_asignacion_actual, err := ConsultarEstadoActualAsingacion(id_asiganacion)

	//error al consultar el estado de la asignacion
	if err != nil {
		outputError = fmt.Errorf("error al consultar el estado de la asignacion")
		return nil, outputError
	}

	//si estado_asignacion es nulo y el codigo_estado es EA se agrega el estado
	if estado_asignacion_actual == nil && codigo_estado == "EA" {
		agregarEstado(codigo_estado, id_asiganacion)
		mapResponse := make(map[string]interface{})
		mapResponse["Message"] = "Se agrego esl estado EA"
		return mapResponse, nil
	}

	//si estado_asignacion es nulo y el codigo_estado no es EA no se puede asiganar el estado
	if estado_asignacion_actual == nil && codigo_estado != "EA" {
		fmt.Println("no se puede asiganar el estado")
		outputError = fmt.Errorf("error no se pude asginar el estado : %s, a una asigancion nueva", codigo_estado)
		return nil, outputError
	}

	// si el estado es el mismo no hay cambios
	if estado_asignacion_actual.EstadoAsignacionEvaluador.CodigoAbreviacion == codigo_estado {
		fmt.Println("no se puede asiganar el estado")
		outputError = fmt.Errorf("error no se pude asginar el estado : %s, a una asigancion nueva", codigo_estado)
		return nil, outputError
	}

	////verificar si el estado  esta en el mapa de estados asignables
	estados, existe := verificarSecuencia(estados_asignables, estado_asignacion_actual.EstadoAsignacionEvaluador.CodigoAbreviacion)

	if !existe {
		outputError = fmt.Errorf("el estado %s no se puede asignar, no se encuentra en la lista de asignables", estado_asignacion_actual.EstadoAsignacionEvaluador.CodigoAbreviacion)
		return nil, outputError
	}

	//verificar si el estado  esta en el slice de estados asignables
	existe_en_slice := verificarSlice(estados, codigo_estado)

	if existe_en_slice {
		desabilitarEstado(estado_asignacion_actual)
		agregarEstado(codigo_estado, id_asiganacion)
		mapResponse := make(map[string]interface{})
		mapResponse["Message"] = fmt.Sprintf("Se cambio  el estado de %s a %s", estado_asignacion_actual.EstadoAsignacionEvaluador.CodigoAbreviacion, codigo_estado)
		return mapResponse, nil
	}
	return nil, outputError
}

func consultarEstadoAsignacionEvaluacion(codigo_estado string) (cambio_estado_asignacion_evaluador *models.EstadoAsignacionEvaluador, outputError error) {

	defer func() {
		if err := recover(); err != nil {
			outputError = fmt.Errorf("%v", err)
			panic(outputError.Error())
		}
	}()

	var respuesta_peticion = make(map[string]interface{})
	var lista_asignacion_evaluador []models.EstadoAsignacionEvaluador
	fmt.Print(beego.AppConfig.String("urlEvaluacionCumplidosCrud") + "/estado_asignacion_evaluador?query=CodigoAbreviacion:" + codigo_estado)
	if response, err := helpers.GetJsonWSO2Test(beego.AppConfig.String("urlEvaluacionCumplidosCrud")+"/estado_asignacion_evaluador?query=CodigoAbreviacion:"+codigo_estado, &respuesta_peticion); err == nil && response == 200 {

		helpers.LimpiezaRespuestaRefactor(respuesta_peticion, &lista_asignacion_evaluador)

		if lista_asignacion_evaluador[0].Id != 0 {
			cambio_estado_asignacion_evaluador = &lista_asignacion_evaluador[0]

		} else {
			outputError = fmt.Errorf("error al consultar el estado de la asignacion")
			return nil, outputError
		}

	}
	return cambio_estado_asignacion_evaluador, nil

}

func consultarAsignacionesPorId(id_asiganacion int) (asignacion *models.AsignacionEvaluador, outputError error) {
	defer func() {
		if err := recover(); err != nil {
			outputError = fmt.Errorf("%v", err)
			panic(outputError)
		}
	}()

	var respuestaPeticion map[string]interface{}
	var asignaciones_evaluador []models.AsignacionEvaluador

	query := fmt.Sprintf("/asignacion_evaluador?query=Id:%d", id_asiganacion)
	if response, err := helpers.GetJsonWSO2Test(beego.AppConfig.String("urlEvaluacionCumplidosCrud")+query, &respuestaPeticion); err == nil && response == 200 {

		helpers.LimpiezaRespuestaRefactor(respuestaPeticion, &asignaciones_evaluador)
		if len(asignaciones_evaluador) > 0 && asignaciones_evaluador[0].EvaluacionId != nil {
			asignacion = &asignaciones_evaluador[0]

		}
	} else {
		outputError = fmt.Errorf("error al consultar asignaciones")
		return nil, outputError

	}
	return asignacion, nil
}

func ConsultarEstadoActualAsingacion(id_asiganacion int) (estado_asignacion *models.CambioEstadoASignacionEnvaluacion, outputError error) {
	defer func() {
		if err := recover(); err != nil {
			outputError = fmt.Errorf("%v", err)
			panic(outputError)
		}
	}()

	var respuestaPeticion map[string]interface{}
	var cambio_estado_asignacion_evaluador []models.CambioEstadoASignacionEnvaluacion

	query := fmt.Sprintf("/cambio_estado_asignacion_evaluador/?query=AsignacionEvaluadorId.Id:%d,Activo:true", id_asiganacion)
	if response, err := helpers.GetJsonWSO2Test(beego.AppConfig.String("urlEvaluacionCumplidosCrud")+query, &respuestaPeticion); err == nil && response == 200 {

		helpers.LimpiezaRespuestaRefactor(respuestaPeticion, &cambio_estado_asignacion_evaluador)
		if len(cambio_estado_asignacion_evaluador) > 0 && cambio_estado_asignacion_evaluador[0].Id != 0 {
			estado_asignacion = &cambio_estado_asignacion_evaluador[0]
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

func desabilitarEstado(estadoAsignacion *models.CambioEstadoASignacionEnvaluacion) (outputError error) {
	defer func() {
		if err := recover(); err != nil {
			outputError = fmt.Errorf("%v", err)
			panic(outputError)
		}
	}()
	var respuestaPeticion map[string]interface{}
	if estadoAsignacion != nil {
		estadoAsignacion.Activo = false

		query := fmt.Sprintf("/cambio_estado_asignacion_evaluador/%d", estadoAsignacion.Id)

		if response := helpers.SendJson(beego.AppConfig.String("urlEvaluacionCumplidosCrud")+query, "PUT", &respuestaPeticion, estadoAsignacion); response == nil {
		} else {
			outputError = fmt.Errorf("error al desabilitar el estado")
			return outputError
		}

		if respuestaPeticion["Success"] == false {
			outputError = fmt.Errorf("error al desabilitar el estado")
			return outputError
		}

	} else {
		outputError = fmt.Errorf("error al desabilitar el estado")
		return outputError
	}

	return nil

}

func agregarEstado(codigo_abrevicaion string, id_asiganacion int) (outputError error) {
	defer func() {
		if err := recover(); err != nil {
			outputError = fmt.Errorf("%v", err)
			panic(outputError)
		}
	}()
	var respuestaPeticion map[string]interface{}

	estadoAsignacion, err := consultarEstadoAsignacionEvaluacion(codigo_abrevicaion)

	if err != nil {
		outputError = fmt.Errorf("error al consultar el estado de la asignacion")
		return outputError
	}
	estadoAsignacionEvaluador := models.EstadoAsignacionEvaluador{Id: estadoAsignacion.Id}
	asignacionEvaluador := models.AsignacionEvaluador{Id: id_asiganacion}

	cambio_estado_asignacion_evaluador := models.CambioEstadoASignacionEnvaluacionPeticion{
		EstadoAsignacionEvaluador: estadoAsignacionEvaluador,
		AsignacionEvaluadorId:     asignacionEvaluador,
		Activo:                    true,
	}

	fmt.Println(cambio_estado_asignacion_evaluador.Activo)
	if response := helpers.SendJson(beego.AppConfig.String("urlEvaluacionCumplidosCrud")+"/cambio_estado_asignacion_evaluador", "POST", &respuestaPeticion, cambio_estado_asignacion_evaluador); response == nil {

		fmt.Println(response)
	} else {
		outputError = fmt.Errorf("error al desabilitar el estado")
		return outputError
	}
	return nil
}

func verificarSecuencia(estados map[string][]string, abreviacion string) (estado []string, existe bool) {
	// Asignamos directamente el valor del mapa a la variable estado y el valor de existencia
	estado, existe = estados[abreviacion]
	return estado, existe
}

func verificarSlice(estados []string, abreviacion string) (existe bool) {

	for _, estado := range estados {

		if estado == abreviacion {
			return true
		}
	}
	return false
}

// func Consultar_evaluacion_id(id_evaluacion int) (evaluacion *models.Evaluacion, outputError error) {
// 	defer func() {
// 		if err := recover(); err != nil {
// 			outputError = fmt.Errorf("%v", err)
// 			panic(outputError)
// 		}
// 	}()

// 	var respuestaPeticion map[string]interface{}
// 	var evaluaciones []models.Evaluacion

// 	query := fmt.Sprintf("/evaluacion/?query=Id:%d", id_evaluacion)
// 	fmt.Println(beego.AppConfig.String("urlEvaluacionCumplidosCrud") + query)
// 	if response, err := helpers.GetJsonWSO2Test(beego.AppConfig.String("urlEvaluacionCumplidosCrud")+query, &respuestaPeticion); err == nil && response == 200 {

// 		helpers.LimpiezaRespuestaRefactor(respuestaPeticion, &evaluaciones)
// 		if len(evaluaciones) > 0 && evaluaciones[0].Id != 0 {
// 			evaluacion = &evaluaciones[0]

// 		}
// 	} else {
// 		outputError = fmt.Errorf("error al consultar asignaciones")
// 		return nil, outputError

// 	}
// 	fmt.Println(evaluacion)
// 	return evaluacion, nil

// }
