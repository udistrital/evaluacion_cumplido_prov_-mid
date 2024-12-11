package helpers

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/astaxie/beego"
	"github.com/udistrital/evaluacion_cumplido_prov_mid/models"
)

// Funcion para calcular la clasificacion de la evaluacion, por el momento este dato se calcula desde el cliente
func CalcularClasificacionEvaluacion(resultados_finales models.Resultado) (resultado_clasificacion models.Clasificacion, puntaje_total_evaluacion int, outputError error) {
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

	if response, err := GetJsonTest(beego.AppConfig.String("UrlEvaluacionesCumplidosProveedoresCrud")+"/clasificacion/?query=CodigoAbreviacion:"+codigo_abreviacion_clasificacion+",Activo:true&limit=1", &respuesta_clasificacion); err != nil && response != 200 {
		outputError = fmt.Errorf("Error al obtener la clasificación de la evaluación")
		return resultado_clasificacion, 0, outputError
	}

	LimpiezaRespuestaRefactor(respuesta_clasificacion, &clasificacion)

	resultado_clasificacion = clasificacion[0]

	return resultado_clasificacion, valor_clasificacion, nil
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
	if response, err := GetJsonTest(beego.AppConfig.String("UrlEvaluacionesCumplidosProveedoresCrud")+"/resultado_evaluacion/?query=AsignacionEvaluadorId.Id:"+strconv.Itoa(asignacion_evaluacion_id)+",Activo:true&limit=1", &respuesta_resultado_evaluacion); err != nil && response != 200 {
		outputError = fmt.Errorf("Error al obtener el resultado de la evaluación")
		return resultado_evaluacion, outputError
	}

	data := respuesta_resultado_evaluacion["Data"].([]interface{})
	if len(data[0].(map[string]interface{})) == 0 {
		outputError = fmt.Errorf("El Evaluador no tiene registrado un resultado de la evaluación")
		return resultado_evaluacion, outputError
	}

	LimpiezaRespuestaRefactor(respuesta_resultado_evaluacion, &resultado)

	resultado_evaluacion = resultado[0]

	return resultado_evaluacion, nil
}

func ObtenerItemsEvaluador(asignacion_evaluador_id int) (items_evaluador []models.Item, items_evaluador_formateados string, outputError error) {
	defer func() {
		if err := recover(); err != nil {
			outputError = fmt.Errorf("%v", err)
			panic(outputError.Error())
		}
	}()

	var respuesta_items_evaluador map[string]interface{}
	var items []models.AsignacionEvaluadorItem

	if response, err := GetJsonTest(beego.AppConfig.String("UrlEvaluacionesCumplidosProveedoresCrud")+"/asignacion_evaluador_item/?query=AsignacionEvaluadorId.Id:"+strconv.Itoa(asignacion_evaluador_id)+",Activo:true&limit=-1", &respuesta_items_evaluador); err != nil && response != 200 {
		outputError = fmt.Errorf("Error al obtener los items del evaluador")
		return items_evaluador, items_evaluador_formateados, outputError
	}

	data := respuesta_items_evaluador["Data"].([]interface{})
	if len(data[0].(map[string]interface{})) == 0 {
		outputError = fmt.Errorf("El Evaluador no tiene items asignados")
		return items_evaluador, items_evaluador_formateados, outputError
	}

	LimpiezaRespuestaRefactor(respuesta_items_evaluador, &items)

	for _, item := range items {
		items_evaluador = append(items_evaluador, item.ItemId)
		items_evaluador_formateados += strconv.Itoa(item.Id) + ", "
	}

	return items_evaluador, strings.TrimSuffix(items_evaluador_formateados, ", "), nil
}

func ObtenerItemsEvaluacion(evaluacion_id int) (items_evaluacion []models.Item, items_evaluador_formateados string, outputError error) {
	defer func() {
		if err := recover(); err != nil {
			outputError = fmt.Errorf("%v", err)
			panic(outputError.Error())
		}
	}()

	var respuesta_items_evaluacion map[string]interface{}

	//fmt.Println("URL items evaluacion: ", beego.AppConfig.String("UrlEvaluacionesCumplidosProveedoresCrud")+"/item/?query=EvaluacionId.Id:"+strconv.Itoa(evaluacion_id)+",Activo:true&limit=-1")
	if response, err := GetJsonTest(beego.AppConfig.String("UrlEvaluacionesCumplidosProveedoresCrud")+"/item/?query=EvaluacionId.Id:"+strconv.Itoa(evaluacion_id)+",Activo:true&limit=-1", &respuesta_items_evaluacion); err != nil && response != 200 {
		outputError = fmt.Errorf("Error al obtener los items de la evaluación")
		return items_evaluacion, items_evaluador_formateados, outputError
	}

	if len(respuesta_items_evaluacion) == 0 {
		outputError = fmt.Errorf("La evaluación no tiene items asignados")
		return items_evaluacion, items_evaluador_formateados, outputError
	}

	LimpiezaRespuestaRefactor(respuesta_items_evaluacion, &items_evaluacion)

	for _, item := range items_evaluacion {
		items_evaluador_formateados += strconv.Itoa(item.Id) + ", "
	}

	return items_evaluacion, strings.TrimSuffix(items_evaluador_formateados, ", "), nil
}
