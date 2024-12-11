package services

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

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

func ObtenerResultadoFinalEvaluacion(evaluacion_id int) (resultados_finales models.ResultadoFinalEvaluacion, outputError error) {
	defer func() {
		if err := recover(); err != nil {
			outputError = fmt.Errorf("%v", err)
			panic(outputError.Error())
		}
	}()

	// Validar que la evaluacion este en estado Pendiente Revision Evaluadores

	var respuesta_cambio_estado map[string]interface{}
	//fmt.Println("URL: ", beego.AppConfig.String("UrlEvaluacionesCumplidosProveedoresCrud")+"/cambio_estado_evaluacion/?query=EvaluacionId.Id:"+strconv.Itoa(evaluacion_id)+",EstadoEvaluacionId.CodigoAbreviacion:AE,Activo:true&limit=-1&sortby=FechaCreacion&order=desc")

	if response, err := helpers.GetJsonTest(beego.AppConfig.String("UrlEvaluacionesCumplidosProveedoresCrud")+"/cambio_estado_evaluacion/?query=EvaluacionId.Id:"+strconv.Itoa(evaluacion_id)+",EstadoEvaluacionId.CodigoAbreviacion:PRE,Activo:true&limit=-1&sortby=FechaCreacion&order=desc", &respuesta_cambio_estado); err != nil && response != 200 {
		outputError = fmt.Errorf("Error al obtener el cambio de estado de la evaluación")
		return resultados_finales, outputError
	}

	data := respuesta_cambio_estado["Data"].([]interface{})
	if len(data[0].(map[string]interface{})) == 0 {
		outputError = fmt.Errorf("La evaluación no se encuentra en estado Aprobado evaluadores")
		return resultados_finales, outputError
	}

	// Obtener todos los evaluadores de la evaluacion
	var respuesta_evaluadores map[string]interface{}
	var evaluadores []models.AsignacionEvaluador

	//fmt.Println("URL evaluadores: ", beego.AppConfig.String("UrlEvaluacionesCumplidosProveedoresCrud")+"/asignacion_evaluador/?query=EvaluacionId.Id:"+strconv.Itoa(evaluacion_id)+",Activo:true&limit=-1")
	if response, err := helpers.GetJsonTest(beego.AppConfig.String("UrlEvaluacionesCumplidosProveedoresCrud")+"/asignacion_evaluador/?query=EvaluacionId.Id:"+strconv.Itoa(evaluacion_id)+",Activo:true&limit=-1", &respuesta_evaluadores); err != nil && response != 200 {
		outputError = fmt.Errorf("Error al obtener los evaluadores de la evaluación")
		return resultados_finales, outputError
	}

	helpers.LimpiezaRespuestaRefactor(respuesta_evaluadores, &evaluadores)

	resultados_finales, err := ProcesarResultadosEvaluaciones(evaluadores)
	if err != nil {
		outputError = fmt.Errorf(err.Error())
		return resultados_finales, outputError
	}

	return resultados_finales, nil
}

func ProcesarResultadosEvaluaciones(evaluadores []models.AsignacionEvaluador) (resultados_finales models.ResultadoFinalEvaluacion, outputError error) {
	defer func() {
		if err := recover(); err != nil {
			outputError = fmt.Errorf("%v", err)
			panic(outputError.Error())
		}
	}()

	resultadosPorPregunta := make(map[string]map[string]float64)
	valorAsignadoPorPregunta := make(map[string]map[string]int)
	categoriaTituloPorPregunta := make(map[string]map[string]string)
	supervisorRespuestas := make(map[string]string)

	// Obtener el resultado de la evaluacion de cada evaluador

	for _, evaluador := range evaluadores {
		nombre_persona, error := helpers.ObtenerNombrePersonaNatural(strconv.Itoa(evaluador.PersonaId))
		if error != nil {
			outputError = fmt.Errorf(error.Error())
			return resultados_finales, outputError
		}
		_, items_evaluador, errorItems := helpers.ObtenerItemsEvaluador(evaluador.Id)
		if errorItems != nil {
			outputError = fmt.Errorf(errorItems.Error())
			return resultados_finales, outputError
		}
		resultado_evaluacion, err := helpers.ObtenerResultadoEvaluacion(evaluador.Id)
		if err != nil {
			outputError = fmt.Errorf(err.Error())
			return resultados_finales, outputError
		}

		// Convertir el resultado evaluacion a la estructura Resultado

		var resultado models.Resultado
		error_json := json.Unmarshal([]byte(resultado_evaluacion.ResultadoEvaluacion), &resultado)
		if error_json != nil {
			outputError = fmt.Errorf("Error al convertir el resultado de la evaluación")
			return resultados_finales, outputError
		}

		for _, item := range resultado.ResultadosIndividuales {
			pregunta := strings.TrimSpace(strings.ToUpper(item.Respuesta.Pregunta))
			cumplimiento := strings.TrimSpace(strings.ToUpper(item.Respuesta.Cumplimiento))
			valorAsignado := item.Respuesta.ValorAsignado

			// Inicializar la pregunta si no existe
			if _, exists := resultadosPorPregunta[pregunta]; !exists {

				if strings.ToUpper(item.Categoria) == "GESTIÓN" && strings.ToUpper(item.Titulo) == "PROCEDIMIENTOS" {
					resultadosPorPregunta[pregunta] = map[string]float64{"EXCELENTE": 0.0, "BUENO": 0.0, "REGULAR": 0.0, "MALO": 0.0}
					valorAsignadoPorPregunta[pregunta] = map[string]int{"EXCELENTE": 0, "BUENO": 0, "REGULAR": 0, "MALO": 0}
					categoriaTituloPorPregunta[pregunta] = map[string]string{"Categoria": strings.ToUpper(item.Categoria), "Titulo": strings.ToUpper(item.Titulo)}
				} else {
					resultadosPorPregunta[pregunta] = map[string]float64{"SI": 0.0, "NO": 0.0}
					valorAsignadoPorPregunta[pregunta] = map[string]int{"SI": 0, "NO": 0}
					categoriaTituloPorPregunta[pregunta] = map[string]string{"Categoria": strings.ToUpper(item.Categoria), "Titulo": strings.ToUpper(item.Titulo)}
				}
			}

			// Sumar el porcentaje de evaluación a la opción correspondiente
			resultadosPorPregunta[pregunta][cumplimiento] += evaluador.PorcentajeEvaluacion

			// Registrar el valor asignado para el cumplimiento
			valorAsignadoPorPregunta[pregunta][cumplimiento] = valorAsignado

			//Registar las respuestas del supervisor
			if evaluador.RolAsignacionEvaluadorId.CodigoAbreviacion == "SP" {
				supervisorRespuestas[pregunta] = cumplimiento
			}

		}

		// Registrar el evaluador, su cargo y sus items
		resultados_finales.Evaluadores = append(resultados_finales.Evaluadores, struct {
			Nombre string
			Cargo  string
			Items  string
			Rol    string
		}{nombre_persona, evaluador.Cargo, items_evaluador, evaluador.RolAsignacionEvaluadorId.Nombre})

	}

	//Determinar la respuesta Final
	// Determinar la respuesta final
	for pregunta, votos := range resultadosPorPregunta {
		var respuestaFinal string
		var mayorPuntaje float64
		empatados := []string{}
		// Evaluar las opciones disponibles
		for opcion, puntaje := range votos {
			if puntaje > mayorPuntaje {
				// Nueva opción con mayor puntaje
				mayorPuntaje = puntaje
				respuestaFinal = opcion
				empatados = []string{opcion} // Reinicia la lista de empatados
			} else if puntaje == mayorPuntaje {
				// Añadir a la lista de empatados
				empatados = append(empatados, opcion)
			}
		}

		// Resolver empates si hay más de una opción empatada
		if len(empatados) > 1 {
			for _, opcion := range empatados {
				if opcion == supervisorRespuestas[pregunta] {
					respuestaFinal = opcion
					break
				}
			}
		}

		// Obtener el valor asignado de la respuesta ganadora
		valorGanador := valorAsignadoPorPregunta[pregunta][respuestaFinal]

		// Guardar los resultados en la estructura final
		resultados_finales.Resultados = append(resultados_finales.Resultados, struct {
			Categoria     string
			Titulo        string
			Pregunta      string
			Cumplimiento  string
			ValorAsignado int
		}{
			Categoria:     categoriaTituloPorPregunta[pregunta]["Categoria"],
			Titulo:        categoriaTituloPorPregunta[pregunta]["Titulo"],
			Pregunta:      pregunta,
			Cumplimiento:  strings.ToUpper(respuestaFinal),
			ValorAsignado: valorGanador,
		})
	}

	return resultados_finales, nil
}
