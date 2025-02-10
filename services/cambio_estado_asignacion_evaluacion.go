package services

import (
	"fmt"
	"strconv"

	"github.com/astaxie/beego"
	"github.com/udistrital/evaluacion_cumplido_prov_mid/helpers"
	"github.com/udistrital/evaluacion_cumplido_prov_mid/models"
)

func CambioEstadoAsignacionEvaluacion(id_asiganacion int, codigo_estado string) (mapResponse map[string]interface{}, outputError error) {

	estados_asignables := map[string][]string{
		"EAG": {"ERE", "EAG"},
		"ERE": {"EAP"},
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

	//si estado_asignacion es nulo y el codigo_estado es EAG se agrega el estado
	if estado_asignacion_actual == nil && codigo_estado == "EAG" {
		agregarEstadoAsignacion(codigo_estado, id_asiganacion)
		mapResponse := make(map[string]interface{})
		mapResponse["Message"] = "Se agrego el estado EAG"
		return mapResponse, nil
	}

	//si estado_asignacion es nulo y el codigo_estado no es EAG no se puede asiganar el estado
	if estado_asignacion_actual == nil && codigo_estado != "EAG" {
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
	estados, existe := verificarSecuenciaAsignacion(estados_asignables, estado_asignacion_actual.EstadoAsignacionEvaluador.CodigoAbreviacion)

	if !existe {
		outputError = fmt.Errorf("el estado %s no se puede asignar, no se encuentra en la lista de asignables", estado_asignacion_actual.EstadoAsignacionEvaluador.CodigoAbreviacion)
		return nil, outputError
	}

	//verificar si el estado  esta en el slice de estados asignables
	existe_en_slice := verificarSliceAsignacion(estados, codigo_estado)

	if existe_en_slice {

		err := desabilitarEstadoAsignacion(estado_asignacion_actual)

		if err != nil {
			outputError = fmt.Errorf("error al desabilitar dea assignacion el estado")
			return nil, outputError
		}

		err = agregarEstadoAsignacion(codigo_estado, id_asiganacion)

		if err != nil {
			outputError = fmt.Errorf("error al agregar estado de asignacion")
			return nil, outputError
		}

		if codigo_estado == "ERE" {
			EnviarNotificacionRealizacionEvaluacion(estado_asignacion_actual.AsignacionEvaluadorId.PersonaId, strconv.Itoa(estado_asignacion_actual.AsignacionEvaluadorId.EvaluacionId.VigenciaContrato), strconv.Itoa(estado_asignacion_actual.AsignacionEvaluadorId.EvaluacionId.ContratoSuscritoId))
			cambiar_estado_evaluacion, err := VerificarYCambiarEstadoEvaluacion(estado_asignacion_actual.AsignacionEvaluadorId.EvaluacionId.Id, codigo_estado)

			if err != nil {
				outputError = fmt.Errorf("error al verificar el cambio de estado de la evaluacion")
				return nil, outputError
			}

			if cambiar_estado_evaluacion {

				if err != nil {
					outputError = fmt.Errorf("error al consultar el estado de la evaluacion")
					return nil, outputError
				}

				_, err = CambioEstadoEvaluacion(estado_asignacion_actual.AsignacionEvaluadorId.EvaluacionId.Id, "PRE")

				if err != nil {
					outputError = fmt.Errorf("error al cambiar el estado de la evaluacion")
					return nil, outputError

				}
			}

		}

		if codigo_estado == "EAP" {
			cambiar_estado_evaluacion, err := VerificarYCambiarEstadoEvaluacion(estado_asignacion_actual.AsignacionEvaluadorId.EvaluacionId.Id, codigo_estado)

			if err != nil {
				outputError = fmt.Errorf("error al verificar el cambio de estado de la evaluacion")
				return nil, outputError
			}

			if cambiar_estado_evaluacion {

				if err != nil {
					outputError = fmt.Errorf("error al consultar el estado de la evaluacion")
					return nil, outputError
				}

				_, err = CambioEstadoEvaluacion(estado_asignacion_actual.AsignacionEvaluadorId.EvaluacionId.Id, "AEV")

				if err != nil {
					outputError = fmt.Errorf("error al cambiar el estado de la evaluacion")
					return nil, outputError

				}

				_, _ = EnviarNotificacionesFinalizacionEvaluacion(estado_asignacion_actual.AsignacionEvaluadorId.EvaluacionId.Id, strconv.Itoa(estado_asignacion_actual.AsignacionEvaluadorId.EvaluacionId.ContratoSuscritoId), strconv.Itoa(estado_asignacion_actual.AsignacionEvaluadorId.EvaluacionId.VigenciaContrato))
			}

		}

		mapResponse := make(map[string]interface{})
		mapResponse["Message"] = fmt.Sprintf("Se cambio  el estado de %s a %s", estado_asignacion_actual.EstadoAsignacionEvaluador.CodigoAbreviacion, codigo_estado)
		return mapResponse, nil
	} else {
		outputError = fmt.Errorf("error al asignar  estado de asignacion")
		return nil, outputError
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
	//fmt.Println("Url Estado Asignacion", beego.AppConfig.String("UrlEvaluacionCumplidoCrud")+"/estado_asignacion_evaluador?query=CodigoAbreviacion:"+codigo_estado)
	if response, err := helpers.GetJsonWSO2Test(beego.AppConfig.String("UrlEvaluacionCumplidoCrud")+"/estado_asignacion_evaluador?query=CodigoAbreviacion:"+codigo_estado+",Activo:true", &respuesta_peticion); err == nil && response == 200 {

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

func ConsultarAsignacionesPorIdEvaluacion(id_evaluacion int) (asignacion *[]models.AsignacionEvaluador, outputError error) {
	defer func() {
		if err := recover(); err != nil {
			outputError = fmt.Errorf("%v", err)
			panic(outputError)
		}
	}()

	var respuestaPeticion map[string]interface{}
	var asignaciones_evaluador []models.AsignacionEvaluador

	query := fmt.Sprintf("/asignacion_evaluador?query=EvaluacionId.Id:%d,Activo:true&limit=-1", id_evaluacion)
	//fmt.Println("Url Asignacion", beego.AppConfig.String("UrlEvaluacionCumplidoCrud")+query)
	if response, err := helpers.GetJsonWSO2Test(beego.AppConfig.String("UrlEvaluacionCumplidoCrud")+query, &respuestaPeticion); err == nil && response == 200 {

		helpers.LimpiezaRespuestaRefactor(respuestaPeticion, &asignaciones_evaluador)
		if len(asignaciones_evaluador) > 0 && asignaciones_evaluador[0].EvaluacionId != nil {
			asignacion = &asignaciones_evaluador

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

	//fmt.Println("Url Estado Asignacion", beego.AppConfig.String("UrlEvaluacionCumplidoCrud")+"/cambio_estado_asignacion_evaluador/?query=AsignacionEvaluadorId.Id:"+strconv.Itoa(id_asiganacion)+",Activo:true")
	query := fmt.Sprintf("/cambio_estado_asignacion_evaluador/?query=AsignacionEvaluadorId.Id:%d,Activo:true", id_asiganacion)
	fmt.Println(beego.AppConfig.String("UrlEvaluacionCumplidoCrud") + query)
	if response, err := helpers.GetJsonWSO2Test(beego.AppConfig.String("UrlEvaluacionCumplidoCrud")+query, &respuestaPeticion); err == nil && response == 200 {

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

func desabilitarEstadoAsignacion(estadoAsignacion *models.CambioEstadoASignacionEnvaluacion) (outputError error) {
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

		if response := helpers.SendJson(beego.AppConfig.String("UrlEvaluacionCumplidoCrud")+query, "PUT", &respuestaPeticion, estadoAsignacion); response == nil {
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

func agregarEstadoAsignacion(codigo_abrevicaion string, id_asiganacion int) (outputError error) {
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
	fmt.Println(cambio_estado_asignacion_evaluador.EstadoAsignacionEvaluador.Id)
	if response := helpers.SendJson(beego.AppConfig.String("UrlEvaluacionCumplidoCrud")+"/cambio_estado_asignacion_evaluador", "POST", &respuestaPeticion, cambio_estado_asignacion_evaluador); response == nil {

		fmt.Println(response)
	} else {
		outputError = fmt.Errorf("error al desabilitar el estado")
		return outputError
	}
	return nil
}

func verificarSecuenciaAsignacion(estados map[string][]string, abreviacion string) (estado []string, existe bool) {

	estado, existe = estados[abreviacion]
	return estado, existe
}

func verificarSliceAsignacion(estados []string, abreviacion string) (existe bool) {

	for _, estado := range estados {

		if estado == abreviacion {
			return true
		}
	}
	return false
}

func VerificarYCambiarEstadoEvaluacion(id_evaluacion int, estado_abreviacion string) (cambiar_estado_evaluacion bool, outputError error) {

	cambiar_estado_evaluacion = true

	asiganciones, err := ConsultarAsignacionesPorIdEvaluacion(id_evaluacion)

	if err != nil {
		outputError = fmt.Errorf("error al consultar asignaciones")
		return false, outputError
	}

	fmt.Println("Asignaciones", asiganciones)
	for _, asigancion := range *asiganciones {

		estado_asiganacion, err := ConsultarEstadoActualAsingacion(asigancion.Id)

		if err != nil {
			fmt.Printf("error al consultar el estado de la asignacion %s", estado_asiganacion.Id)

		}

		if estado_asiganacion.EstadoAsignacionEvaluador.CodigoAbreviacion != estado_abreviacion {
			cambiar_estado_evaluacion = false
			break
		}

	}

	fmt.Println("Cambiar Estado Evaluacion", cambiar_estado_evaluacion)
	return cambiar_estado_evaluacion, outputError

}
