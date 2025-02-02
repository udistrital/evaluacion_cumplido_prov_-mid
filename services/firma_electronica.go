package services

import (
	"encoding/json"
	"fmt"

	"github.com/astaxie/beego"
	"github.com/udistrital/evaluacion_cumplido_prov_mid/helpers"
	"github.com/udistrital/evaluacion_cumplido_prov_mid/models"
)

func EjecutarProcesoDefirma(peticion_firma models.PeticionFirmaElectronica) (map_response map[string]interface{}, outputError error) {

	defer func() {
		if err := recover(); err != nil {
			outputError = fmt.Errorf("%v", err)
			panic(outputError.Error())
		}
	}()
	var query string
	var estapa_evalucion = 2
	map_response = make(map[string]interface{})
	estado_actual, err := ConsultarEstadoActualAsingacion(peticion_firma.AsignacionId)

	if err != nil {
		outputError = fmt.Errorf("error al consultar el estado de la evaluacion")
	}

	if estado_actual.EstadoAsignacionEvaluador.CodigoAbreviacion == "ER" {
		//Consultar información de la persona natural
		info_persona_natural, err := ConsultarInfoPersonaNatural(peticion_firma.PersonaId)

		if err != nil {
			outputError = fmt.Errorf("error al consultar la información de la persona natural")
			return nil, outputError
		}

		//consultar la asignacion
		asignacion, err := ConsultarAsignacion(peticion_firma.AsignacionId)

		if err != nil {
			outputError = fmt.Errorf("error al consultar la información de la asignación")
			return nil, outputError
		}

		//consultar la evaluacion
		evaluacion, err := ConsultarEvaluacion(asignacion.EvaluacionId.Id)

		if err != nil {
			outputError = fmt.Errorf("error al consultar la información de la evaluación")
			return nil, outputError
		}

		//consultar el documento en documento crud
		documento_crud, err := consultarDocumentocrud(evaluacion.DocumentoId)

		if err != nil {
			outputError = fmt.Errorf("error al consultar el documento en el documento crud")
			return nil, outputError
		}

		primer_firmante, length, err := VerificarPrimerFirmante(asignacion.EvaluacionId.Id, "EAP")

		if primer_firmante && length == 1 {
			estapa_evalucion = 1
			query = "/firma_electronica_mid/v1/firma_electronica"

		}
		if primer_firmante && length > 1 {
			estapa_evalucion = 1
			query = "/firma_electronica_mid/v1/firma_multiple"
		}

		ultimo_firmante, err := VerificarUltimoFirmanteFirmante(asignacion.EvaluacionId.Id, "EAP")

		if ultimo_firmante && length > 1 {
			query = "/firma_electronica_mid/v1/firma_multiple"
			estapa_evalucion = 3

		}

		//consultar el documento en el gestor documental
		documento_gestor_documental, err := consultarGestorDocumental(documento_crud.Enlace)

		//consultar firmantes
		metadata := documento_crud.Metadatos

		lista_firmante, err := obtenerFirmantesDesdeJSON(metadata)

		//crear la peticion de firma

		peticion_firma_electronica := crearPeticionFirmaElectronica(*info_persona_natural, *documento_gestor_documental, *evaluacion, estapa_evalucion, lista_firmante)

		fmt.Println(peticion_firma_electronica)

		respuesta_firma, err := FirmarDocumento(peticion_firma_electronica, query)

		if err != nil {

			return nil, err
		}

		_, err = CambioEstadoAsignacionEvaluacion(peticion_firma.AsignacionId, "EAP")

		if err != nil {
			outputError = fmt.Errorf("error al firmar el documento")
		}

		guardarDocumentoFirmado(respuesta_firma.Res.Id, asignacion.EvaluacionId.Id)

		if err != nil {

			return nil, err
		}
		map_response["Message"] = "Se cambio el estado correctamente"
		map_response["Data"] = respuesta_firma

	} else {

		map_response["Message"] = "El estado de la asignación no permite la firma"

	}

	return map_response, nil
}

func crearPeticionFirmaElectronica(persona models.Persona, documento_gestor_documental models.DocumentoEnlace, evaluacion models.Evaluacion, estapa_firma int, lista_firmantes []models.Firmante) (peticion_firma []models.PeticionFirmaElectronicaCrud) {

	var metadatos = models.Metadatos{}

	metadatos["firmantes"] = []interface{}{
		map[string]interface{}{
			"Nombre":         persona.PrimerApellido + " " + persona.SegundoNombre + " " + persona.PrimerApellido + " " + persona.SegundoApellido,
			"Cargo":          "Evaluador",
			"TipoId":         persona.TipoDocumento.Abreviatura,
			"Identificacion": persona.Id,
		},
	}

	if len(lista_firmantes) > 0 {
		for _, firmante := range lista_firmantes {
			metadatos["firmantes"] = append(metadatos["firmantes"].([]interface{}), map[string]interface{}{
				"Nombre":         firmante.Nombre,
				"Cargo":          firmante.Cargo,
				"TipoId":         firmante.TipoId,
				"Identificacion": firmante.Identificacion,
			})
		}
	}
	metadatos["representantes"] = []interface{}{}
	var peticion = models.PeticionFirmaElectronicaCrud{
		IdTipoDocumento: 158,
		Nombre:          fmt.Sprintf("Evaluacion del contrato %d , con vigencia %d", evaluacion.ContratoSuscritoId, evaluacion.VigenciaContrato),
		Metadatos:       metadatos,
		Representantes:  []models.Firmante{},
		Firmantes: []models.Firmante{
			{Nombre: persona.PrimerApellido + " " + persona.SegundoNombre + " " + persona.PrimerApellido + " " + persona.SegundoApellido, Cargo: "Evaluador", TipoId: persona.TipoDocumento.Abreviatura, Identificacion: persona.Id},
		},
		Descripcion: "Firma de la evaluación del contrato",
		EtapaFirma:  estapa_firma,
		File:        documento_gestor_documental.File,
	}

	peticion_firma = append(peticion_firma, peticion)

	return peticion_firma
}
func FirmarDocumento(peticion_firma []models.PeticionFirmaElectronicaCrud, query string) (respuestaPeticion *models.RespuestaFirmaElectronica, outputError error) {

	defer func() {
		if err := recover(); err != nil {
			outputError = fmt.Errorf("%v", err)
			panic(outputError.Error())
		}
	}()

	//var respuesta_firma models.RespuestaFirMaElectronica
	fmt.Println(beego.AppConfig.String("UrlFirmaElectronica") + query)

	if err := helpers.SendJson(beego.AppConfig.String("UrlFirmaElectronica")+query, "POST", &respuestaPeticion, peticion_firma); err == nil {

	} else {
		fmt.Println(err)
		return nil, err
	}
	if respuestaPeticion.Res.Id == 0 {
		outputError = fmt.Errorf("error al firmar el documento")
		return nil, outputError
	}
	return respuestaPeticion, nil

}

func guardarDocumentoFirmado(id_documento_firmado int, evaluacuion_id int) (outputError error) {

	defer func() {
		if err := recover(); err != nil {
			outputError = fmt.Errorf("%v", err)
			panic(outputError.Error())
		}
	}()

	evaluacion, err := ConsultarEvaluacion(evaluacuion_id)

	if err != nil {
		outputError = fmt.Errorf("error consultar la evaluacion")
		return outputError
	}

	if evaluacion.Id == 0 {
		outputError = fmt.Errorf("no hay resultado para  la evaluacion con id %d", evaluacuion_id)
		return outputError
	}
	evaluacion.DocumentoId = id_documento_firmado

	var map_response = make(map[string]interface{})
	query := fmt.Sprintf("/evaluacion/%d", evaluacion.Id)

	if err := helpers.SendJson(beego.AppConfig.String("UrlEvaluacionCumplidoCrud")+query, "PUT", &map_response, evaluacion); err != nil {
		return nil
	}
	outputError = fmt.Errorf("error al guardar el documento firmado")
	return outputError
}

func ConsultarInfoPersonaNatural(numero_documento string) (info_persona_natural *models.Persona, outputError error) {

	defer func() {
		if err := recover(); err != nil {
			outputError = fmt.Errorf("%v", err)
			panic(outputError.Error())
		}
	}()

	var consulta_info_persona_natural []models.Persona
	query := fmt.Sprintf("/informacion_persona_natural?query=Id:%s", numero_documento)
	fmt.Print(beego.AppConfig.String("UrlAdministrativaAmazonApi") + query)
	if response, err := helpers.GetJsonWSO2Test(beego.AppConfig.String("UrlAdministrativaAmazonApi")+query, &consulta_info_persona_natural); err == nil && response == 200 {

		if consulta_info_persona_natural != nil && consulta_info_persona_natural[0].Id != "" {

			info_persona_natural = &consulta_info_persona_natural[0]

		} else {
			outputError = fmt.Errorf("error al consultar la información de la persona natural")
			return nil, outputError
		}
	} else {
		outputError = fmt.Errorf("error al consultar la información de la persona natural")
		return info_persona_natural, outputError
	}
	fmt.Println(info_persona_natural)
	return info_persona_natural, nil
}

// /Se consulta la asigancion de un evaluador
func ConsultarAsignacion(id_asiganacion int) (asignacion *models.AsignacionEvaluador, outputError error) {

	defer func() {
		if err := recover(); err != nil {
			outputError = fmt.Errorf("%v", err)
			panic(outputError.Error())
		}
	}()
	var map_response = make(map[string]interface{})
	var consulta_asignaciones []models.AsignacionEvaluador

	query := fmt.Sprintf("/asignacion_evaluador?query=Id:%d", id_asiganacion)
	fmt.Println(beego.AppConfig.String("UrlEvaluacionCumplidoCrud") + query)
	if response, err := helpers.GetJsonWSO2Test(beego.AppConfig.String("UrlEvaluacionCumplidoCrud")+query, &map_response); err == nil && response == 200 {

		helpers.LimpiezaRespuestaRefactor(map_response, &consulta_asignaciones)

		if consulta_asignaciones != nil && consulta_asignaciones[0].Id != 0 {

			asignacion = &consulta_asignaciones[0]

		} else {
			outputError = fmt.Errorf("no hay información de la asignación")
			return nil, outputError
		}

	} else {
		outputError = fmt.Errorf("error al consultar la información de la asignación")
		return asignacion, outputError
	}
	fmt.Println(asignacion)
	return asignacion, nil
}

// /Se consulta la evalcuion
func ConsultarEvaluacion(id_evalacion int) (evaluacion *models.Evaluacion, outputError error) {

	defer func() {
		if err := recover(); err != nil {
			outputError = fmt.Errorf("%v", err)
			panic(outputError.Error())
		}
	}()

	var consulta_evaluaciones []models.Evaluacion
	var map_response = make(map[string]interface{})
	query := fmt.Sprintf("/evaluacion/?query=Id:%d", id_evalacion)
	if response, err := helpers.GetJsonWSO2Test(beego.AppConfig.String("UrlEvaluacionCumplidoCrud")+query, &map_response); err == nil && response == 200 {
		helpers.LimpiezaRespuestaRefactor(map_response, &consulta_evaluaciones)
		if consulta_evaluaciones != nil && consulta_evaluaciones[0].Id != 0 {

			evaluacion = &consulta_evaluaciones[0]

		} else {
			outputError = fmt.Errorf("error al consultar la información de la persona natural")
			return nil, outputError
		}
	} else {
		outputError = fmt.Errorf("error al consultar la información de la persona natural")
		return nil, outputError
	}
	fmt.Println(evaluacion)
	return evaluacion, nil
}

// /Se consulta el codumnetro en documento crud y se retorna el documento de gestor documental
func consultarDocumentocrud(id_documento int) (documento *models.DocumentoCrud, outputError error) {

	defer func() {
		if err := recover(); err != nil {
			outputError = fmt.Errorf("%v", err)
			panic(outputError.Error())
		}
	}()

	var peticion_consulta_documento []models.DocumentoCrud

	query := fmt.Sprintf("/documento/?limit=-1&query=Id.in:%d", id_documento)
	fmt.Println(beego.AppConfig.String("UrlDocumentosCrud") + query)
	if response, err := helpers.GetJsonTest(beego.AppConfig.String("UrlDocumentosCrud")+query, &peticion_consulta_documento); response == 200 {

		if peticion_consulta_documento != nil && peticion_consulta_documento[0].Id != 0 {

			documento = &peticion_consulta_documento[0]
			return documento, nil

		}

	} else {
		return nil, err
	}

	return nil, outputError
}

// /Se consulta el documento en el gestor documental
func consultarGestorDocumental(enlace string) (documento *models.DocumentoEnlace, outputError error) {

	defer func() {
		if err := recover(); err != nil {
			outputError = fmt.Errorf("%v", err)
			panic(outputError.Error())
		}
	}()

	var peticion_consulta_documento models.DocumentoEnlace
	query := fmt.Sprintf("/document/%s", enlace)
	fmt.Println(beego.AppConfig.String("UrlGestorDocumental") + query)
	if response, err := helpers.GetJsonWSO2Test(beego.AppConfig.String("UrlGestorDocumental")+query, &peticion_consulta_documento); err == nil && response == 200 {

		if peticion_consulta_documento.File != "" {
			documento = &peticion_consulta_documento
			return documento, nil
		}

	} else {
		outputError = fmt.Errorf("error al consultar el documento en el gestor documental")
		return nil, outputError
	}
	return nil, nil
}

func VerificarPrimerFirmante(id_evaluacion int, estado_abreviacion string) (primer_firmante bool, length int, outputError error) {

	primer_firmante = true

	asiganciones, err := ConsultarAsignacionesPorIdEvaluacion(id_evaluacion)

	if err != nil {
		outputError = fmt.Errorf("error al consultar asignaciones")
		return false, 0, outputError
	}
	length = len(*asiganciones)

	for _, asigancion := range *asiganciones {

		estado_asiganacion, err := ConsultarEstadoActualAsingacion(asigancion.Id)

		if err != nil {
			fmt.Printf("error al consultar el estado de la asignacion %s", estado_asiganacion.Id)

		}

		if estado_asiganacion.EstadoAsignacionEvaluador.CodigoAbreviacion == estado_abreviacion {
			primer_firmante = false
			break
		}

	}

	return primer_firmante, length, outputError

}

func VerificarUltimoFirmanteFirmante(id_evaluacion int, estado_abreviacion string) (ultimo_firmante bool, outputError error) {

	ultimo_firmante = true
	var estados []string

	asiganciones, err := ConsultarAsignacionesPorIdEvaluacion(id_evaluacion)

	if err != nil {
		outputError = fmt.Errorf("error al consultar asignaciones")
		return false, outputError
	}

	for _, asigancion := range *asiganciones {

		estado_asiganacion, err := ConsultarEstadoActualAsingacion(asigancion.Id)

		if err != nil {
			fmt.Printf("error al consultar el estado de la asignacion %s", estado_asiganacion.Id)

		}

		if estado_asiganacion.EstadoAsignacionEvaluador.CodigoAbreviacion != "EAP" {
			estados = append(estados, estado_asiganacion.EstadoAsignacionEvaluador.CodigoAbreviacion)
		}

	}

	if len(estados) > 1 {
		ultimo_firmante = false
	}
	return ultimo_firmante, outputError

}

func obtenerFirmantesDesdeJSON(metadata string) (lista_firmante []models.Firmante, outputError error) {
	var metadatos map[string]interface{}

	if err := json.Unmarshal([]byte(metadata), &metadatos); err != nil {
		return nil, err
	}

	if firmantesData, ok := metadatos["firmantes"].([]interface{}); ok {

		for _, item := range firmantesData {
			if firmante, ok := item.(map[string]interface{}); ok {
				lista_firmante = append(lista_firmante, models.Firmante{
					Cargo:          firmante["Cargo"].(string),
					Identificacion: firmante["Identificacion"].(string),
					Nombre:         firmante["Nombre"].(string),
					TipoId:         firmante["TipoId"].(string),
				})
			}
		}
	}

	return lista_firmante, nil
}
