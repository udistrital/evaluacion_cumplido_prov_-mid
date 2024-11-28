// @APIVersion 1.0.0
// @Title API Test
// @Description API para responder con 'Pong' al hacer un GET a /ping
// @Contact email@example.com
// @License Apache 2.0
// @LicenseUrl http://www.apache.org/licenses/LICENSE-2.0.html

package controllers

import (
	"fmt"
	"strconv"

	"github.com/astaxie/beego"
	"github.com/xuri/excelize/v2"
)

type TestController struct {
	beego.Controller
}

func (c *TestController) URLMapping() {
	c.Mapping("Ping", c.Ping)
	c.Mapping("UploadExcel", c.UploadExcel)
}

// Ping responde con "Pong" a la solicitud GET en /ping
// @router /pings [get]
func (c *TestController) Ping() {
	c.Ctx.WriteString("Pong")
}

// UploadExcel maneja la carga de un archivo Excel
// @router /upload [post]
func (c *TestController) UploadExcel() {

	file, _, err := c.GetFile("file")
	if err != nil {
		c.Ctx.ResponseWriter.WriteHeader(400)
		c.Ctx.WriteString("Error al obtener el archivo")

		return
	}
	defer file.Close()

	f, err := excelize.OpenReader(file)

	if err != nil {
		fmt.Println("Error al abrir el archivo")
		return
	}

	// rows, err := f.GetRows("Hoja1")
	// if err != nil {
	// 	fmt.Print("Error al obtenr hoja")
	// 	fmt.Print(err)
	// 	return
	// }

	// for _, row := range rows {

	// 	for _, col := range row {
	// 		fmt.Print(col, "\t")

	// 	}
	// }

	for i := 1; i < 100; i++ {
		cell, err := f.GetCellValue("Hoja1", "A"+strconv.Itoa(i))
		if err != nil {
			fmt.Println("Errro al imprimir  celdas ")
			return
		}
		fmt.Println(cell)
		if len(cell) < 1 {
			return
		}
	}

}
