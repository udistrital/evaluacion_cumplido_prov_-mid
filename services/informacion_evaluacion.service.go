package services

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/astaxie/beego"
	"github.com/udistrital/evaluacion_cumplido_prov_mid/helpers"
	"github.com/udistrital/evaluacion_cumplido_prov_mid/models"
)

func ObtenerInformacionEvaluacion(asignacion_evaluacion_id string) (informacion_evaluacion models.InformacionEvaluacion, outputError error) {
	defer func() {
		if err := recover(); err != nil {
			outputError = fmt.Errorf("%v", err)
			panic(outputError)
		}
	}()

	var respuesta_asignacion_evaluador map[string]interface{}
	var asignacion_evaluadores []models.AsignacionEvaluador

	// Se Busca la asignacion evaluador por id
	//fmt.Println("URL evaluadores: ", beego.AppConfig.String("UrlEvaluacionesCumplidosProveedoresCrud")+"/asignacion_evaluador/?query=EvaluacionId.Id:"+strconv.Itoa(evaluacion_id)+",Activo:true&limit=-1")
	if response, err := helpers.GetJsonTest(beego.AppConfig.String("UrlEvaluacionesCumplidosProveedoresCrud")+"/asignacion_evaluador/?query=Id:"+asignacion_evaluacion_id+",Activo:true&limit=-1", &respuesta_asignacion_evaluador); err != nil && response != 200 {
		outputError = fmt.Errorf("Error al obtener la asignación del evaluador")
		return informacion_evaluacion, outputError
	}

	data := respuesta_asignacion_evaluador["Data"].([]interface{})
	if len(data[0].(map[string]interface{})) == 0 {
		outputError = fmt.Errorf(fmt.Sprintf("No se encontró la asignación del evaluador con el id %s", asignacion_evaluacion_id))
		return informacion_evaluacion, outputError
	}

	helpers.LimpiezaRespuestaRefactor(respuesta_asignacion_evaluador, &asignacion_evaluadores)

	// Obtener el nombre del evaluador
	nombre_evaluador, error := helpers.ObtenerNombrePersonaNatural(strconv.Itoa(asignacion_evaluadores[0].PersonaId))
	if error != nil {
		outputError = fmt.Errorf(error.Error())
		return informacion_evaluacion, outputError
	}

	// Obtener el resultado evaluacion del evaluador
	var resultado models.Resultado
	resultado_evaluacion, err := helpers.ObtenerResultadoEvaluacion(asignacion_evaluadores[0].Id)
	if err != nil {
		informacion_evaluacion.ResultadoEvaluacion = resultado

	} else {
		// Convertir el resultado evaluacion a la estructura Resultado
		error_json := json.Unmarshal([]byte(resultado_evaluacion.ResultadoEvaluacion), &resultado)
		if error_json != nil {
			outputError = fmt.Errorf("Error al convertir el resultado de la evaluación")
			return informacion_evaluacion, outputError
		}
		informacion_evaluacion.FechaEvaluacion = resultado_evaluacion.FechaCreacion.Format("2006-01-02")
		informacion_evaluacion.ResultadoEvaluacion = resultado
	}

	// Obtener evaluadores
	evaluadores, error := ObtenerEvaluadores(asignacion_evaluadores[0])
	if error != nil {
		outputError = fmt.Errorf(error.Error())
		return informacion_evaluacion, outputError
	}

	// Calcular el puntaje total de la evaluacion
	clasificacion, puntaje_total_evaluacion, error_clasificacion := helpers.CalcularClasificacionEvaluacion(resultado)
	if error_clasificacion != nil {
		informacion_evaluacion.Clasificacion = ""
		informacion_evaluacion.PuntajeTotalEvaluacion = 0
	} else {
		informacion_evaluacion.Clasificacion = clasificacion.Nombre
		informacion_evaluacion.PuntajeTotalEvaluacion = puntaje_total_evaluacion
	}

	// Obtener los datos del contrato
	contrato_general, error_contrato := helpers.ObtenerContratoGeneral(strconv.Itoa(asignacion_evaluadores[0].EvaluacionId.ContratoSuscritoId), strconv.Itoa(asignacion_evaluadores[0].EvaluacionId.VigenciaContrato))
	if error_contrato != nil {
		outputError = fmt.Errorf(error_contrato.Error())
		return informacion_evaluacion, outputError
	}

	// Obtener la dependencia evaluadora
	dependencia_evaluadora, error_dependencia := helpers.ObtenerDependenciasSupervisor(strconv.Itoa(contrato_general.Supervisor.Documento))
	if error_dependencia != nil {
		informacion_evaluacion.DependenciaEvaluadora = ""
	} else {
		for _, dependencia := range dependencia_evaluadora {
			if contrato_general.Supervisor.DependenciaSupervisor == dependencia.Codigo {
				informacion_evaluacion.DependenciaEvaluadora = dependencia.Nombre
			}
		}
	}
	//Obtener los datos del proveedor
	var informacion_proveedor []models.InformacionProveedor
	if response, err := helpers.GetJsonTest(beego.AppConfig.String("UrlAmazonApi")+"/informacion_proveedor/?query=Id:"+strconv.Itoa(contrato_general.Contratista), &informacion_proveedor); (err != nil) && (response != 200) {
		outputError = fmt.Errorf("Error al obtener la informacion del proveedor")
		return informacion_evaluacion, outputError
	}

	//Obtener items evaluados, si es supervisor se obtienen los items evaluados de la evaluacion, si es evaluador se obtienen los items evaluados por el evaluador
	if asignacion_evaluadores[0].RolAsignacionEvaluadorId.CodigoAbreviacion == "SP" {
		items_evaluacion, _, error_items := helpers.ObtenerItemsEvaluacion(asignacion_evaluadores[0].EvaluacionId.Id)
		if error_items != nil {
			informacion_evaluacion.ItemsEvaluados = []models.Item{}
		} else {
			informacion_evaluacion.ItemsEvaluados = items_evaluacion
		}
	} else {
		items_evaluador, _, error_items := helpers.ObtenerItemsEvaluador(asignacion_evaluadores[0].Id)
		if error_items != nil {
			informacion_evaluacion.ItemsEvaluados = []models.Item{}
		} else {
			informacion_evaluacion.ItemsEvaluados = items_evaluador
		}
	}

	// LLenar los datos del modelo a retornar
	informacion_evaluacion.NombreEvaluador = nombre_evaluador
	informacion_evaluacion.Cargo = asignacion_evaluadores[0].Cargo
	informacion_evaluacion.DependenciaEvaluadora = dependencia_evaluadora[0].Nombre
	informacion_evaluacion.EmpresaProveedor = informacion_proveedor[0].NomProveedor
	informacion_evaluacion.ObjetoContrato = contrato_general.ObjetoContrato
	informacion_evaluacion.Evaluadores = evaluadores
	informacion_evaluacion.CodigoAbreviacionRol = asignacion_evaluadores[0].RolAsignacionEvaluadorId.CodigoAbreviacion

	return informacion_evaluacion, nil

}

func ObtenerEvaluadores(asignacion_evaluador models.AsignacionEvaluador) (evaluadores []models.Evaluador, outputError error) {
	defer func() {
		if err := recover(); err != nil {
			outputError = fmt.Errorf("%v", err)
			panic(outputError)
		}
	}()

	// Validar si el evaluador es supervsior o evaluador
	if asignacion_evaluador.RolAsignacionEvaluadorId.CodigoAbreviacion == "SP" {

		// Obtener los evaluadores de la evaluacion
		var respuesta_evaluadores map[string]interface{}
		var evaluadores_asignacion []models.AsignacionEvaluador

		if response, err := helpers.GetJsonTest(beego.AppConfig.String("UrlEvaluacionesCumplidosProveedoresCrud")+"/asignacion_evaluador/?query=EvaluacionId.Id:"+strconv.Itoa(asignacion_evaluador.EvaluacionId.Id)+",Activo:true&limit=-1", &respuesta_evaluadores); err != nil && response != 200 {
			outputError = fmt.Errorf("Error al obtener los evaluadores de la evaluación")
			return evaluadores, outputError
		}

		if len(respuesta_evaluadores) == 0 {
			outputError = fmt.Errorf("No se encontraron evaluadores para la evaluación")
			return evaluadores, outputError
		}

		helpers.LimpiezaRespuestaRefactor(respuesta_evaluadores, &evaluadores_asignacion)

		for _, evaluador := range evaluadores_asignacion {
			var datos_evaluador models.Evaluador
			datos_evaluador.Rol = evaluador.RolAsignacionEvaluadorId.CodigoAbreviacion
			datos_evaluador.Documento = strconv.Itoa(evaluador.PersonaId)
			datos_evaluador.Cargo = evaluador.Cargo
			var respuesta_cambio_estado_asignacion_evaluador map[string]interface{}
			var cambio_estado_asignacion_evaluador []models.CambioEstadoAsignacionEvaluador

			// Obtener el estado de la evaluacion de cada evaluador
			if response, err := helpers.GetJsonTest(beego.AppConfig.String("UrlEvaluacionesCumplidosProveedoresCrud")+"/cambio_estado_asignacion_evaluador/?query=AsignacionEvaluadorId.Id:"+strconv.Itoa(evaluador.Id)+",Activo:true&limit=-1&sortby=FechaCreacion&order=desc", &respuesta_cambio_estado_asignacion_evaluador); err != nil && response != 200 {
				outputError = fmt.Errorf("Error al obtener el cambio de estado de la asignación del evaluador")
				return evaluadores, outputError
			}

			if len(respuesta_cambio_estado_asignacion_evaluador) == 0 {
				datos_evaluador.EstadoEvaluacion = ""
			} else {
				helpers.LimpiezaRespuestaRefactor(respuesta_cambio_estado_asignacion_evaluador, &cambio_estado_asignacion_evaluador)
				datos_evaluador.EstadoEvaluacion = cambio_estado_asignacion_evaluador[0].EstadoAsignacionEvaluadorId.Nombre
			}

			//Obtenemos los items asignados al evaluador
			_, items_evaluador, errorItems := helpers.ObtenerItemsEvaluador(evaluador.Id)
			if errorItems != nil {
				datos_evaluador.ItemsEvaluados = ""
			} else {
				datos_evaluador.ItemsEvaluados = items_evaluador
			}

			// Obtener el resultado evaluacion del evaluador
			var resultado models.Resultado
			resultado_evaluacion, err := helpers.ObtenerResultadoEvaluacion(evaluador.Id)
			if err == nil {
				// Convertir el resultado evaluacion a la estructura Resultado
				error_json := json.Unmarshal([]byte(resultado_evaluacion.ResultadoEvaluacion), &resultado)
				if error_json != nil {
					outputError = fmt.Errorf("Error al convertir el resultado de la evaluación")
					return evaluadores, outputError
				}
				datos_evaluador.Observaciones = resultado_evaluacion.Observaciones
			}

			// Calcular el puntaje total de la evaluacion
			_, puntaje_total_evaluacion, error_clasificacion := helpers.CalcularClasificacionEvaluacion(resultado)
			if error_clasificacion != nil {
				datos_evaluador.PuntajeEvaluacion = 0
			} else {
				datos_evaluador.PuntajeEvaluacion = puntaje_total_evaluacion
			}

			evaluadores = append(evaluadores, datos_evaluador)

		}
	} else {
		var datos_evaluador models.Evaluador
		datos_evaluador.Documento = strconv.Itoa(asignacion_evaluador.PersonaId)
		datos_evaluador.Cargo = asignacion_evaluador.Cargo
		var respuesta_cambio_estado_asignacion_evaluador map[string]interface{}
		var cambio_estado_asignacion_evaluador []models.CambioEstadoAsignacionEvaluador

		// Obtener el estado de la evaluacion de cada evaluador
		if response, err := helpers.GetJsonTest(beego.AppConfig.String("UrlEvaluacionesCumplidosProveedoresCrud")+"/cambio_estado_asignacion_evaluador/?query=AsignacionEvaluadorId.Id:"+strconv.Itoa(asignacion_evaluador.Id)+",Activo:true&limit=-1&sortby=FechaCreacion&order=desc", &respuesta_cambio_estado_asignacion_evaluador); err != nil && response != 200 {
			outputError = fmt.Errorf("Error al obtener el cambio de estado de la asignación del evaluador")
			return evaluadores, outputError
		}

		if len(respuesta_cambio_estado_asignacion_evaluador) == 0 {
			datos_evaluador.EstadoEvaluacion = ""
		} else {
			helpers.LimpiezaRespuestaRefactor(respuesta_cambio_estado_asignacion_evaluador, &cambio_estado_asignacion_evaluador)
			datos_evaluador.EstadoEvaluacion = cambio_estado_asignacion_evaluador[0].EstadoAsignacionEvaluadorId.Nombre
		}

		//Obtenemos los items asignados al evaluador
		_, items_evaluador, errorItems := helpers.ObtenerItemsEvaluador(asignacion_evaluador.Id)
		if errorItems != nil {
			datos_evaluador.ItemsEvaluados = ""
		} else {
			datos_evaluador.ItemsEvaluados = items_evaluador
		}

		// Obtener el resultado evaluacion del evaluador
		var resultado models.Resultado
		resultado_evaluacion, err := helpers.ObtenerResultadoEvaluacion(asignacion_evaluador.Id)
		if err == nil {
			// Convertir el resultado evaluacion a la estructura Resultado
			error_json := json.Unmarshal([]byte(resultado_evaluacion.ResultadoEvaluacion), &resultado)
			if error_json != nil {
				outputError = fmt.Errorf("Error al convertir el resultado de la evaluación")
				return evaluadores, outputError
			}
		}

		// Calcular el puntaje total de la evaluacion
		_, puntaje_total_evaluacion, error_clasificacion := helpers.CalcularClasificacionEvaluacion(resultado)
		if error_clasificacion != nil {
			datos_evaluador.PuntajeEvaluacion = 0
		} else {
			datos_evaluador.PuntajeEvaluacion = puntaje_total_evaluacion
		}

		evaluadores = append(evaluadores, datos_evaluador)
	}
	return evaluadores, nil
}
