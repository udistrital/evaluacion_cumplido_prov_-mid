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

func EnviarNotificacionesAsignacionEvaluacion(evaluacionId int) (errores []string, outputError error) {
	defer func() {
		if err := recover(); err != nil {
			outputError = fmt.Errorf("%v", err)
			panic(outputError.Error())
		}
	}()

	// Obtener todos los evaluadores de la evaluacion
	var respuesta_evaluadores map[string]interface{}
	var evaluadores []models.AsignacionEvaluador

	//fmt.Println("URL evaluadores: ", beego.AppConfig.String("UrlEvaluacionCumplidoCrud")+"/asignacion_evaluador/?query=EvaluacionId.Id:"+strconv.Itoa(evaluacion_id)+",Activo:true&limit=-1")
	if response, err := helpers.GetJsonTest(beego.AppConfig.String("UrlEvaluacionCumplidoCrud")+"/asignacion_evaluador/?query=EvaluacionId.Id:"+strconv.Itoa(evaluacionId)+",Activo:true&limit=-1", &respuesta_evaluadores); err != nil && response != 200 {
		outputError = fmt.Errorf("Error al obtener los evaluadores de la evaluación")
		return errores, outputError
	}

	// Verificar si la evaluación tiene evaluadores asignados
	data := respuesta_evaluadores["Data"].([]interface{})
	if len(data[0].(map[string]interface{})) == 0 {
		outputError = fmt.Errorf("La evaluación no tiene evaluadores asignados")
		return errores, outputError
	}

	helpers.LimpiezaRespuestaRefactor(respuesta_evaluadores, &evaluadores)

	// Obtener el nombre del proveedor
	contrato_general, err := helpers.ObtenerContratoGeneral(strconv.Itoa(evaluadores[0].EvaluacionId.ContratoSuscritoId), strconv.Itoa(evaluadores[0].EvaluacionId.VigenciaContrato))
	if err != nil {
		outputError = fmt.Errorf(err.Error())
		return errores, outputError
	}
	var informacion_proveedor []models.InformacionProveedor
	if response, err := helpers.GetJsonTest(beego.AppConfig.String("UrlAdministrativaAmazonApi")+"/informacion_proveedor/?query=Id:"+strconv.Itoa(contrato_general.Contratista), &informacion_proveedor); (err != nil) && (response != 200) {
		outputError = fmt.Errorf("Error al obtener la información del proveedor")
		return errores, outputError
	}

	// Enviar notificaciones a los evaluadores
	for _, evaluador := range evaluadores {
		// Obtener el correo del evaluador
		var email string
		var autenticacion_persona models.AutenticacionPersona
		var respuesta_autenticacion map[string]interface{}

		if beego.AppConfig.String("runmode") != "prod" {
			// Poner aqui el email que se usara para realizar las pruebas de las notificaciones
			email = ""
		} else {
			body_autenticacion := map[string]interface{}{
				"numero": evaluador.PersonaId,
			}
			if err := helpers.SendJson(beego.AppConfig.String("UrlAutenticacionMid")+"/token/documentoToken", "POST", &respuesta_autenticacion, body_autenticacion); err != nil {
				errores = append(errores, fmt.Sprintf("Error al obtener el correo del evaluador %v", evaluador.PersonaId))
				continue
			}

			json_autenticacion, err := json.Marshal(respuesta_autenticacion)
			if err != nil {
				errores = append(errores, fmt.Sprintf("Error al obtener el correo del evaluador %v", evaluador.PersonaId))
				continue
			}
			err = json.Unmarshal(json_autenticacion, &autenticacion_persona)
			if err != nil {
				errores = append(errores, fmt.Sprintf("Error al obtener el correo del evaluador %v", evaluador.PersonaId))
				continue
			}

			if autenticacion_persona.Email == "" {
				errores = append(errores, fmt.Sprintf("El evaluador %v no tiene correo registrado", evaluador.PersonaId))
				continue
			}
			email = autenticacion_persona.Email
		}

		//Obtener los items evaluados del evaluador
		_, items_evaluador_formateados, err := helpers.ObtenerItemsEvaluador(evaluador.Id)
		if err != nil {
			errores = append(errores, fmt.Sprintf("Error al obtener los items evaluados del evaluador %v", evaluador.PersonaId))
			continue
		}

		// Enviar notificación al correo del evaluador
		var body_enviar_notificacion models.NotificacionAsignacionEvaluacionEmail
		body_enviar_notificacion.Source = "notificacionescumplidosproveedores@udistrital.edu.co"
		body_enviar_notificacion.Template = "PLANTILLA_ASIGNACION_EVALUACION"
		body_enviar_notificacion.Destinations = make([]models.DestinationAsignacionEvaluacion, 1)
		fmt.Println("Email: ", email)
		body_enviar_notificacion.Destinations[0].Destination.ToAddresses = append(body_enviar_notificacion.Destinations[0].Destination.ToAddresses, email)
		body_enviar_notificacion.Destinations[0].ReplacementTemplateData.RolEvaluador = evaluador.RolAsignacionEvaluadorId.Nombre
		body_enviar_notificacion.Destinations[0].ReplacementTemplateData.NombreProveedor = informacion_proveedor[0].NomProveedor
		body_enviar_notificacion.Destinations[0].ReplacementTemplateData.ItemsEvaluar = items_evaluador_formateados
		body_enviar_notificacion.Destinations[0].Attachments = []string{}
		body_enviar_notificacion.Destinations[0].Destination.BccAddresses = []string{}
		body_enviar_notificacion.Destinations[0].Destination.CcAddresses = []string{}
		body_enviar_notificacion.DefaultTemplateData.RolEvaluador = evaluador.RolAsignacionEvaluadorId.Nombre
		body_enviar_notificacion.DefaultTemplateData.NombreProveedor = informacion_proveedor[0].NomProveedor
		body_enviar_notificacion.DefaultTemplateData.ItemsEvaluar = items_evaluador_formateados

		var respuesta map[string]interface{}
		if err := helpers.SendJsonTls(beego.AppConfig.String("UrlNotificacionesMid")+"/email/enviar_templated_email", "POST", &respuesta, body_enviar_notificacion); err != nil {
			errores = append(errores, fmt.Sprintf("Error al enviar la notificación al correo del evaluador %v", evaluador.PersonaId))
			continue
		}
		fmt.Println("Respuesta: ", respuesta)
	}
	return errores, nil
}

func EnviarNotificacionRealizacionEvaluacion(documento_evaluador string, vigencia string, numero_contrato string) (outputError error) {
	defer func() {
		if err := recover(); err != nil {
			outputError = fmt.Errorf("%v", err)
			panic(outputError.Error())
		}
	}()

	// Obtener el nombre del proveedor
	contrato_general, err := helpers.ObtenerContratoGeneral(numero_contrato, vigencia)
	if err != nil {
		outputError = fmt.Errorf(err.Error())
		return outputError
	}
	var informacion_proveedor []models.InformacionProveedor
	if response, err := helpers.GetJsonTest(beego.AppConfig.String("UrlAdministrativaAmazonApi")+"/informacion_proveedor/?query=Id:"+strconv.Itoa(contrato_general.Contratista), &informacion_proveedor); (err != nil) && (response != 200) {
		outputError = fmt.Errorf("Error al obtener la información del proveedor")
		return outputError
	}

	// Obtener el correo del evaluador
	var email string

	if beego.AppConfig.String("runmode") != "prod" {
		// Poner aqui el email que se usara para realizar las pruebas de las notificaciones
		email = "fctrujilloo@udistrital.edu.co"
	} else {
		var autenticacion_persona models.AutenticacionPersona
		var respuesta_autenticacion map[string]interface{}
		body_autenticacion := map[string]interface{}{
			"numero": documento_evaluador,
		}
		if err := helpers.SendJson(beego.AppConfig.String("UrlAutenticacionMid")+"/token/documentoToken", "POST", &respuesta_autenticacion, body_autenticacion); err != nil {
			outputError = fmt.Errorf("Error al obtener el correo del evaluador %v", documento_evaluador)
			return outputError
		}

		json_autenticacion, err := json.Marshal(respuesta_autenticacion)
		if err != nil {
			outputError = fmt.Errorf("Error al obtener el correo del evaluador %v", documento_evaluador)
			return outputError
		}
		err = json.Unmarshal(json_autenticacion, &autenticacion_persona)
		if err != nil {
			outputError = fmt.Errorf("Error al obtener el correo del evaluador %v", documento_evaluador)
			return outputError
		}

		if autenticacion_persona.Email == "" {
			outputError = fmt.Errorf("El evaluador %v no tiene correo registrado", documento_evaluador)
			return outputError
		}
		email = autenticacion_persona.Email
	}

	// Enviar notificación al correo del evaluador
	var body_enviar_notificacion models.NotificacionRealizacionEvaluacionEmail
	body_enviar_notificacion.Source = "notificacionescumplidosproveedores@udistrital.edu.co"
	body_enviar_notificacion.Template = "PLANTILLA_REALIZACION_EVALUACION"
	body_enviar_notificacion.Destinations = make([]models.DestinationRealizacionEvaluacion, 1)
	fmt.Println("Email: ", email)
	body_enviar_notificacion.Destinations[0].Destination.ToAddresses = append(body_enviar_notificacion.Destinations[0].Destination.ToAddresses, email)
	body_enviar_notificacion.Destinations[0].ReplacementTemplateData.NombreProveedor = informacion_proveedor[0].NomProveedor
	body_enviar_notificacion.Destinations[0].ReplacementTemplateData.Vigencia = vigencia
	body_enviar_notificacion.Destinations[0].ReplacementTemplateData.NumeroContrato = numero_contrato
	body_enviar_notificacion.Destinations[0].Attachments = []string{}
	body_enviar_notificacion.Destinations[0].Destination.BccAddresses = []string{}
	body_enviar_notificacion.Destinations[0].Destination.CcAddresses = []string{}
	body_enviar_notificacion.DefaultTemplateData.NombreProveedor = informacion_proveedor[0].NomProveedor
	body_enviar_notificacion.DefaultTemplateData.Vigencia = vigencia
	body_enviar_notificacion.DefaultTemplateData.NumeroContrato = numero_contrato

	var respuesta map[string]interface{}
	if err := helpers.SendJsonTls(beego.AppConfig.String("UrlNotificacionesMid")+"/email/enviar_templated_email", "POST", &respuesta, body_enviar_notificacion); err != nil {
		outputError = fmt.Errorf("Error al enviar la notificación al correo del evaluador %v", documento_evaluador)
		return outputError
	}
	fmt.Println("Respuesta: ", respuesta)
	return nil
}

func EnviarNotificacionesFinalizacionEvaluacion(evaluacionId int, numero_contrato string, vigencia string) (errores []string, outputError error) {
	defer func() {
		if err := recover(); err != nil {
			outputError = fmt.Errorf("%v", err)
			panic(outputError.Error())
		}
	}()

	var nombre_evaluadores []string

	//Obtener el nombre de los evaluadores
	var respuesta_evaluadores map[string]interface{}
	var evaluadores []models.AsignacionEvaluador

	//fmt.Println("URL evaluadores: ", beego.AppConfig.String("UrlEvaluacionCumplidoCrud")+"/asignacion_evaluador/?query=EvaluacionId.Id:"+strconv.Itoa(evaluacion_id)+",Activo:true&limit=-1")
	if response, err := helpers.GetJsonTest(beego.AppConfig.String("UrlEvaluacionCumplidoCrud")+"/asignacion_evaluador/?query=EvaluacionId.Id:"+strconv.Itoa(evaluacionId)+",Activo:true&limit=-1", &respuesta_evaluadores); err != nil && response != 200 {
		outputError = fmt.Errorf("Error al obtener los evaluadores de la evaluación")
		return errores, outputError
	}

	// Verificar si la evaluación tiene evaluadores asignados
	data := respuesta_evaluadores["Data"].([]interface{})
	if len(data[0].(map[string]interface{})) == 0 {
		outputError = fmt.Errorf("La evaluación no tiene evaluadores asignados")
		return errores, outputError
	}

	helpers.LimpiezaRespuestaRefactor(respuesta_evaluadores, &evaluadores)

	for _, evaluador := range evaluadores {
		nombre_evaluador, err := helpers.ObtenerNombrePersonaNatural(evaluador.PersonaId)
		if err != nil {
			outputError = fmt.Errorf(err.Error())
			return errores, outputError
		}
		nombre_evaluadores = append(nombre_evaluadores, nombre_evaluador)
	}

	// Obtener el nombre del proveedor
	contrato_general, err := helpers.ObtenerContratoGeneral(numero_contrato, vigencia)
	if err != nil {
		outputError = fmt.Errorf(err.Error())
		return errores, outputError
	}
	var informacion_proveedor []models.InformacionProveedor
	if response, err := helpers.GetJsonTest(beego.AppConfig.String("UrlAdministrativaAmazonApi")+"/informacion_proveedor/?query=Id:"+strconv.Itoa(contrato_general.Contratista), &informacion_proveedor); (err != nil) && (response != 200) {
		outputError = fmt.Errorf("Error al obtener la información del proveedor")
		return errores, outputError
	}

	// Enviar notificaciones a los evaluadores
	for _, evaluador := range evaluadores {
		var email string
		// Obtener el correo del evaluador
		var autenticacion_persona models.AutenticacionPersona
		var respuesta_autenticacion map[string]interface{}
		if beego.AppConfig.String("runmode") != "prod" {
			// Poner aqui el email que se usara para realizar las pruebas de las notificaciones
			email = "fctrujilloo@udistrital.edu.co"
		} else {
			body_autenticacion := map[string]interface{}{
				"numero": evaluador.PersonaId,
			}
			if err := helpers.SendJson(beego.AppConfig.String("UrlAutenticacionMid")+"/token/documentoToken", "POST", &respuesta_autenticacion, body_autenticacion); err != nil {
				errores = append(errores, fmt.Sprintf("Error al obtener el correo del evaluador %v", evaluador.PersonaId))
				continue
			}

			json_autenticacion, err := json.Marshal(respuesta_autenticacion)
			if err != nil {
				errores = append(errores, fmt.Sprintf("Error al obtener el correo del evaluador %v", evaluador.PersonaId))
				continue
			}
			err = json.Unmarshal(json_autenticacion, &autenticacion_persona)
			if err != nil {
				errores = append(errores, fmt.Sprintf("Error al obtener el correo del evaluador %v", evaluador.PersonaId))
				continue
			}

			if autenticacion_persona.Email == "" {
				errores = append(errores, fmt.Sprintf("El evaluador %v no tiene correo registrado", evaluador.PersonaId))
				continue
			}
			email = autenticacion_persona.Email
		}

		// Enviar notificación al correo del evaluador
		var body_enviar_notificacion models.NotificacionFinalizacionEvaluacionEmail
		body_enviar_notificacion.Source = "notificacionescumplidosproveedores@udistrital.edu.co"
		body_enviar_notificacion.Template = "PLANTILLA_EVALUACION_FINALIZADA"
		body_enviar_notificacion.Destinations = make([]models.DestinationFinalizacionEvaluacion, 1)
		fmt.Println("Email: ", email)
		body_enviar_notificacion.Destinations[0].Destination.ToAddresses = append(body_enviar_notificacion.Destinations[0].Destination.ToAddresses, email)
		body_enviar_notificacion.Destinations[0].ReplacementTemplateData.NombreProveedor = informacion_proveedor[0].NomProveedor
		body_enviar_notificacion.Destinations[0].ReplacementTemplateData.Vigencia = vigencia
		body_enviar_notificacion.Destinations[0].ReplacementTemplateData.NumeroContrato = numero_contrato
		body_enviar_notificacion.Destinations[0].ReplacementTemplateData.NombreEvaluadores = strings.Join(nombre_evaluadores, ", ")
		body_enviar_notificacion.Destinations[0].Attachments = []string{}
		body_enviar_notificacion.Destinations[0].Destination.BccAddresses = []string{}
		body_enviar_notificacion.Destinations[0].Destination.CcAddresses = []string{}
		body_enviar_notificacion.DefaultTemplateData.NombreProveedor = informacion_proveedor[0].NomProveedor
		body_enviar_notificacion.DefaultTemplateData.Vigencia = vigencia
		body_enviar_notificacion.DefaultTemplateData.NumeroContrato = numero_contrato
		body_enviar_notificacion.DefaultTemplateData.NombreEvaluadores = strings.Join(nombre_evaluadores, ", ")

		var respuesta map[string]interface{}
		if err := helpers.SendJsonTls(beego.AppConfig.String("UrlNotificacionesMid")+"/email/enviar_templated_email", "POST", &respuesta, body_enviar_notificacion); err != nil {
			errores = append(errores, fmt.Sprintf("Error al enviar la notificación al correo del evaluador %v", evaluador.PersonaId))
			continue
		}

		fmt.Println("Respuesta: ", respuesta)

	}
	return errores, nil

}
