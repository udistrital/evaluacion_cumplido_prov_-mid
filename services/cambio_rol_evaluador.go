package services

import (
	"fmt"

	"github.com/astaxie/beego"
	"github.com/udistrital/evaluacion_cumplido_prov_mid/helpers"
	"github.com/udistrital/evaluacion_cumplido_prov_mid/models"
)

func CambiarRolAsignacionEvaluador(idEvaluacion string) (resultado_map map[string]interface{}, outputError error) {

	defer func() {
		if err := recover(); err != nil {
			outputError = fmt.Errorf("%v", err)
			panic(outputError)
		}
	}()
	resultado_map = make(map[string]interface{})
	//Consulta el estado actual de la evaluacion
	estadoEvaluacion, err := consultarCambioEstadoEvaluacion(idEvaluacion)

	if err != nil {
		return nil, fmt.Errorf("error al consultar asignaciones")
	}
	//consutla las asignaciones con rol de evaluador
	listaAsiganaciones, err := consultarAsignacionesPorIdEvaluacionYEstadoEvaluador(idEvaluacion)

	if err != nil {
		return nil, fmt.Errorf("error al consultar asignaciones")
	}

	if estadoEvaluacion.EstadoEvaluacionId.CodigoAbreviacion == "EPR" {
		lista, err := agregarRolEvaluador(listaAsiganaciones)
		if err != nil {
			return nil, fmt.Errorf("error al consultar asignaciones")
		}
		resultado_map["rolesNoAgregado"] = lista
	}

	if estadoEvaluacion.EstadoEvaluacionId.CodigoAbreviacion == "AEV" {
		lista, err := eliminarRolEvaluador(listaAsiganaciones)
		if err != nil {
			return nil, fmt.Errorf("error al consultar asignaciones")
		}
		resultado_map["rolesNoEliminados"] = lista

	}
	return nil, nil
}

func consultarAsignacionesPorIdEvaluacionYEstadoEvaluador(idEvaluacion string) (asignaciones []models.AsignacionEvaluador, outputError error) {
	defer func() {
		if err := recover(); err != nil {
			outputError = fmt.Errorf("%v", err)
			panic(outputError)
		}
	}()
	var respuestaPeticion map[string]interface{}
	var listaAsignacionEvaluador []models.AsignacionEvaluador

	fmt.Println(beego.AppConfig.String("UrlEvaluacionCumplidoCrud") + "/asignacion_evaluador?query=EvaluacionId.Id:" + idEvaluacion + ",RolAsignacionEvaluadorId.CodigoAbreviacion:EV")
	if response, err := helpers.GetJsonWSO2Test(beego.AppConfig.String("UrlEvaluacionCumplidoCrud")+"/asignacion_evaluador?query=EvaluacionId.Id:"+idEvaluacion+",RolAsignacionEvaluadorId.CodigoAbreviacion:EV", &respuestaPeticion); err == nil && response == 200 {
		helpers.LimpiezaRespuestaRefactor(respuestaPeticion, &listaAsignacionEvaluador)
		if len(listaAsignacionEvaluador) > 0 && listaAsignacionEvaluador[0].EvaluacionId != nil {
			asignaciones = listaAsignacionEvaluador

		}
	} else {
		return asignaciones, fmt.Errorf("error al consultar asignaciones")

	}
	return asignaciones, nil
}

func consultarAsignacionesPorDocumentoEvaluadorYEstadoEvaluador(personaId string) (asignacionesPendienesPorEvaluar bool, outputError error) {
	defer func() {
		if err := recover(); err != nil {
			outputError = fmt.Errorf("%v", err)
			panic(outputError)
		}
	}()
	asignacionesPendienesPorEvaluar = false
	var respuestaPeticion map[string]interface{}
	var listaAsignacionEvaluador []models.AsignacionEvaluador
	fmt.Println(beego.AppConfig.String("UrlEvaluacionCumplidoCrud") + "/asignacion_evaluador?query=EvaluacionId.Id:" + personaId + ",RolAsignacionEvaluadorId.CodigoAbreviacion:EV")
	if response, err := helpers.GetJsonWSO2Test(beego.AppConfig.String("UrlEvaluacionCumplidoCrud")+"/asignacion_evaluador?query=PersonaId:"+personaId+",RolAsignacionEvaluadorId.CodigoAbreviacion:EV", &respuestaPeticion); err == nil && response == 200 {
		helpers.LimpiezaRespuestaRefactor(respuestaPeticion, &listaAsignacionEvaluador)
		if len(listaAsignacionEvaluador) > 0 && listaAsignacionEvaluador[0].EvaluacionId != nil {

			for _, asignacion := range listaAsignacionEvaluador {
				if asignacion.RolAsignacionEvaluadorId.CodigoAbreviacion == "EPR" {
					asignacionesPendienesPorEvaluar = true
					break

				} else {
					asignacionesPendienesPorEvaluar = false
				}
			}

		}
	} else {
		return asignacionesPendienesPorEvaluar, fmt.Errorf("error al consultar asignaciones")

	}
	return asignacionesPendienesPorEvaluar, nil
}

func agregarRolEvaluador(asignaciones []models.AsignacionEvaluador) (listaNoAgregados []models.PeticionAutenticacion, outputError error) {

	for _, asignacion := range asignaciones {

		activo, autenticacion, err := verificarRolEvaluador(asignacion.PersonaId)
		if err != nil {
			return listaNoAgregados, fmt.Errorf("error al verificar rol de evaluador")
		}

		if !activo {
			var autenticacionResponse = make(map[string]interface{})

			peticionAutenticacion := models.PeticionAutenticacion{User: autenticacion.Email, Rol: "EVALUADOR_CUMPLIDO_PROV"}
			if response := helpers.SendJson(beego.AppConfig.String("UrlAutenticacionMid")+"/rol/add", "POST", &autenticacionResponse, peticionAutenticacion); response == nil {
			} else {
				listaNoAgregados = append(listaNoAgregados, peticionAutenticacion)
			}
		}

	}
	return listaNoAgregados, nil
}

func eliminarRolEvaluador(asignaciones []models.AsignacionEvaluador) (listaNoEliminados []models.PeticionAutenticacion, outputError error) {

	for _, asignacion := range asignaciones {
		activo, autenticacion, err := verificarRolEvaluador(asignacion.PersonaId)
		if err != nil {
			return listaNoEliminados, fmt.Errorf("error al verificar rol de evaluador")
		}
		asignacionePendientes, err := consultarAsignacionesPorDocumentoEvaluadorYEstadoEvaluador(asignacion.PersonaId)

		if err != nil {
			return listaNoEliminados, fmt.Errorf("error al verificar asignaciones pendientes")
		}

		if activo && !asignacionePendientes {
			var autenticacionResponse = make(map[string]interface{})
			peticionAutenticacion := models.PeticionAutenticacion{User: autenticacion.Email, Rol: "EVALUADOR_CUMPLIDO_PROV"}
			if response := helpers.SendJson(beego.AppConfig.String("UrlAutenticacionMid")+"/rol/remove", "POST", &autenticacionResponse, peticionAutenticacion); response == nil {
			} else {
				listaNoEliminados = append(listaNoEliminados, peticionAutenticacion)
			}
		}

	}
	return listaNoEliminados, nil
}

func consultarCambioEstadoEvaluacion(idEvalacion string) (cambiosEstadoEvaluacion *models.CambioEstadoEvaluacion, outputError error) {
	defer func() {
		if err := recover(); err != nil {
			outputError = fmt.Errorf("%v", err)
			panic(outputError)
		}
	}()

	resultado_map := make(map[string]interface{})
	var listCambiosEstadoEvaluacion []models.CambioEstadoEvaluacion

	if response, err := helpers.GetJsonWSO2Test(beego.AppConfig.String("UrlEvaluacionCumplidoCrud")+"/cambio_estado_evaluacion/?query=Activo:true,EvaluacionId.id:"+idEvalacion, &resultado_map); err == nil && response == 200 {
		helpers.LimpiezaRespuestaRefactor(resultado_map, &listCambiosEstadoEvaluacion)
		//fmt.Println(listCambiosEstadoEvaluacion)
		cambiosEstadoEvaluacion = &listCambiosEstadoEvaluacion[0]

	} else {
		return nil, fmt.Errorf("Error al consultar asignaciones")
	}

	return cambiosEstadoEvaluacion, nil
}

func verificarRolEvaluador(documento string) (activo bool, autenticacionResponse models.Autentiacion, outputError error) {
	defer func() {
		if err := recover(); err != nil {
			outputError = fmt.Errorf("%v", err)
			panic(outputError)
		}
	}()

	peticionAutenticacion := models.PeticionAutenticacion{Numero: documento}
	autenticacionResponse = models.Autentiacion{}
	//fmt.Println(beego.AppConfig.String("UrlAutenticacionMid"))
	fmt.Print(beego.AppConfig.String("UrlAutenticacionMid") + "/token/documentoToken")
	if response := helpers.SendJson(beego.AppConfig.String("UrlAutenticacionMid")+"/token/documentoToken", "POST", &autenticacionResponse, peticionAutenticacion); response == nil {

		for _, rol := range autenticacionResponse.Role {
			if rol == "EVALUADOR_CUMPLIDO_PROV" {
				activo = true
				break
			} else {
				activo = false
			}
		}

	} else {
		return false, autenticacionResponse, fmt.Errorf("error al consultar roles en autenticacion")
	}

	return activo, autenticacionResponse, nil
}
