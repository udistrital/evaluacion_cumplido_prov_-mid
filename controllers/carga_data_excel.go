package controllers

import (
	"github.com/astaxie/beego"
	"github.com/udistrital/evaluacion_cumplido_prov_mid/services"
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

// UploadExcel handles the uploading of Excel files for evaluations.
// @Title UploadExcel
// @Description Upload an Excel file to process evaluation items by evaluacion_id
// @Param   file           formData file    true  "Excel file containing the evaluation items"
// @Param   idEvaluacion   formData int     true  "ID of the evaluation"
// @Success 200 {object} map[string]interface{} "Items processed successfully"
// @Failure 400 "Bad request: file or idEvaluacion is invalid"
// @Failure 500 "Internal server error"
// @router /upload [post]
func (c *CargaDataExcelController) UploadExcel() {

	defer errorhandler.HandlePanic(&c.Controller)

	file, _, err := c.GetFile("file")

	if err != nil {
		c.Ctx.Output.SetStatus(400)
		c.Data["json"] = requestresponse.APIResponseDTO(false, 400, nil, err.Error())
	}

	defer file.Close()

	idEvaluacion, err := c.GetInt("idEvaluacion")
	if err != nil {
		c.Ctx.Output.SetStatus(400)
		c.Data["json"] = requestresponse.APIResponseDTO(false, 400, nil, "El id de la evaluación es obligatorio")
		return
	}

	response, intemsNoAgregados, err := services.CargaDataExcel(file, idEvaluacion)

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
