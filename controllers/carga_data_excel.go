package controllers

import (
	"github.com/astaxie/beego"
	"github.com/udistrital/evaluacion_cumplido_prov_mid/helpers"
)

// Carga-Data-ExcelController operations for Carga-Data-Excel
type CargaDataExcelController struct {
	beego.Controller
}

// URLMapping ...
func (c *CargaDataExcelController) URLMapping() {
	c.Mapping("UploadExcel", c.UploadExcel)
}

// @router /upload [post]
func (c *CargaDataExcelController) UploadExcel() {

	file, _, err := c.GetFile("file")
	if err != nil {
		c.Ctx.ResponseWriter.WriteHeader(400)
		c.Ctx.WriteString("Error al obtener el archivo")

		return
	}

	defer file.Close()

	helpers.CargaDataExcel(file)
}
