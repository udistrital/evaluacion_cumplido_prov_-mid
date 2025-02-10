package services

import (
	"fmt"
	"strconv"

	"github.com/astaxie/beego"
	"github.com/udistrital/evaluacion_cumplido_prov_mid/helpers"
	"github.com/udistrital/evaluacion_cumplido_prov_mid/models"
)

func ObtenerListaDeAsignaciones(documento string) (mapResponse map[string]interface{}, outputError error) {

	defer func() {
		if err := recover(); err != nil {
			outputError = fmt.Errorf("%v", err)
			panic(outputError)
		}
	}()

	mapResponse = make(map[string]interface{})
	var listaAsignacionEvaluador []models.AsignacionEvaluador
	var listaConsultasAsignaciones []models.AsignacionEvaluacion
	listaAsignacionEvaluador, err := consultarAsignaciones(documento)
	var listaAsignaciones []models.AsignacionEvaluacion

	if err != nil {
		return nil, fmt.Errorf("Error al consultar asignaciones")
	}

	if len(listaAsignacionEvaluador) > 0 && listaAsignacionEvaluador[0].EvaluacionId != nil {
		for _, asignacion := range listaAsignacionEvaluador {

			listaContratoGeneral, err := obtenerContratoGeneral(asignacion.EvaluacionId.ContratoSuscritoId, asignacion.EvaluacionId.VigenciaContrato)

			if err != nil {
				return nil, fmt.Errorf("Error al consultar detalles del contrato")

			}

			if len(listaContratoGeneral) > 0 {
				var contratoGeneral = listaContratoGeneral[0]

				respuesta, err := obtenerProveedor(contratoGeneral.Contratista, asignacion, listaContratoGeneral)
				listaConsultasAsignaciones = append(listaConsultasAsignaciones, respuesta)
				if err != nil {
					return nil, fmt.Errorf("Error al consultar detalles del contrato")

				}

			}

		}

	}

	dependencias, err := obtenerDependencias(documento)
	if err != nil {
		return nil, fmt.Errorf("Error al consultar asignaciones")
	}

	var listaContratosSupervisor []models.Contrato

	for _, dependencia := range dependencias {
		contrato_dependencia, _ := consultarContratosPorDependencia(dependencia.Codigo)
		listaContratosSupervisor = append(listaContratosSupervisor, contrato_dependencia...)
	}

	var contratosDepedencia []models.ContratoGeneral

	for _, contrato := range listaContratosSupervisor {

		numeroContrato, _ := strconv.Atoi(contrato.NumeroContrato)
		numeroVigencia, _ := strconv.Atoi(contrato.Vigencia)

		contratoGeneral, _ := obtenerContratoGeneral(numeroContrato, numeroVigencia)
		contratosDepedencia = append(contratosDepedencia, contratoGeneral...)
	}

	listaSinAsignaciones, err := consulartasingAsingnaciones(contratosDepedencia)

	if err != nil {
		return nil, fmt.Errorf("error al consultar asignaciones")
	}
	mapResponse["SinAsignaciones"] = limpiarSinAsignaciones(listaConsultasAsignaciones, listaSinAsignaciones)

	for _, asignacion := range listaConsultasAsignaciones {

		listaAsignaciones = append(listaAsignaciones, asignacion)

	}

	mapResponse["Asignaciones"] = listaAsignaciones

	return mapResponse, nil
}

func consulartasingAsingnaciones(contratosDepedencia []models.ContratoGeneral) (listaSinAsignaciones []models.AsignacionEvaluacion, outputError error) {

	for _, contrato := range contratosDepedencia {

		nombre_dependencia, _ := ObtenerDependencia(contrato.Supervisor.DependenciaSupervisor)

		var listaProveedor []models.InformacionProveedor
		if response, err := helpers.GetJsonWSO2Test(beego.AppConfig.String("UrlAdministrativaAmazonApi")+"/informacion_proveedor/?query=Id:"+strconv.Itoa(contrato.Contratista), &listaProveedor); err == nil && response == 200 {
			asignacionEvaluacion := models.AsignacionEvaluacion{
				AsignacionEvaluacionId: 0,
				NombreProveedor:        listaProveedor[0].NomProveedor,
				Dependencia:            nombre_dependencia,
				TipoContrato:           contrato.TipoContrato.TipoContrato,
				NumeroContrato:         contrato.ContratoSuscrito[0].NumeroContratoSuscrito,
				VigenciaContrato:       strconv.Itoa(contrato.VigenciaContrato),
				EvaluacionId:           0,
			}
			listaSinAsignaciones = append(listaSinAsignaciones, asignacionEvaluacion)

		} else {
			return nil, fmt.Errorf("Error al consultar asignaciones")

		}
	}
	return listaSinAsignaciones, nil
}

func consultarAsignaciones(documento string) (asignaciones []models.AsignacionEvaluador, outputError error) {
	defer func() {
		if err := recover(); err != nil {
			outputError = fmt.Errorf("%v", err)
			panic(outputError)
		}
	}()
	var respuestaPeticion map[string]interface{}
	var listaAsignacionEvaluador []models.AsignacionEvaluador
	fmt.Println(beego.AppConfig.String("UrlEvaluacionCumplidoCrud") + "/asignacion_evaluador?query=personaId:" + documento)

	if response, err := helpers.GetJsonWSO2Test(beego.AppConfig.String("UrlEvaluacionCumplidoCrud")+"/asignacion_evaluador?query=personaId:"+documento, &respuestaPeticion); err == nil && response == 200 {
		helpers.LimpiezaRespuestaRefactor(respuestaPeticion, &listaAsignacionEvaluador)
		if len(listaAsignacionEvaluador) > 0 && listaAsignacionEvaluador[0].EvaluacionId != nil {
			asignaciones = listaAsignacionEvaluador

		}
	} else {
		return asignaciones, fmt.Errorf("Error al consultar asignaciones")

	}
	return asignaciones, nil
}

func obtenerContratoGeneral(contratoSuscritoId int, vigenciaContrato int) (contratoGeneral []models.ContratoGeneral, outputError error) {
	defer func() {
		if err := recover(); err != nil {
			outputError = fmt.Errorf("%v", err)
			panic(outputError)
		}
	}()

	if response, err := helpers.GetJsonWSO2Test(beego.AppConfig.String("UrlAdministrativaAmazonApi")+"/contrato_general/?query=ContratoSuscrito.NumeroContratoSuscrito:"+strconv.Itoa(contratoSuscritoId)+",VigenciaContrato:"+strconv.Itoa(vigenciaContrato), &contratoGeneral); err == nil && response == 200 {
	} else {
		return contratoGeneral, fmt.Errorf("Error al consultar asignaciones")

	}
	return contratoGeneral, nil
}

func obtenerContratoGeneralPorNumeroDecontrato(contratoSuscritoId int, vigenciaContrato int) (contratoGeneral []models.ContratoGeneral, outputError error) {
	defer func() {
		if err := recover(); err != nil {
			outputError = fmt.Errorf("%v", err)
			panic(outputError)
		}
	}()

	if response, err := helpers.GetJsonWSO2Test(beego.AppConfig.String("UrlAdministrativaAmazonApi")+"/contrato_general/?query=ContratoSuscrito.Id:"+strconv.Itoa(contratoSuscritoId)+",VigenciaContrato:"+strconv.Itoa(vigenciaContrato), &contratoGeneral); err == nil && response == 200 {
	} else {
		return contratoGeneral, fmt.Errorf("Error al consultar asignaciones")

	}
	return contratoGeneral, nil
}

func obtenerProveedor(contratistaId int, asignacion models.AsignacionEvaluador, listaContratoGeneral []models.ContratoGeneral) (asisgnaciones models.AsignacionEvaluacion, outputError error) {
	defer func() {
		if err := recover(); err != nil {
			outputError = fmt.Errorf("%v", err)
			panic(outputError)
		}
	}()

	var listaProveedor []models.InformacionProveedor
	contratoGeneral := listaContratoGeneral[0]
	if response, err := helpers.GetJsonWSO2Test(beego.AppConfig.String("UrlAdministrativaAmazonApi")+"/informacion_proveedor/?query=Id:"+strconv.Itoa(contratistaId), &listaProveedor); err == nil && response == 200 {
		estado, err := obtenerEstadoAsignacionEvaluacion(asignacion.Id)

		if err != nil {
			return asisgnaciones, fmt.Errorf("Error al consultar estado de asignaciones")

		}

		nombre_dependencia, _ := ObtenerDependencia(contratoGeneral.Supervisor.DependenciaSupervisor)
		estadoEvaluacion, _ := ObtenerEstadoEvaluacion(asignacion.EvaluacionId.Id)
		asignacionEvaluacion := models.AsignacionEvaluacion{
			AsignacionEvaluacionId:    asignacion.Id,
			NombreProveedor:           listaProveedor[0].NomProveedor,
			Dependencia:               nombre_dependencia,
			TipoContrato:              contratoGeneral.TipoContrato.TipoContrato,
			NumeroContrato:            contratoGeneral.ContratoSuscrito[0].NumeroContratoSuscrito,
			VigenciaContrato:          strconv.Itoa(contratoGeneral.VigenciaContrato),
			EvaluacionId:              asignacion.EvaluacionId.Id,
			EstadoAsignacionEvaluador: estado[0].EstadoAsignacionEvaluador,
			EstadoEvaluacion:          &estadoEvaluacion,
			RolEvaluador:              asignacion.RolAsignacionEvaluadorId.CodigoAbreviacion,
		}
		asisgnaciones = asignacionEvaluacion
	} else {
		return asisgnaciones, fmt.Errorf("Error al consultar asignaciones")

	}

	return asisgnaciones, nil
}

func obtenerDependencias(documento string) (dependencias []models.Dependencia, outputError error) {

	defer func() {
		if err := recover(); err != nil {
			outputError = fmt.Errorf("%v", err)
			panic(outputError)
		}
	}()

	var respuesta models.DependenciasRespuesta
	fmt.Println(beego.AppConfig.String("UrlAdministrativaJBPM") + "/dependencias_supervisor/" + documento)
	if response, err := helpers.GetJsonWSO2Test(beego.AppConfig.String("UrlAdministrativaJBPM")+"/dependencias_supervisor/"+documento, &respuesta); err == nil && response == 200 {
		dependencias = respuesta.Dependencias.Dependencia
	} else {
		return dependencias, fmt.Errorf("Error al consultar asignaciones")
	}
	return dependencias, nil
}

func consultarContratosPorDependencia(dependencia string) (contratos []models.Contrato, outputError error) {

	defer func() {
		if err := recover(); err != nil {
			outputError = fmt.Errorf("%v", err)
			panic(outputError)
		}
	}()

	var respuesta models.ContratosRespuesta

	if response, err := helpers.GetJsonWSO2Test(beego.AppConfig.String("UrlAdministrativaJBPM")+"/contratos_proveedor_dependencia/"+dependencia, &respuesta); err == nil && response == 200 {
		contratos = respuesta.Contratos.Contrato
	} else {
		return contratos, fmt.Errorf("Error al consultar depenendiencias")
	}
	return contratos, nil

}

func limpiarSinAsignaciones(Asignaciones, SinAsignaciones []models.AsignacionEvaluacion) []models.AsignacionEvaluacion {
	asignacionesMap := make(map[string]bool)
	for _, a := range Asignaciones {
		key := fmt.Sprintf("%s-%s", a.NumeroContrato, a.VigenciaContrato)
		asignacionesMap[key] = true
	}

	var filtroSinAsignaciones []models.AsignacionEvaluacion
	for _, sa := range SinAsignaciones {
		key := fmt.Sprintf("%s-%s", sa.NumeroContrato, sa.VigenciaContrato)
		if !asignacionesMap[key] {
			filtroSinAsignaciones = append(filtroSinAsignaciones, sa)
		}
	}

	return filtroSinAsignaciones
}

func obtenerEstadoAsignacionEvaluacion(AsignacionId int) (listaCambiosEstado []models.CambioEstadoASignacionEnvaluacion, outputError error) {
	defer func() {
		if err := recover(); err != nil {
			outputError = fmt.Errorf("%v", err)
			panic(outputError)
		}
	}()
	var query = fmt.Sprintf("/cambio_estado_asignacion_evaluador/?query=AsignacionEvaluadorId.Id:%d,Activo:true", AsignacionId)

	fmt.Println(beego.AppConfig.String("UrlEvaluacionCumplidoCrud") + query)
	var respuestaPeticion map[string]interface{}
	listaCambiosEstado = make([]models.CambioEstadoASignacionEnvaluacion, 0)

	if response, err := helpers.GetJsonTest(beego.AppConfig.String("UrlEvaluacionCumplidoCrud")+query, &respuestaPeticion); err == nil && response == 200 {

		helpers.LimpiezaRespuestaRefactor(respuestaPeticion, &listaCambiosEstado)
	} else {
		return listaCambiosEstado, fmt.Errorf("Error al consultar cambios de estado")

	}
	return listaCambiosEstado, nil
}

func ObtenerDependencia(codigoDependencia string) (nombreDependencia string, outputError error) {
	defer func() {
		if err := recover(); err != nil {
			outputError = fmt.Errorf("%v", err)
			panic(outputError)
		}
	}()
	query := fmt.Sprintf("/dependencia_SIC/?query=ESFCODIGODEP:%s&limit=1", codigoDependencia)
	fmt.Println(beego.AppConfig.String("UrlAdministrativaAmazonApi") + query)
	var dependencia []models.DependenciaSic
	if response, err := helpers.GetJsonWSO2Test(beego.AppConfig.String("UrlAdministrativaAmazonApi")+query, &dependencia); err == nil && response == 200 {
		nombreDependencia = dependencia[0].ESFDEPENCARGADA
	} else {
		return nombreDependencia, fmt.Errorf("Error al consultar dependencia")
	}

	return nombreDependencia, nil

}

func ObtenerEstadoEvaluacion(idEvaluacion int) (estadoEvaluacion models.EstadoEvaluacion, outputError error) {

	defer func() {
		if err := recover(); err != nil {
			outputError = fmt.Errorf("%v", err)
			panic(outputError)
		}
	}()

	var respuestaPeticion map[string]interface{}
	var estadosEvaluacion []models.CambioEstadoEvaluacion
	query := fmt.Sprintf("/cambio_estado_evaluacion/?query=EvaluacionId.Id:%d,Activo:true", idEvaluacion)
	fmt.Println(beego.AppConfig.String("urlEvaluacionCumplidoCrud") + query)
	if response, err := helpers.GetJsonTest(beego.AppConfig.String("urlEvaluacionCumplidoCrud")+query, &respuestaPeticion); err == nil && response == 200 {
		helpers.LimpiezaRespuestaRefactor(respuestaPeticion, &estadosEvaluacion)
	} else {
		return estadoEvaluacion, fmt.Errorf("Error al consultar cambios de estado")

	}

	if len(estadosEvaluacion) > 0 && estadosEvaluacion[0].EstadoEvaluacionId != nil {
		estadoEvaluacion = *estadosEvaluacion[0].EstadoEvaluacionId
		return estadoEvaluacion, nil
	}
	return estadoEvaluacion, nil

}
