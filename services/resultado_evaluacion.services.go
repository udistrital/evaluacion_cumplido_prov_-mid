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

// Funcion para calcular la clasificacion de la evaluacion, por el momento este dato se calcula desde el cliente
func CalcularClasificacionEvaluacion(resultados_finales models.Resultado) (resultado_clasificacion models.Clasificacion, outputError error) {
	defer func() {
		if err := recover(); err != nil {
			outputError = fmt.Errorf("%v", err)
			panic(outputError.Error())
		}
	}()

	var valor_clasificacion int
	var codigo_abreviacion_clasificacion string

	// Calcular la clasificación de la evaluación
	for _, item := range resultados_finales.ResultadosIndividuales {
		valor_clasificacion += item.Respuesta.ValorAsignado
	}

	switch {
	case valor_clasificacion >= 0 && valor_clasificacion <= 45:
		codigo_abreviacion_clasificacion = "ML"
	case valor_clasificacion >= 46 && valor_clasificacion <= 79:
		codigo_abreviacion_clasificacion = "BN"
	case valor_clasificacion >= 80 && valor_clasificacion <= 100:
		codigo_abreviacion_clasificacion = "EX"
	}

	var respuesta_clasificacion map[string]interface{}
	var clasificacion []models.Clasificacion

	if response, err := helpers.GetJsonTest(beego.AppConfig.String("UrlEvaluacionesCumplidosProveedoresCrud")+"/clasificacion/?query=CodigoAbreviacion:"+codigo_abreviacion_clasificacion+",Activo:true&limit=1", &respuesta_clasificacion); err != nil && response != 200 {
		outputError = fmt.Errorf("Error al obtener la clasificación de la evaluación")
		return resultado_clasificacion, outputError
	}

	helpers.LimpiezaRespuestaRefactor(respuesta_clasificacion, &clasificacion)

	resultado_clasificacion = clasificacion[0]

	return resultado_clasificacion, nil
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

func ObtenerResultadoEvaluacion(asignacion_evaluacion_id int) (resultado_evaluacion models.ResultadoEvaluacion, outputError error) {
	defer func() {
		if err := recover(); err != nil {
			outputError = fmt.Errorf("%v", err)
			panic(outputError.Error())
		}
	}()

	var respuesta_resultado_evaluacion map[string]interface{}
	var resultado []models.ResultadoEvaluacion

	//fmt.Println("URL resultado evaluacion: ", beego.AppConfig.String("UrlEvaluacionesCumplidosProveedoresCrud")+"/resultado_evaluacion/?query=AsignacionEvaluadorId.Id:"+strconv.Itoa(asignacion_evaluacion_id)+",Activo:true&limit=1")
	if response, err := helpers.GetJsonTest(beego.AppConfig.String("UrlEvaluacionesCumplidosProveedoresCrud")+"/resultado_evaluacion/?query=AsignacionEvaluadorId.Id:"+strconv.Itoa(asignacion_evaluacion_id)+",Activo:true&limit=1", &respuesta_resultado_evaluacion); err != nil && response != 200 {
		outputError = fmt.Errorf("Error al obtener el resultado de la evaluación")
		return resultado_evaluacion, outputError
	}

	data := respuesta_resultado_evaluacion["Data"].([]interface{})
	if len(data[0].(map[string]interface{})) == 0 {
		outputError = fmt.Errorf("El Evaluador no tiene registrado un resultado de la evaluación")
		return resultado_evaluacion, outputError
	}

	helpers.LimpiezaRespuestaRefactor(respuesta_resultado_evaluacion, &resultado)

	resultado_evaluacion = resultado[0]

	return resultado_evaluacion, nil
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
		nombre_persona, error := ObtenerNombrePersonaNatural(strconv.Itoa(evaluador.PersonaId))
		if error != nil {
			outputError = fmt.Errorf(error.Error())
			return resultados_finales, outputError
		}
		items_evaluador, errorItems := ObtenerItemsEvaluador(evaluador.Id)
		if errorItems != nil {
			outputError = fmt.Errorf(errorItems.Error())
			return resultados_finales, outputError
		}
		resultado_evaluacion, err := ObtenerResultadoEvaluacion(evaluador.Id)
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
			pregunta := item.Respuesta.Pregunta
			cumplimiento := item.Respuesta.Cumplimiento
			valorAsignado := item.Respuesta.ValorAsignado

			// Inicializar la pregunta si no existe
			if _, exists := resultadosPorPregunta[pregunta]; !exists {

				if strings.ToUpper(item.Categoria) == "GESTIÓN" && strings.ToUpper(item.Titulo) == "PROCEDIMIENTOS" {
					resultadosPorPregunta[pregunta] = map[string]float64{"Excelente": 0.0, "Bueno": 0.0, "Regular": 0.0, "Malo": 0.0}
					valorAsignadoPorPregunta[pregunta] = map[string]int{"Excelente": 0, "Bueno": 0, "Regular": 0, "Malo": 0}
					categoriaTituloPorPregunta[pregunta] = map[string]string{"Categoria": item.Categoria, "Titulo": item.Titulo}
				} else {
					resultadosPorPregunta[pregunta] = map[string]float64{"Si": 0.0, "No": 0.0}
					valorAsignadoPorPregunta[pregunta] = map[string]int{"Si": 0, "No": 0}
					categoriaTituloPorPregunta[pregunta] = map[string]string{"Categoria": item.Categoria, "Titulo": item.Titulo}
				}
			}

			// Sumar el porcentaje de evaluación a la opción correspondiente
			resultadosPorPregunta[pregunta][cumplimiento] += evaluador.PorcentajeEvaluacion

			// Registrar el valor asignado para el cumplimiento
			valorAsignadoPorPregunta[pregunta][cumplimiento] = valorAsignado

			//Registar las respuestas del supervisor
			if evaluador.RolAsignacionEvaluadorId.CodigoAbreviacion == "SPR" {
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

func ObtenerNombrePersonaNatural(documento_persona string) (nombre_persona string, outputError error) {
	defer func() {
		if err := recover(); err != nil {
			outputError = fmt.Errorf("%v", err)
			panic(outputError.Error())
		}
	}()

	var informacion []models.InformacionPersonaNatural

	if response, err := helpers.GetJsonTest(beego.AppConfig.String("UrlAmazonApi")+"/informacion_persona_natural/?fields=PrimerNombre,SegundoNombre,PrimerApellido,SegundoApellido&limit=0&query=Id:"+documento_persona, &informacion); err != nil && response != 200 {
		outputError = fmt.Errorf("Error al obtener la información de la persona")
		return nombre_persona, outputError
	}

	nombre_persona = informacion[0].PrimerNombre + " " + informacion[0].SegundoNombre + " " + informacion[0].PrimerApellido + " " + informacion[0].SegundoApellido

	return nombre_persona, nil
}

func ObtenerItemsEvaluador(asignacion_evaluador_id int) (items_evaluador string, outputError error) {
	defer func() {
		if err := recover(); err != nil {
			outputError = fmt.Errorf("%v", err)
			panic(outputError.Error())
		}
	}()

	var respuesta_items_evaluador map[string]interface{}
	var items []models.AsignacionEvaluadorItem

	if response, err := helpers.GetJsonTest(beego.AppConfig.String("UrlEvaluacionesCumplidosProveedoresCrud")+"/asignacion_evaluador_item/?query=AsignacionEvaluadorId.Id:"+strconv.Itoa(asignacion_evaluador_id)+",Activo:true&limit=-1", &respuesta_items_evaluador); err != nil && response != 200 {
		outputError = fmt.Errorf("Error al obtener los items del evaluador")
		return items_evaluador, outputError
	}

	helpers.LimpiezaRespuestaRefactor(respuesta_items_evaluador, &items)

	for _, item := range items {
		items_evaluador += strconv.Itoa(item.Id) + ", "
	}

	return strings.TrimSuffix(items_evaluador, ", "), nil
}
