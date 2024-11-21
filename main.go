package main

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/plugins/cors"
	//apistatus "github.com/udistrital/utils_oas/apiStatusLib"
	//auditoria "github.com/udistrital/utils_oas/auditoria"
	_ "github.com/udistrital/evaluacion_cumplido_prov_mid/routers"
	//"github.com/udistrital/utils_oas/customerrorv2"
)

func main() {
	AllowedOrigins := []string{"*", "*.udistrital.edu.co"}
	if beego.BConfig.RunMode == beego.AppConfig.String("runmode") {
		beego.BConfig.WebConfig.DirectoryIndex = true
		beego.BConfig.WebConfig.StaticDir["/swagger"] = "swagger"
	}
	beego.InsertFilter("*", beego.BeforeRouter, cors.Allow(&cors.Options{
		AllowOrigins: AllowedOrigins,
		AllowMethods: []string{"PUT", "PATCH", "GET", "POST", "OPTIONS", "DELETE"},
		AllowHeaders: []string{"Origin", "x-requested-with",
			"content-type",
			"accept",
			"origin",
			"authorization",
			"x-csrftoken"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))
	//beego.ErrorController(&customerrorv2.CustomErrorController{})
	//apistatus.Init()
	//	auditoria.InitMiddleware()
	beego.Run()
}
