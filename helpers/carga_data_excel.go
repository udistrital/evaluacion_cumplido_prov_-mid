package helpers

import (
	"fmt"
	"mime/multipart"
	"strconv"

	"github.com/udistrital/evaluacion_cumplido_prov_mid/models"
	"github.com/udistrital/evaluacion_cumplido_prov_mid/services"
	"github.com/xuri/excelize/v2"
)

// / CargaDataExcel lee  y  carga la data de un archivo excel
func CargaDataExcel(excel multipart.File) (response map[string]interface{}, itemsNoAGregados []models.ItemEvaluacion, outputError error) {

	defer func() {
		if err := recover(); err != nil {
			outputError = fmt.Errorf("%v", err)
			panic(outputError)
		}
	}()

	f, err := excelize.OpenReader(excel)

	if err != nil {
		outputError = fmt.Errorf("Error al abrir el archivo: %v", err)
		return nil, nil, outputError
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
		posicionTipoNecessidad := "G"
		posicionFichaTecnica := "H"

		valorUnidad := obtenerUnidadMedida(obtenerCelda(f, posicionUnidad, i))
		valorTipoNecessidad := obtenerTipoNecessidad(obtenerCelda(f, posicionTipoNecessidad, i))

		item := models.ItemEvaluacion{
			Identificador: obtenerCelda(f, posicionIdentificador, i),
			Nombre:        obtenerCelda(f, posicionNombre, i),
			Cantidad:      parseFloat(obtenerCelda(f, posicionCantidad, i), 64),
			ValorInitario: parseFloat(obtenerCelda(f, posicionValorInitario, i), 64),
			Iva:           parseFloat(obtenerCelda(f, posicionIva, i), 64),
			FichaTecnica:  obtenerCelda(f, posicionFichaTecnica, i),
			EvaluacionId: models.Evaluacion{
				Id: 1,
			},
		}

		if valorUnidad != 0 {
			item.Unidad = valorUnidad
		}

		if valorTipoNecessidad != 0 {
			item.TipoNecessidad = valorTipoNecessidad
		}

		if verificarExistencia(itemsAAGregar, item.Identificador) {
			itemsAAGregar = append(itemsAAGregar, item)
		} else {
			itemsNoAGregados = append(itemsNoAGregados, item)
		}

	}

	response, err = services.GuardarItem(itemsAAGregar)

	if err != nil {
		outputError = fmt.Errorf("Error al guardar item: %v", err)
		return nil, nil, outputError
	}
	return response, itemsNoAGregados, nil
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
	unidadRespuesta, _ := services.ObternerUnidadMedida(unidad)

	if unidadRespuesta != nil {
		idUnidad = *unidadRespuesta
		return idUnidad
	}
	return 0
}

func obtenerTipoNecessidad(tipoNecessidad string) (idTipoNecessidad int) {

	if tipoNecessidad == "BIEN" {
		return 1
	}
	if tipoNecessidad == "SERVICIO" {
		return 2
	}
	if tipoNecessidad == "BIENES Y SERVICIOS" {
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
