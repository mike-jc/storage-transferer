package controllers

import (
	"github.com/astaxie/beego"
	"gitlab.com/24sessions/lib-go-logger/logger/services"
)

var LogMain *logger.Logger

func RegisterRoutes() {
	healthCheckController := new(HealthCheckController)
	healthCheckController.Prepare()
	beego.Router("/healthcheck", healthCheckController, "get:HealthCheck")
	beego.Router("/api/v1/healthcheck", healthCheckController, "get:HealthCheck")
}
