package helpers

import (
	"fmt"

	"github.com/astaxie/beego"
	"github.com/udistrital/evaluacion_cumplido_prov_mid/models"
)

func ObtenerContratoGeneral(numero_contrato_suscrito string, vigencia_contrato string) (contrato_general models.ContratoGeneral, outputError error) {
	defer func() {
		if err := recover(); err != nil {
			outputError = fmt.Errorf("%v", err)
			panic(outputError)
		}
	}()

	var contrato []models.ContratoGeneral
	//fmt.Println("Url contrato general: ", beego.AppConfig.String("UrlAdministrativaAmazonApi")+"/contrato_general/?query=ContratoSuscrito.NumeroContratoSuscrito:"+numero_contrato_suscrito+",VigenciaContrato:"+vigencia_contrato)
	if response, err := GetJsonTest(beego.AppConfig.String("UrlAdministrativaAmazonApi")+"/contrato_general/?query=ContratoSuscrito.NumeroContratoSuscrito:"+numero_contrato_suscrito+",VigenciaContrato:"+vigencia_contrato, &contrato); (err == nil) && (response == 200) {
		if len(contrato) > 0 {
			return contrato[0], nil
		} else {
			outputError = fmt.Errorf("No se encontró contrato")
			return contrato[0], outputError
		}
	} else {
		outputError = fmt.Errorf("Error al obtener el contrato general")
		return contrato_general, outputError
	}
}

func ObtenerDependenciasSupervisor(documento_supervisor string) (dependencias_supervisor []models.Dependencia, outputError error) {
	defer func() {
		if err := recover(); err != nil {
			outputError = fmt.Errorf("%v", err)
		}
	}()

	var respuesta_peticion map[string]interface{}
	fmt.Println("Url dependencias: ", beego.AppConfig.String("UrlAdministrativaJBPM")+"/dependencias_supervisor/"+documento_supervisor)
	if response, err := GetJsonWSO2Test(beego.AppConfig.String("UrlAdministrativaJBPM")+"/dependencias_supervisor/"+documento_supervisor, &respuesta_peticion); err == nil && response == 200 {
		if respuesta_peticion != nil {
			if dependenciasMap, ok := respuesta_peticion["dependencias"].(map[string]interface{}); ok {

				for _, depList := range dependenciasMap {

					if list, ok := depList.([]interface{}); ok {

						for _, dep := range list {

							depMap := dep.(map[string]interface{})
							dependencia := models.Dependencia{

								Codigo: depMap["codigo"].(string),
								Nombre: depMap["nombre"].(string),
							}
							dependencias_supervisor = append(dependencias_supervisor, dependencia)
						}

					} else {
						outputError = fmt.Errorf("No se encontraron dependencias para el supervisor con documento: " + documento_supervisor)
						return dependencias_supervisor, outputError
					}
				}
			}
		} else {
			outputError = fmt.Errorf("No se encontraron dependencias para el supervisor con documento: " + documento_supervisor)
			return dependencias_supervisor, outputError
		}
	}
	return dependencias_supervisor, nil
}

func ObtenerNombrePersonaNatural(documento_persona string) (nombre_persona string, outputError error) {
	defer func() {
		if err := recover(); err != nil {
			outputError = fmt.Errorf("%v", err)
			panic(outputError.Error())
		}
	}()

	var informacion []models.InformacionPersonaNatural
	//fmt.Println("Url informacion persona: ", beego.AppConfig.String("UrlAdministrativaAmazonApi")+"/informacion_persona_natural/?fields=PrimerNombre,SegundoNombre,PrimerApellido,SegundoApellido&limit=0&query=Id:"+documento_persona)
	if response, err := GetJsonTest(beego.AppConfig.String("UrlAdministrativaAmazonApi")+"/informacion_persona_natural/?fields=PrimerNombre,SegundoNombre,PrimerApellido,SegundoApellido&limit=0&query=Id:"+documento_persona, &informacion); err != nil && response != 200 {
		outputError = fmt.Errorf("Error al obtener la información de la persona")
		return nombre_persona, outputError
	}

	if len(informacion) == 0 {
		outputError = fmt.Errorf(fmt.Sprintf("No se encontró información de la persona con documento: %s", documento_persona))
		return nombre_persona, outputError
	}

	nombre_persona = informacion[0].PrimerNombre + " " + informacion[0].SegundoNombre + " " + informacion[0].PrimerApellido + " " + informacion[0].SegundoApellido

	return nombre_persona, nil
}
