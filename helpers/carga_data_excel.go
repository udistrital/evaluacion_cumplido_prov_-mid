package helpers

import (
	"fmt"
	"mime/multipart"
	"strconv"

	"github.com/udistrital/evaluacion_cumplido_prov_mid/models"
	"github.com/xuri/excelize/v2"
)

// / CargaDataExcel lee  y  carga la data de un archivo excel
func CargaDataExcel(excel multipart.File) {

	f, err := excelize.OpenReader(excel)

	if err != nil {
		fmt.Errorf("Error al abrir el archivo: %v", err)
	}

	var items []models.ItemEvaluacion

	for i := 2; ; i++ {

		Id := obtenerCelda(f, 65, i)
		if Id == "" {
			break
		}

		item := models.ItemEvaluacion{
			Id:             obtenerCelda(f, 65, i),
			Nombre:         obtenerCelda(f, 66, i),
			Cantidad:       parseFloat(obtenerCelda(f, 67, i), 64),
			ValorInitario:  parseFloat(obtenerCelda(f, 68, i), 64),
			Iva:            parseFloat(obtenerCelda(f, 69, i), 64),
			Unidad:         1,
			TipoNecessidad: 1,
			FichaTecnica:   obtenerCelda(f, 72, i),
		}
		items = append(items, item)
		fmt.Printf("%+v\n", items)
		// return nil
	}
}

func obtenerCelda(excel *excelize.File, indexLetra int, indexCelda int) (cell string) {

	cell, err := excel.GetCellValue("Informacion", obtenerColumnna(indexLetra)+strconv.Itoa(indexCelda))

	if err != nil {
		fmt.Errorf("Error al leer celda:  %v", err)
	}

	return cell
}

func obtenerColumnna(index int) (columna string) {

	return string(rune(index))
}

func parseFloat(s string, bitSize int) float64 {
	f, _ := strconv.ParseFloat(s, bitSize)
	return f
}

func parseInt(s string, bitSize int) int {
	i, _ := strconv.Atoi(s)
	return i
}
