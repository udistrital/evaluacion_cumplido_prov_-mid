package controllers

import (
	"github.com/astaxie/beego"
	"github.com/udistrital/evaluacion_cumplido_prov_mid/helpers"
	"github.com/udistrital/utils_oas/errorhandler"
	"github.com/udistrital/utils_oas/requestresponse"
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

	defer errorhandler.HandlePanic(&c.Controller)

	file, _, err := c.GetFile("file")

	if err != nil {
		c.Ctx.Output.SetStatus(400)
		c.Data["json"] = requestresponse.APIResponseDTO(false, 400, nil, err.Error())
	}

	defer file.Close()

	response, intemsNoAgregados, err := helpers.CargaDataExcel(file)

	responseMap := map[string]interface{}{
		"itemsNoAgregados": intemsNoAgregados,
	}

	if err == nil {

		c.Ctx.Output.SetStatus(200)
		c.Data["json"] = requestresponse.APIResponseDTO(true, 200, responseMap, response)
	} else {
		c.Ctx.Output.SetStatus(400)
		c.Data["json"] = requestresponse.APIResponseDTO(false, 400, err)
	}
	c.ServeJSON()

}
