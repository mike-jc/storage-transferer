package main

import (
	"github.com/astaxie/beego"
	"service-recordingStorage/controllers"
	"service-recordingStorage/resources"
	"service-recordingStorage/system"
	"service-recordingStorage/workers"
)

func init() {
	system.SetAppDirToCurrentDir(0)

	// logger
	lg := resources.InitLogger(resources.AppTypeRegular)
	controllers.LogMain = lg

	// routing
	controllers.RegisterRoutes()

	// workers
	workers.StartWorkers(lg)
}

func main() {
	if beego.BConfig.RunMode == "dev" {
		beego.BConfig.WebConfig.DirectoryIndex = true
		beego.BConfig.WebConfig.StaticDir["/swagger"] = "swagger"
	}

	beego.Run()
}
