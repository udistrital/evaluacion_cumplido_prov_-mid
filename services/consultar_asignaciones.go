package services

import (
	"fmt"

	"github.com/astaxie/beego"
	"github.com/udistrital/evaluacion_cumplido_prov_mid/helpers"
	"github.com/udistrital/evaluacion_cumplido_prov_mid/models"
)

func ConsultarAsignaciones(documento string) (outputError error) {

	defer func() {
		if err := recover(); err != nil {
			outputError = fmt.Errorf("%v", err)
			panic(outputError)
		}
	}()

	var AsignacionesEvaluaciones []models.AsignacionEvaluacion

	fmt.Println("URL: ", beego.AppConfig.String("UrlEvaluacionesCumplidosProveedoresCrud")+"/supervisor/contratos-supervisor/"+documento)
	if response, err := helpers.GetJsonWSO2Test(beego.AppConfig.String("UrlEvaluacionesCumplidosProveedoresCrud")+"/supervisor/contratos-supervisor/"+documento, &AsignacionesEvaluaciones); err == nil && response == 200 {
		println("AsignacionesEvaluaciones: ", response)
		println("AsignacionesEvaluaciones: ", response)

	}
	fmt.Println("AsignacionesEvaluaciones: ", AsignacionesEvaluaciones)
	return nil
}
