package services

import (
	"bufio"
	"bytes"
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"math"
	"strconv"
	"strings"
	"time"

	"github.com/astaxie/beego"
	"github.com/phpdave11/gofpdf"
	"github.com/udistrital/evaluacion_cumplido_prov_mid/helpers"
	"github.com/udistrital/evaluacion_cumplido_prov_mid/models"
	"github.com/udistrital/utils_oas/xlsx2pdf"
	excelize "github.com/xuri/excelize/v2"
)

func GenerarDocumentoEvaluacion(evaluacion_id int) (outputError error) {
	defer func() {
		if err := recover(); err != nil {
			outputError = fmt.Errorf("%v", err)
			panic(outputError)
		}
	}()

	informacion_evaluacion, error_informacion := ObtenerInformacionDocumento(evaluacion_id)
	if error_informacion != nil {
		outputError = fmt.Errorf("Error al obtener la informacion del documento")
		return outputError
	}

	// Abrir el archivo de la plantilla
	filePath := "static/plantilla/plantilla.xlsx"
	f, err := excelize.OpenFile(filePath)
	if err != nil {
		outputError = fmt.Errorf("Error al abrir el archivo")
		return outputError
	}

	sheetName := "GC-PR-006-FR-028"
	var anchoColumna float64

	// Llenar la informacion del proveedor
	fecha_actual := time.Now().Format("02/01/2006")
	f.SetCellValue(sheetName, "E4", "")
	f.SetCellValue(sheetName, "G6", fecha_actual)
	f.SetCellValue(sheetName, "D8", informacion_evaluacion.EmpresaProveedor+" - NIT "+informacion_evaluacion.Documento)
	f.SetCellValue(sheetName, "D6", informacion_evaluacion.Dependencia)
	f.SetCellValue(sheetName, "D9", informacion_evaluacion.ObjetoContrato)

	styleObjetoContrato, err := f.NewStyle(&excelize.Style{
		Alignment: &excelize.Alignment{
			WrapText:   true,
			Vertical:   "center",
			Horizontal: "center",
		},
	})

	if err != nil {
		outputError = fmt.Errorf("Error al crear el estilo")
		return outputError
	}

	if err := f.SetCellStyle(sheetName, "D9", "D9", styleObjetoContrato); err != nil {
		outputError = fmt.Errorf("Error al aplicar el estilo a la celda")
		return outputError
	}

	// Registrar los resultados de la evaluación
	if err := registrarResultadosEvaluacion(f, sheetName, informacion_evaluacion); err != nil {
		outputError = fmt.Errorf("Error al registrar los resultados de la evaluación")
		return outputError
	}

	// Registrar los evaluadores
	if err := registrarEvaluadores(f, sheetName, informacion_evaluacion); err != nil {
		outputError = fmt.Errorf("Error al combinar celdas por evaluador")
		return outputError
	}

	for _, cell := range []string{"D", "E", "F", "G"} {
		widht, err := f.GetColWidth(sheetName, cell)
		if err != nil {
			outputError = fmt.Errorf("Error al obtener el ancho de la columna del objeto contrato")
			return outputError
		}
		anchoColumna += widht
	}

	// Ajustar la altura de la celda del objeto contrato
	error_ajustar_alture := AjustarAlturaCelda(f, sheetName, "D9", informacion_evaluacion.ObjetoContrato, anchoColumna)
	if error_ajustar_alture != nil {
		outputError = fmt.Errorf("Error al ajustar la altura de la celda")
		return outputError
	}

	if err := f.SaveAs(fmt.Sprintf("static/documento/evaluacion_%s.xlsx", informacion_evaluacion.EmpresaProveedor)); err != nil {
		outputError = fmt.Errorf("Error al guardar el archivo")
		return outputError
	}

	// Crear Pdf

	template, err := excelize.OpenFile(fmt.Sprintf("static/documento/evaluacion_%s.xlsx", informacion_evaluacion.EmpresaProveedor))
	if err != nil {
		fmt.Println(err)
		return
	}

	pdf := gofpdf.New("P", "mm", "A4", "")

	ExcelPdf := xlsx2pdf.Excel2PDF{
		Excel:    template,
		Pdf:      pdf,
		Sheets:   make(map[string]xlsx2pdf.SheetInfo),
		WFx:      2,
		HFx:      2.925,
		FontDims: xlsx2pdf.FontDims{Size: 0.85},
		Header:   func() {},
		Footer:   func() {},
		CustomSize: xlsx2pdf.PageFormat{
			Orientation: "P",
			Wd:          297,
			Ht:          210,
		},
	}

	dim, _ := template.GetSheetDimension(sheetName)
	_, maxrow, _ := excelize.CellNameToCoordinates(strings.Split(dim, ":")[1])
	for r := 1; r <= maxrow; r++ {
		h, _ := template.GetRowHeight(sheetName, r)
		template.SetRowHeight(sheetName, r, h*1.046)
	}

	ExcelPdf.ConvertSheets()

	// Codificar el archivo en base64
	encodeFile, outputError := encodePDF(pdf)
	if outputError != nil {
		return outputError
	}

	// Decodificar el archivo en base64 y guardarlo
	decoded, err := base64.StdEncoding.DecodeString(encodeFile)
	if err != nil {
		return fmt.Errorf("Error al decodificar el archivo")
	}

	// Guardar el archivo
	err = ioutil.WriteFile("static/documento/evaluacion_"+informacion_evaluacion.EmpresaProveedor+".pdf", decoded, 0644)
	if err != nil {
		return fmt.Errorf("Error al guardar el archivo")
	}

	// if err := pdf.OutputFileAndClose("static/documento/prueba.pdf"); err != nil {
	// 	fmt.Println(err)
	// 	return
	// }

	return nil
}

func encodePDF(pdf *gofpdf.Fpdf) (encodedFile string, outputError error) {
	defer func() {
		if err := recover(); err != nil {
			outputError = fmt.Errorf("%v", err)
			panic(outputError)
		}
	}()
	var buffer bytes.Buffer
	writer := bufio.NewWriter(&buffer)
	//pdf.OutputFileAndClose("/home/faidercamilo/go/src/github.com/udistrital/prueba.pdf") // para guardar el archivo localmente
	err := pdf.Output(writer)
	if err != nil {
		outputError = fmt.Errorf("Error al generar el PDF:", err)
		return encodedFile, outputError
	}
	err = writer.Flush()
	if err != nil {
		outputError = fmt.Errorf("Error al hacer flush del writer:", err)
	}
	encodedFile = base64.StdEncoding.EncodeToString(buffer.Bytes())
	//fmt.Println(encodedFile)
	return encodedFile, nil
}

func ExcelToPDFWithStyles(inputFile, outputFile string) error {
	// Abrir el archivo Excel
	excel, err := excelize.OpenFile(inputFile)
	if err != nil {
		return fmt.Errorf("Error al abrir el archivo Excel: %v", err)
	}

	// Crear un nuevo PDF
	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.SetFont("Arial", "", 10)

	for _, sheet := range excel.GetSheetList() {
		// Crear una nueva página para cada hoja
		pdf.AddPage()
		pdf.SetFont("Arial", "B", 12)
		pdf.Cell(190, 10, fmt.Sprintf("Hoja: %s", sheet))
		pdf.Ln(12)

		// Obtener las filas de la hoja
		rows, err := excel.GetRows(sheet)
		if err != nil {
			return fmt.Errorf("Error al leer filas de la hoja %s: %v", sheet, err)
		}

		// Ajustar el ancho y alto de las celdas
		pageWidth := 190.0 // Ancho total de la página
		colWidths := pageWidth / float64(len(rows[0]))
		rowHeight := 8.0

		// Procesar celdas combinadas
		mergedCells, err := excel.GetMergeCells(sheet)
		if err != nil {
			return fmt.Errorf("Error al obtener celdas combinadas: %v", err)
		}
		mergedRanges := map[string]string{}
		for _, merged := range mergedCells {
			mergedRanges[merged.GetStartAxis()] = merged.GetEndAxis()
		}

		// Procesar cada fila
		for rIdx, row := range rows {
			for cIdx, cellValue := range row {
				cellName, _ := excelize.CoordinatesToCellName(cIdx+1, rIdx+1)
				styleID, _ := excel.GetCellStyle(sheet, cellName)

				// Obtener estilos
				style, err := excel.GetStyle(styleID)
				if err != nil {
					return fmt.Errorf("Error al obtener estilo de celda: %v", err)
				}

				// Procesar colores
				fillColor := "255 255 255" // Blanco por defecto
				if len(style.Fill.Color) > 0 {
					fillColor = style.Fill.Color[0]
				}
				r, g, b := parseHexColor(fillColor)

				// Dibujar la celda
				pdf.SetFillColor(r, g, b)
				pdf.SetTextColor(0, 0, 0) // Negro por defecto
				pdf.SetDrawColor(0, 0, 0)

				if _, isMerged := mergedRanges[cellName]; isMerged {
					// Si la celda es parte de un rango combinado, saltar para evitar duplicados
					continue
				}

				// Ajustar texto
				text := cellValue
				if len(text) > 50 {
					text = fmt.Sprintf("%s...", text[:50])
				}

				// Dibujar la celda con texto y bordes
				pdf.CellFormat(colWidths, rowHeight, text, "1", 0, "C", true, 0, "")
			}
			pdf.Ln(rowHeight) // Mover a la siguiente fila
		}
	}

	// Guardar el archivo PDF
	if err := pdf.OutputFileAndClose(outputFile); err != nil {
		return fmt.Errorf("Error al guardar el archivo PDF: %v", err)
	}

	return nil
}

func parseHexColor(hexColor string) (int, int, int) {
	hexColor = strings.TrimPrefix(hexColor, "#")
	if len(hexColor) == 6 {
		var r, g, b int
		fmt.Sscanf(hexColor, "%02x%02x%02x", &r, &g, &b)
		return r, g, b
	}
	return 255, 255, 255
}

func registrarResultadosEvaluacion(f *excelize.File, sheetName string, informacion_evaluacion models.InformacionDocumentoEvaluacion) error {

	startRow := 14
	endRow := 34

	var preguntaAnterior string
	row := startRow

	for row <= endRow {

		var preguntaActual string
		//var categoria string
		var titulo string

		//Recuperar la pregunta de la fila actual
		preguntaActual, err := f.GetCellValue(sheetName, fmt.Sprintf("D%d", row))
		if err != nil {
			return fmt.Errorf("Error al obtener el valor de la celda: %w", err)
		}

		//Recuperar el titulo de la fila actual
		titulo, err = f.GetCellValue(sheetName, fmt.Sprintf("C%d", row))
		if err != nil {
			return fmt.Errorf("Error al obtener el valor de la celda: %w", err)
		}

		//Recuperar la categoria de la fila actual
		// categoria, err = f.GetCellValue(sheetName, fmt.Sprintf("B%d", row))
		// if err != nil {
		// 	return fmt.Errorf("Error al obtener el valor de la celda: %w", err)
		// }

		preguntaActual = limpiarPregunta(preguntaActual)
		//fmt.Println("Pregunta actual: ", preguntaActual)
		// fmt.Printf("Cadena original: %q\n", preguntaActual)
		// for i, c := range preguntaActual {
		// 	fmt.Printf("Posición %d: '%c' (U+%04X)\n", i, c, c)
		// }

		//Verificar si la pregunta anterior es igual a la actual se debe saltar a la siguiente
		if preguntaAnterior == preguntaActual && preguntaAnterior != "" {
			row += 2
			continue
		}

		preguntaAnterior = preguntaActual
		// Buscar la pregunta en la información de la evaluación
		for _, evaluacion := range informacion_evaluacion.ResultadoFinalEvaluacion.Resultados {
			if strings.ToLower(strings.TrimSpace(evaluacion.Pregunta)) == strings.ToLower(preguntaActual) &&
				strings.ToLower(strings.TrimSpace(evaluacion.Titulo)) == strings.ToLower(strings.TrimSpace(titulo)) {
				// fmt.Println("------------------------- Resultados -----------------------------------------------")
				// fmt.Println("Pregunta actual: ", strings.ToLower(preguntaActual))
				// fmt.Println("Pregunta evaluacion: ", strings.ToLower(strings.TrimSpace(evaluacion.Pregunta)))
				// fmt.Println("Titulo evaluacion: ", strings.ToLower(strings.TrimSpace(evaluacion.Titulo)))
				// fmt.Println("Titulo actual: ", strings.ToLower(strings.TrimSpace(titulo)))
				// fmt.Println("Resultado: ", evaluacion.Cumplimiento)
				// fmt.Println("Pregunta: ", strings.ToLower(strings.TrimSpace(evaluacion.Pregunta)) == strings.ToLower(preguntaActual))
				// fmt.Println("Titulo: ", strings.ToLower(strings.TrimSpace(evaluacion.Titulo)) == strings.ToLower(strings.TrimSpace(titulo)))
				// fmt.Println("------------------------------------------------------------------------------------")

				err := f.SetCellValue(sheetName, fmt.Sprintf("E%d", row), evaluacion.Cumplimiento)
				if err != nil {
					return fmt.Errorf("Error al establecer el valor de la celda: %w", err)
				}

				err_valor := f.SetCellValue(sheetName, fmt.Sprintf("F%d", row), evaluacion.ValorAsignado)
				if err_valor != nil {
					return fmt.Errorf("Error al establecer el valor de la celda: %w", err_valor)
				}

			}
		}
		row += 2
	}
	return nil
}

func limpiarPregunta(pregunta string) string {
	pregunta = strings.TrimSpace(pregunta)
	if pregunta[0] == '(' {
		pregunta = pregunta[6:]
	}
	return pregunta
}

func registrarEvaluadores(f *excelize.File, sheetName string, informacion_evaluacion models.InformacionDocumentoEvaluacion) error {
	rowStart := 10 // La fila inicial para insertar
	rowSupervisor := 10
	rowOffset := 2 // Número de filas a insertar cada vez
	//colores_fondo := []string{"#FDE9D9", "#EBF1DE", "#DCE6F1", "#F1F1F1"}
	//indice_color := 0
	anchoColumna := 0.0
	for _, evaluador := range informacion_evaluacion.ResultadoFinalEvaluacion.Evaluadores {

		if strings.ToLower(evaluador.Rol) == "evaluador" {
			//Configurar los campos de cada evaluador

			// Insertar filas
			err := f.InsertRows(sheetName, rowStart, rowOffset)
			if err != nil {
				return fmt.Errorf("error al insertar las filas: %w", err)
			}

			// Combinar columnas específicas
			for i := 0; i < rowOffset; i++ {
				currentRow := rowStart + i

				// Combinar columnas B y C
				err = f.MergeCell(sheetName, fmt.Sprintf("B%d", currentRow), fmt.Sprintf("C%d", currentRow))
				if err != nil {
					return fmt.Errorf("error al combinar columnas B y C en la fila %d: %w", currentRow, err)
				}

				styleId, err := f.GetCellStyle(sheetName, fmt.Sprintf("B%d", currentRow+2))
				if err != nil {
					return fmt.Errorf("error al obtener el estilo de la celda: %w", err)
				}

				style, err := f.GetStyle(styleId)
				if err != nil {
					return fmt.Errorf("error al obtener el estilo: %w", err)
				}

				// style.Fill = excelize.Fill{
				// 	Type:    "pattern",
				// 	Color:   []string{colores_fondo[indice_color]},
				// 	Pattern: 1,
				// }

				newStyleId, err := f.NewStyle(style)

				// Copiar el estilo de la celda
				err = f.SetCellStyle(sheetName, fmt.Sprintf("B%d", currentRow), fmt.Sprintf("G%d", currentRow), newStyleId)
				if err != nil {
					return fmt.Errorf("error al aplicar el estilo a la celda: %w", err)
				}

				if i == 0 {
					//Configurar la columa de los items evaluados

					altura, err := f.GetRowHeight(sheetName, currentRow+2)
					if err != nil {
						return fmt.Errorf("error al obtener la altura de la fila: %w", err)
					}
					// Ajustar la altura de la fila
					if err := f.SetRowHeight(sheetName, currentRow, altura); err != nil {
						return fmt.Errorf("Error al ajustar la altura de la fila: %v", err)
					}

					// Combinar columnas B y C
					err = f.MergeCell(sheetName, fmt.Sprintf("B%d", currentRow), fmt.Sprintf("C%d", currentRow))
					if err != nil {
						return fmt.Errorf("error al combinar columnas B y C en la fila %d: %w", currentRow, err)
					}

					// Combinar columnas D y E
					err = f.MergeCell(sheetName, fmt.Sprintf("D%d", currentRow), fmt.Sprintf("E%d", currentRow))
					if err != nil {
						return fmt.Errorf("error al combinar columnas D y E en la fila %d: %w", currentRow, err)
					}

					// Combinar columnas E y F
					err = f.MergeCell(sheetName, fmt.Sprintf("E%d", currentRow), fmt.Sprintf("F%d", currentRow))
					if err != nil {
						return fmt.Errorf("error al combinar columnas E y G en la fila %d: %w", currentRow, err)
					}

					// Combinar columnas F y G
					err = f.MergeCell(sheetName, fmt.Sprintf("F%d", currentRow), fmt.Sprintf("G%d", currentRow))
					if err != nil {
						return fmt.Errorf("error al combinar columnas G y I en la fila %d: %w", currentRow, err)
					}

					f.SetCellValue(sheetName, fmt.Sprintf("B%d", currentRow), "ITEM EVALUADO (*)")
					for _, cell := range []string{"D", "E", "F", "G"} {
						widht, err := f.GetColWidth(sheetName, cell)
						if err != nil {
							return fmt.Errorf("Error al obtener el ancho de la columna del objeto contrato")
						}
						anchoColumna += widht
					}

					AjustarAlturaCelda(f, sheetName, fmt.Sprintf("D%d", currentRow), evaluador.Items, anchoColumna)
					f.SetCellValue(sheetName, fmt.Sprintf("D%d", currentRow), evaluador.Items)

				}

				if i == 1 {
					//Configurar columna del nombre evaluador y cargo

					altura, err := f.GetRowHeight(sheetName, currentRow+2)
					if err != nil {
						return fmt.Errorf("error al obtener la altura de la fila: %w", err)
					}
					// Ajustar la altura de la fila
					if err := f.SetRowHeight(sheetName, currentRow, altura); err != nil {
						return fmt.Errorf("Error al ajustar la altura de la fila: %v", err)
					}
					columnas := []string{"D", "E", "F", "G"}
					for _, columna := range columnas {
						// Obtener el valor de la celda original
						celdaOrigen := fmt.Sprintf("%s%d", columna, currentRow+2) // Fila de referencia para copiar
						valor, err := f.GetCellValue(sheetName, celdaOrigen)
						if err != nil {
							return fmt.Errorf("error al obtener el valor de la celda %s: %w", celdaOrigen, err)
						}

						// Copiar el valor a la nueva celda
						celdaDestino := fmt.Sprintf("%s%d", columna, currentRow) // Nueva fila
						if err := f.SetCellValue(sheetName, celdaDestino, valor); err != nil {
							return fmt.Errorf("error al copiar el valor a la celda %s: %w", celdaDestino, err)
						}

						// Copiar el estilo de la celda original
						estilo, err := f.GetCellStyle(sheetName, celdaOrigen)
						if err != nil {
							return fmt.Errorf("error al obtener el estilo de la celda %s: %w", celdaOrigen, err)
						}

						style, err := f.GetStyle(estilo)
						if err != nil {
							return fmt.Errorf("error al obtener el estilo: %w", err)
						}

						// style.Fill = excelize.Fill{
						// 	Type:    "pattern",
						// 	Color:   []string{colores_fondo[indice_color]},
						// 	Pattern: 1,
						// }

						newStyleId, err := f.NewStyle(style)

						if err := f.SetCellStyle(sheetName, celdaDestino, celdaDestino, newStyleId); err != nil {
							return fmt.Errorf("error al copiar el estilo a la celda %s: %w", celdaDestino, err)
						}
					}
					f.SetCellValue(sheetName, fmt.Sprintf("B%d", currentRow), "OBJETO DEL CONTRATO:")
					anchoColumna, err := f.GetColWidth(sheetName, "E")
					AjustarAlturaCelda(f, sheetName, fmt.Sprintf("E%d", currentRow), fmt.Sprintf("CARGO: %s", evaluador.Cargo), anchoColumna)
					f.SetCellValue(sheetName, fmt.Sprintf("E%d", currentRow), fmt.Sprintf("CARGO: %s", evaluador.Cargo))
					f.SetCellValue(sheetName, fmt.Sprintf("D%d", currentRow), evaluador.Nombre)
				}

			}

			// Actualizar fila inicial para la próxima inserción y la fila del supervisor
			rowStart += rowOffset
			rowSupervisor += rowOffset - 1
		} else {
			// Configurar los campos del supervisor

			// Configurar el estilo celda item evaluado
			estilo, err := f.GetCellStyle(sheetName, fmt.Sprintf("B%d", rowSupervisor))
			if err != nil {
				return fmt.Errorf("Error al obtener el estilo de la celda: %v", err)
			}

			style, err := f.GetStyle(estilo)

			// style.Fill = excelize.Fill{
			// 	Type:    "pattern",
			// 	Color:   []string{colores_fondo[indice_color]},
			// 	Pattern: 1,
			// }

			style.Alignment = &excelize.Alignment{
				WrapText:   true,
				Vertical:   "center",
				Horizontal: "center",
			}

			newStyleId, err := f.NewStyle(style)

			if err != nil {
				return fmt.Errorf("Error al crear el estilo: %v", err)
			}

			if err := f.SetCellStyle(sheetName, fmt.Sprintf("B%d", rowSupervisor), fmt.Sprintf("B%d", rowSupervisor), newStyleId); err != nil {
				return fmt.Errorf("Error al aplicar el estilo a la celda: %v", err)
			}

			// Configurar celda de objeto contrato
			for _, cell := range []string{"D", "E", "F", "G"} {
				estilo, err := f.GetCellStyle(sheetName, fmt.Sprintf("%s%d", cell, rowSupervisor))
				if err != nil {
					return fmt.Errorf("Error al obtener el estilo de la celda: %v", err)
				}

				style, err := f.GetStyle(estilo)
				if err != nil {
					return fmt.Errorf("Error al obtener el estilo: %v", err)
				}

				// style.Fill = excelize.Fill{
				// 	Type:    "pattern",
				// 	Color:   []string{colores_fondo[indice_color]},
				// 	Pattern: 1,
				// }

				style.Alignment = &excelize.Alignment{
					WrapText:   true,
					Vertical:   "center",
					Horizontal: "center",
				}

				newStyleId, err := f.NewStyle(style)
				if err != nil {
					return fmt.Errorf("Error al crear el estilo: %v", err)
				}

				if err := f.SetCellStyle(sheetName, fmt.Sprintf("%s%d", cell, rowSupervisor), fmt.Sprintf("%s%d", cell, rowSupervisor), newStyleId); err != nil {
					return fmt.Errorf("Error al aplicar el estilo a la celda: %v", err)
				}

			}

			// Configurar el estilo de la la casilla NOMBRE DEL SUPERVISOR ENCARGADO DE LA EVALUACIÓN:
			estilo2, err := f.GetCellStyle(sheetName, fmt.Sprintf("B%d", rowSupervisor+1))
			if err != nil {
				return fmt.Errorf("Error al obtener el estilo de la celda: %v", err)
			}

			style2, err := f.GetStyle(estilo2)

			// style2.Fill = excelize.Fill{
			// 	Type:    "pattern",
			// 	Color:   []string{colores_fondo[indice_color]},
			// 	Pattern: 1,
			// }

			newStyleId2, err := f.NewStyle(style2)

			if err != nil {
				return fmt.Errorf("Error al crear el estilo: %v", err)
			}

			if err := f.SetCellStyle(sheetName, fmt.Sprintf("B%d", rowSupervisor+1), fmt.Sprintf("B%d", rowSupervisor+1), newStyleId2); err != nil {
				return fmt.Errorf("Error al aplicar el estilo a la celda: %v", err)
			}

			// Configurar el estilo de las celdas de nombre y cargo del supervisor
			for _, cell := range []string{"D", "E", "F", "G"} {
				estilo, err := f.GetCellStyle(sheetName, fmt.Sprintf("%s%d", cell, rowSupervisor+1))
				if err != nil {
					return fmt.Errorf("Error al obtener el estilo de la celda: %v", err)
				}

				style, err := f.GetStyle(estilo)
				if err != nil {
					return fmt.Errorf("Error al obtener el estilo: %v", err)
				}

				// style.Fill = excelize.Fill{
				// 	Type:    "pattern",
				// 	Color:   []string{colores_fondo[indice_color]},
				// 	Pattern: 1,
				// }

				newStyleId, err := f.NewStyle(style)
				if err != nil {
					return fmt.Errorf("Error al crear el estilo: %v", err)
				}

				if err := f.SetCellStyle(sheetName, fmt.Sprintf("%s%d", cell, rowSupervisor+1), fmt.Sprintf("%s%d", cell, rowSupervisor+1), newStyleId); err != nil {
					return fmt.Errorf("Error al aplicar el estilo a la celda: %v", err)
				}

			}

			anchoColumna, err := f.GetColWidth(sheetName, "E")
			AjustarAlturaCelda(f, sheetName, fmt.Sprintf("E%d", rowSupervisor+1), "CARGO: "+evaluador.Cargo, anchoColumna)
			f.SetCellValue(sheetName, fmt.Sprintf("E%d", rowSupervisor+1), "CARGO: "+evaluador.Cargo)
			f.SetCellValue(sheetName, fmt.Sprintf("D%d", rowSupervisor+1), evaluador.Nombre)
			f.SetCellValue(sheetName, fmt.Sprintf("D%d", rowSupervisor), evaluador.Items)

		}

		// if (indice_color + 1) < len(colores_fondo)-1 {
		// 	indice_color++
		// } else {
		// 	indice_color = 0
		// }
	}

	return nil
}

func AjustarAlturaCelda(f *excelize.File, sheetName, cell, texto string, anchoColumna float64) error {
	// Calcular el número de caracteres aproximados por línea
	caracteresPorLinea := int(anchoColumna * 0.8)

	// Calcular el número de líneas necesarias para el texto
	lineas := calcularLineasTexto(texto, caracteresPorLinea) + 1 // Añadir una línea adicional para evitar cortes

	altura := float64(lineas) * 12.0

	// Obtener la fila de la celda
	row, _ := strconv.Atoi(cell[1:])

	// Ajustar la altura de la fila
	if err := f.SetRowHeight(sheetName, row, altura); err != nil {
		return fmt.Errorf("Error al ajustar la altura de la fila: %v", err)
	}

	// Crear y aplicar el estilo de la celda con wrap text
	// style, err := f.NewStyle(&excelize.Style{
	// 	Alignment: &excelize.Alignment{
	// 		WrapText:   true,     // Permitir ajuste automático del texto
	// 		Vertical:   "center", // Centrar verticalmente
	// 		Horizontal: "center", // Centrar horizontalmente
	// 	},
	// })
	// if err != nil {
	// 	return fmt.Errorf("Error al crear el estilo: %v", err)
	// }

	// if err := f.SetCellStyle(sheetName, cell, cell, style); err != nil {
	// 	return fmt.Errorf("Error al aplicar el estilo a la celda: %v", err)
	// }

	return nil
}

func calcularLineasTexto(texto string, caracteresPorLinea int) int {
	if caracteresPorLinea <= 0 {
		return 1
	}

	// Calcular el número total de caracteres en el texto
	totalCaracteres := len(texto)

	// Calcular cuántas líneas son necesarias para mostrar todo el texto
	lineas := int(math.Ceil(float64(totalCaracteres) / float64(caracteresPorLinea)))

	// Garantizar al menos una línea
	if lineas == 0 {
		lineas = 1
	}

	return lineas
}

func ObtenerInformacionDocumento(evaluacion_id int) (informacion_evaluacion models.InformacionDocumentoEvaluacion, outputError error) {
	defer func() {
		if err := recover(); err != nil {
			outputError = fmt.Errorf("%v", err)
			panic(outputError)
		}
	}()

	// Obtener la evaluacion
	var resultado_evaluacion map[string]interface{}
	var evaluacion []models.Evaluacion
	//fmt.Println("Url evaluacion: ", beego.AppConfig.String("UrlEvaluacionesCumplidosProveedoresCrud")+"/evaluacion/?query=Id:"+strconv.Itoa(evaluacion_id)+",Activo:true&limit=-1")
	if response, err := helpers.GetJsonTest(beego.AppConfig.String("UrlEvaluacionesCumplidosProveedoresCrud")+"/evaluacion/?query=Id:"+strconv.Itoa(evaluacion_id)+",Activo:true&limit=-1", &resultado_evaluacion); err != nil && response != 200 {
		outputError = fmt.Errorf("Error al obtener la evaluacion")
		return informacion_evaluacion, outputError
	}

	data := resultado_evaluacion["Data"].([]interface{})
	if len(data[0].(map[string]interface{})) == 0 {
		outputError = fmt.Errorf("La evaluacion con id %v no se encuentra activa o no existe", evaluacion_id)
		return informacion_evaluacion, outputError
	}

	helpers.LimpiezaRespuestaRefactor(resultado_evaluacion, &evaluacion)

	// Obtener el contrato suscrito
	contrato_general, err := helpers.ObtenerContratoGeneral(strconv.Itoa(evaluacion[0].ContratoSuscritoId), strconv.Itoa(evaluacion[0].VigenciaContrato))
	if err != nil {
		outputError = fmt.Errorf("Error al obtener el contrato general")
		return informacion_evaluacion, outputError
	}

	// Obtener la dependencia del contrato
	dependencia_supervisor, error_contrato := helpers.ObtenerDependenciasSupervisor(strconv.Itoa(contrato_general.Supervisor.Documento))
	if error_contrato != nil {
		outputError = fmt.Errorf("Error al obtener las dependencias del supervisor")
		return informacion_evaluacion, outputError
	}

	// Obtener el resultado final de la evaluacion
	resultado_final_evaluacion, error_resultado := ObtenerResultadoFinalEvaluacion(evaluacion_id)
	if error_resultado != nil {
		outputError = fmt.Errorf("Error al obtener el resultado final de la evaluacion")
		return informacion_evaluacion, outputError
	}

	//fmt.Println("Resultado final evaluacion: ", resultado_final_evaluacion)

	//Obtener los datos del proveedor
	var informacion_proveedor []models.InformacionProveedor
	if response, err := helpers.GetJsonTest(beego.AppConfig.String("UrlAmazonApi")+"/informacion_proveedor/?query=Id:"+strconv.Itoa(contrato_general.Contratista), &informacion_proveedor); (err != nil) && (response != 200) {
		outputError = fmt.Errorf("Error al obtener la informacion del proveedor")
		return informacion_evaluacion, outputError
	}

	for _, dependencia := range dependencia_supervisor {
		if contrato_general.Supervisor.DependenciaSupervisor == dependencia.Codigo {
			informacion_evaluacion.Dependencia = dependencia.Nombre
		}
	}

	informacion_evaluacion.EmpresaProveedor = informacion_proveedor[0].NomProveedor
	informacion_evaluacion.ObjetoContrato = contrato_general.ObjetoContrato
	informacion_evaluacion.ResultadoFinalEvaluacion = resultado_final_evaluacion
	informacion_evaluacion.Documento = informacion_proveedor[0].NumDocumento

	return informacion_evaluacion, nil

}
