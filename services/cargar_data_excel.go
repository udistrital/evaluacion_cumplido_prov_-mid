package services

import (
	"fmt"
	"mime/multipart"
	"strconv"
	"strings"

	"github.com/astaxie/beego"
	"github.com/udistrital/evaluacion_cumplido_prov_mid/helpers"
	"github.com/udistrital/evaluacion_cumplido_prov_mid/models"
	"github.com/xuri/excelize/v2"
)

// / CargaDataExcel lee  y  carga la data de un archivo excel
func CargaDataExcel(excel multipart.File, evaluacionId int) (response string, itemsNoAGregados []models.ItemEvaluacion, outputError error) {

	defer func() {
		if err := recover(); err != nil {
			outputError = fmt.Errorf("%v", err)
			panic(outputError)
		}
	}()

	f, err := excelize.OpenReader(excel)

	if err != nil {
		outputError = fmt.Errorf("Error al abrir el archivo: %v", err)
		return "", nil, outputError
	}

	var itemsAAGregar []models.ItemEvaluacion

	for i := 2; ; i++ {

		Id := obtenerCelda(f, "A", i)
		if Id == "" {
			break
		}

		posicionIdentificador := "A"
		posicionNombre := "B"
		posicionCantidad := "C"
		posicionValorInitario := "D"
		posicionIva := "E"
		posicionUnidad := "F"
		posicionTipoNecesidad := "G"
		posicionFichaTecnica := "H"

		valorUnidad := obtenerUnidadMedida(obtenerCelda(f, posicionUnidad, i))
		valorTipoNecesidad := obtenerTipoNecesidad(obtenerCelda(f, posicionTipoNecesidad, i))

		item := models.ItemEvaluacion{
			Identificador: obtenerCelda(f, posicionIdentificador, i),
			Nombre:        obtenerCelda(f, posicionNombre, i),
			Cantidad:      parseFloat(obtenerCelda(f, posicionCantidad, i), 64),
			ValorUnitario: parseFloat(obtenerCelda(f, posicionValorInitario, i), 64),
			Iva:           parseFloat(obtenerCelda(f, posicionIva, i), 64),
			FichaTecnica:  obtenerCelda(f, posicionFichaTecnica, i),
			EvaluacionId: models.Evaluacion{
				Id: evaluacionId,
			},
		}

		if valorUnidad != 0 {
			item.Unidad = valorUnidad
		}

		if valorTipoNecesidad != 0 {
			item.TipoNecesidad = valorTipoNecesidad
		}

		fmt.Printf("Item #%d: %+v\n", i-1, item)

		if verificarExistencia(itemsAAGregar, item.Identificador) {
			itemsAAGregar = append(itemsAAGregar, item)
		} else {
			itemsNoAGregados = append(itemsNoAGregados, item)
		}

	}

	responseCrud, err := GuardarItems(itemsAAGregar)

	if err != nil {
		outputError = fmt.Errorf("Error al guardar item: %v", err)
		return "", nil, outputError
	}
	return responseCrud["Message"].(string), itemsNoAGregados, nil
}

func obtenerCelda(excel *excelize.File, LetraColumna string, indexCelda int) (cell string) {

	cell, err := excel.GetCellValue("Informacion", LetraColumna+strconv.Itoa(indexCelda))

	if err != nil {
		fmt.Errorf("Error al leer celda:  %v", err)
	}

	return cell
}

func parseFloat(s string, bitSize int) float64 {
	f, _ := strconv.ParseFloat(s, bitSize)
	return f
}

func obtenerUnidadMedida(unidad string) (idUnidad int) {
	unidadRespuesta, _ := ObternerUnidadMedida(unidad)

	if unidadRespuesta != 0 {
		idUnidad = unidadRespuesta
		return idUnidad
	}
	return 0
}

func obtenerTipoNecesidad(tipoNecesidad string) (idTipoNecesidad int) {
	if strings.ToLower(tipoNecesidad) == "bien" {
		return 1
	}
	if strings.ToLower(tipoNecesidad) == "servicio" {
		return 2
	}
	if strings.ToLower(tipoNecesidad) == "bien/servicio" {
		return 3
	}
	return 0
}

func verificarExistencia(listaItems []models.ItemEvaluacion, identificadorItem string) bool {

	for _, item := range listaItems {
		if item.Identificador == identificadorItem {
			return false
		}
	}
	return true

}

func ObternerUnidadMedida(unidad string) (idUnidad int, outputError error) {
	defer func() {
		if err := recover(); err != nil {
			outputError = fmt.Errorf("%v", err)
			panic(outputError.Error())
		}
	}()

	var unidadMedida []models.UnidadMedida

	if response, err := helpers.GetJsonWSO2Test(beego.AppConfig.String("UrlAdministrativaAmazonApi")+"/unidad", &unidadMedida); err != nil && response != 200 {
		outputError = fmt.Errorf("Error al obtener las unidades de medida")
		return 0, outputError
	}

	for _, unidadMedidaItem := range unidadMedida {
		if strings.ToLower(unidadMedidaItem.Unidad) == strings.ToLower(unidad) {
			return unidadMedidaItem.Id, nil
		}
	}
	return 0, outputError
}

func GuardarItems(items []models.ItemEvaluacion) (response map[string]interface{}, err error) {

	var respuesta_peticion map[string]interface{}

	if err := helpers.SendJson(beego.AppConfig.String("UrlEvaluacionCumplidoCrud")+"/item/guardado_multiple", "POST", &respuesta_peticion, items); err == nil {

		return respuesta_peticion, nil
	}

	return nil, err
}
