package services

import (
	"fmt"
	"strconv"

	"github.com/astaxie/beego"
	"github.com/udistrital/evaluacion_cumplido_prov_mid/helpers"
	"github.com/udistrital/evaluacion_cumplido_prov_mid/models"
)

func SubirEvaluacionCumplido(evaluacion_id string) (carga_cumplidos []models.CumplidoProveedor, outputError error) {
	defer func() {
		if err := recover(); err != nil {
			outputError = fmt.Errorf("%v", err)
			panic(outputError)
		}
	}()

	// Verificar que el id de la evaluacion sea un número mayor a 0
	if number, err := strconv.Atoi(evaluacion_id); err != nil {
		outputError = fmt.Errorf("El id de la evaluacion no es un número válido")
		return carga_cumplidos, outputError
	} else if number <= 0 {
		outputError = fmt.Errorf("El id de la evaluacion debe ser un número mayor a 0")
		return carga_cumplidos, outputError
	}

	// Obtener la evaluacion y verificar que exista
	var resultado_evaluacion map[string]interface{}
	var evaluacion []models.Evaluacion
	//fmt.Println("Url evaluacion: ", beego.AppConfig.String("UrlEvaluacionCumplidoCrud")+"/evaluacion/?query=Id:"+strconv.Itoa(evaluacion_id)+",Activo:true&limit=-1")
	if response, err := helpers.GetJsonTest(beego.AppConfig.String("UrlEvaluacionCumplidoCrud")+"/evaluacion/?query=Id:"+evaluacion_id+",Activo:true&limit=-1", &resultado_evaluacion); err != nil && response != 200 {
		outputError = fmt.Errorf("Error al obtener la evaluacion")
		return carga_cumplidos, outputError
	}

	if len(resultado_evaluacion) == 0 {
		outputError = fmt.Errorf("No se encontró la evaluación")
		return carga_cumplidos, outputError
	}

	helpers.LimpiezaRespuestaRefactor(resultado_evaluacion, &evaluacion)

	// Obtener los cumplidos a los que se les puede subir la evaluacion
	var respuesta_peticion map[string]interface{}
	var cambios_estados []models.CambioEstadoCumplido
	//fmt.Println("URL Cambio Estado: ", beego.AppConfig.String("UrlCrudRevisionCumplidosProveedores")+"/cambio_estado_cumplido/?query=CumplidoProveedorId.NumeroContrato:"+strconv.Itoa(evaluacion[0].ContratoSuscritoId)+",CumplidoProveedorId.VigenciaContrato:"+strconv.Itoa(evaluacion[0].VigenciaContrato)+",EstadoCumplidoId.CodigoAbreviacion.in:CD|RC|RO,Activo:true&sortby=FechaCreacion&order=desc&limit=-1")
	if response, err := helpers.GetJsonTest(beego.AppConfig.String("UrlCrudRevisionCumplidosProveedores")+"/cambio_estado_cumplido/?query=CumplidoProveedorId.NumeroContrato:"+strconv.Itoa(evaluacion[0].ContratoSuscritoId)+",CumplidoProveedorId.VigenciaContrato:"+strconv.Itoa(evaluacion[0].VigenciaContrato)+",EstadoCumplidoId.CodigoAbreviacion.in:CD|RC|RO,Activo:true&sortby=FechaCreacion&order=desc&limit=-1", &respuesta_peticion); err != nil && response != 200 {
		outputError = fmt.Errorf("Error al obtener los cambios de estado del cumplido")
		return carga_cumplidos, outputError
	}

	data := respuesta_peticion["Data"].([]interface{})
	if len(data[0].(map[string]interface{})) == 0 {
		outputError = fmt.Errorf("No se encontro ningun cumplido al cual se le pueda cargar la evaluacion")
		return carga_cumplidos, outputError
	}

	helpers.LimpiezaRespuestaRefactor(respuesta_peticion, &cambios_estados)

	// Validar la informacion de pago de cada uno de los cumplidos para validar que sean de tipo unico o total
	var cumplidos_pago_total_unico []models.CumplidoProveedor
	for _, cumplido := range cambios_estados {
		var respuesta_peticion map[string]interface{}
		var informacion_pago_proveedor []models.InformacionPago
		//fmt.Println("URL Informacion Pago: ", beego.AppConfig.String("UrlCrudRevisionCumplidosProveedores")+"/informacion_pago/?query=Activo:true,CumplidoProveedorId.Id:"+strconv.Itoa(cumplido.CumplidoProveedorId.Id))
		if response, err := helpers.GetJsonTest(beego.AppConfig.String("UrlCrudRevisionCumplidosProveedores")+"/informacion_pago/?query=Activo:true,CumplidoProveedorId.Id:"+strconv.Itoa(cumplido.CumplidoProveedorId.Id), &respuesta_peticion); err != nil && response != 200 {
			outputError = fmt.Errorf("Error al obtener la informacion de pago del cumplido")
			return carga_cumplidos, outputError
		}

		data := respuesta_peticion["Data"].([]interface{})
		if len(data[0].(map[string]interface{})) == 0 {
			continue
		}
		helpers.LimpiezaRespuestaRefactor(respuesta_peticion, &informacion_pago_proveedor)

		// if informacion_pago_proveedor[0].TipoPagoId.CodigoAbreviacion == "TPU" || informacion_pago_proveedor[0].TipoPagoId.CodigoAbreviacion == "TPT" {
		// 	cumplidos_pago_total_unico = append(cumplidos_pago_total_unico, cumplido.CumplidoProveedorId)
		// }

		if informacion_pago_proveedor[0].TipoPagoId.Id == 2 || informacion_pago_proveedor[0].TipoPagoId.Id == 3 {
			cumplidos_pago_total_unico = append(cumplidos_pago_total_unico, cumplido.CumplidoProveedorId)
		}
	}

	// Validaciones finales para saber que retornar
	if len(cumplidos_pago_total_unico) == 0 {
		// Si no hay ningun cumplido que cumpla las condiciones se retorna un error indicando que no se encontraron cumplidos con tipo de pago total o unico
		outputError = fmt.Errorf("No se encontraron cumplidos con tipo de pago total o unico")
		return carga_cumplidos, outputError
	}

	if len(cumplidos_pago_total_unico) == 1 {
		// Si solo hay un cumplido que cumpla las condiciones se sube la evaluacion a ese cumplido

		// Verificar si ya hay cargada una evaluacion al cumplido, en caso de que si, se elimina la evaluacion para cargar la nueva
		if err := EliminarEvaluacionCumplido(cumplidos_pago_total_unico[0].Id); err != nil {
			outputError = fmt.Errorf("Error al eliminar la evaluacion anterior del cumplido")
			return carga_cumplidos, outputError
		}

		var body_subir_soporte models.SubirSoporteCumplido
		// Subir el documento de evaluacion al cumplido
		if evaluacion[0].DocumentoId == 0 {
			outputError = fmt.Errorf("La evaluacion no tiene documento asociado")
			return carga_cumplidos, outputError
		}
		body_subir_soporte.CumplidoProveedorId.Id = cumplidos_pago_total_unico[0].Id
		body_subir_soporte.DocumentoId = evaluacion[0].DocumentoId

		var respuesta_peticion_subir_soporte map[string]interface{}
		if err := helpers.SendJson(beego.AppConfig.String("UrlCrudRevisionCumplidosProveedores")+"/soporte_cumplido", "POST", &respuesta_peticion_subir_soporte, body_subir_soporte); err != nil {
			outputError = fmt.Errorf("Error al subir el soporte al cumplido")
			return carga_cumplidos, outputError
		}

		carga_cumplidos = append(carga_cumplidos, cumplidos_pago_total_unico[0])

	}

	if len(cumplidos_pago_total_unico) > 1 {
		// Si hay mas de un cumplido se retorna un error indicando que hay mas de un cumplido que cumple las condiciones y se retornan los cumplidos que cumplen las condiciones
		outputError = fmt.Errorf("Hay mas de un cumplido que cumple las condiciones, por este motivo no se pudo cargar automaticamente la evaluacion")
		carga_cumplidos = cumplidos_pago_total_unico
		return carga_cumplidos, outputError
	}

	return carga_cumplidos, outputError

}

func EliminarEvaluacionCumplido(cumplido_proveedor_id int) (outputError error) {
	defer func() {
		if err := recover(); err != nil {
			outputError = fmt.Errorf("%v", err)
			panic(outputError)
		}
	}()

	// Verificar que el id del cumplido sea un número mayor a 0
	if cumplido_proveedor_id <= 0 {
		outputError = fmt.Errorf("El id del cumplido debe ser un número mayor a 0")
		return outputError
	}

	var respuesta_documentos_cumplido map[string]interface{}
	var soportes_cumplido []models.DocumentosSoporteSimplificado

	// Obtener los documentos del cumplido
	if response, err := helpers.GetJsonTest(beego.AppConfig.String("UrlMidRevisionCumplidosProveedores")+"/solicitud-pago/soportes/"+strconv.Itoa(cumplido_proveedor_id), &respuesta_documentos_cumplido); err != nil && response != 200 {
		outputError = fmt.Errorf("Error al obtener los documentos del cumplido")
		return outputError
	}

	data := respuesta_documentos_cumplido["Data"].([]interface{})
	if len(data[0].(map[string]interface{})) == 0 {
		outputError = fmt.Errorf("No se encontraron documentos del cumplido")
		return outputError
	}

	helpers.LimpiezaRespuestaRefactor(respuesta_documentos_cumplido, &soportes_cumplido)

	// Verificar si hay una evaluacion subida al cumplido
	for _, soporte := range soportes_cumplido {
		if soporte.Documento.CodigoAbreviacionTipoDocumento == "EP" {
			// Eliminar la evaluacion
			var respuesta_peticion map[string]interface{}
			if response, err := helpers.GetJsonTest(beego.AppConfig.String("UrlCrudRevisionCumplidosProveedores")+"/soporte_cumplido/"+strconv.Itoa(soporte.SoporteCumplidoId), &respuesta_peticion); err != nil && response != 200 {
				outputError = fmt.Errorf("Error al eliminar la evaluacion")
				return outputError
			}
		}

		var respuesta_peticion map[string]interface{}
		data := make(map[string]interface{})

		if err := helpers.SendJson(beego.AppConfig.String("UrlCrudRevisionCumplidosProveedores")+"/soporte_cumplido/"+strconv.Itoa(soporte.SoporteCumplidoId), "DELETE", &respuesta_peticion, data); err != nil {
			outputError = fmt.Errorf("Error al eliminar el soporte del cumplido")
			return outputError
		}

		return outputError
	}
	return outputError

}
