package services

import (
	"fmt"

	"github.com/astaxie/beego"
	"github.com/udistrital/evaluacion_cumplido_prov_mid/models"
	"github.com/udistrital/evaluacion_cumplido_prov_mid/utils"
)

func ObternerUnidadMedida(unidad string) (idUnidad *int, outputError error) {

	defer func() {
		if err := recover(); err != nil {
			outputError = fmt.Errorf("%v", err)
			panic(outputError.Error())
		}
	}()

	var unidadMedida []models.UnidadMedida
	//fmt.Println("url", beego.AppConfig.String("UrlAdministrativaAmazonApi")+"/unidad")
	if response, err := utils.GetJsonWSO2Test(beego.AppConfig.String("urlEvaluacionCumplidosCrud")+"/unidad", &unidadMedida); err == nil && response == 200 {

		for _, unidadMedida := range unidadMedida {
			if unidadMedida.Unidad == unidad {
				return &unidadMedida.Id, outputError
			}
		}
	}
	return nil, outputError
}

func GuardarItems(items []models.ItemEvaluacion) (response map[string]interface{}, err error) {

	var respuesta_peticion map[string]interface{}

	if err := utils.SendJson(beego.AppConfig.String("urlEvaluacionCumplidosCrud")+"/item/guardado_multiple", "POST", &respuesta_peticion, items); err == nil {

		return respuesta_peticion, nil
	}

	return nil, err
}
