package controllers

import (
	"github.com/astaxie/beego"
	"gitlab.com/24sessions/lib-go-logger/logger/services"
	"service-recordingStorage/models/data/worker/options"
	"service-recordingStorage/services/rest/api"
	"service-recordingStorage/services/storages"
	"service-recordingStorage/workers"
)

type HealthCheckController struct {
	fileTransferer *workers.FileTransferer
	s3Storage      *storages.S3
	dracoonStorage *storages.Dracoon

	AbstractController
}

func (c *HealthCheckController) Prepare() {
	// FileTransferer
	c.fileTransferer = new(workers.FileTransferer)
	c.fileTransferer.SetLogger(LogMain)
	if err := c.fileTransferer.Init(modelsDataWorkerOptions.DefaultFileTransfererOptions()); err != nil {
		c.GetLogger().LogFatal(logger.CreateError("HealthCheck: can not initialize FileTransferer: " + err.Error()).SetErrorCode(15014))
		return
	}

	// AWS S3 storage
	c.s3Storage = new(storages.S3)
	c.s3Storage.SetLogger(LogMain)
	c.s3Storage.SetBucketName("24s-recordings")

	// Dracoon storage
	client := new(servicesRestApi.Dracoon)
	client.SetLogger(LogMain)
	if err := storages.InitDracoonRestClient(client, nil); err != nil { // test conf
		c.GetLogger().LogFatal(logger.CreateError("Can not initialize Dracoon API client: " + err.Error()).SetErrorCode(15014))
		return
	}
	c.dracoonStorage = new(storages.Dracoon)
	c.dracoonStorage.SetLogger(LogMain)
	c.dracoonStorage.SetClient(client)
}

// @Title Health check
// @Description check system
// @Success 200 {object} modelsData.Status
// @Failure 500 {object} modelsData.Error
// @router /healthcheck [get]
func (c *HealthCheckController) HealthCheck() {

	if beego.BConfig.RunMode == "prod" && workers.RunningStatus["FileTransferer"] < 1 {
		c.error("HealthCheck worker error: file_transferer", "worker_error_file_transferer", "No running workers")
		return
	}

	// FileTransferer's check
	if err := c.fileTransferer.HealthCheck(); err != nil {
		c.error("HealthCheck worker error: file_transferer", "worker_error_file_transferer", err.Error())
		return
	}

	// AWS S3 storage's check
	if err := c.s3Storage.HealthCheck(); err != nil {
		c.error("HealthCheck storage error: AWS S3", "storage_amazon_s3", err.Error())
		return
	}

	// Dracoon storage's check
	if err := c.dracoonStorage.HealthCheck(); err != nil {
		c.error("HealthCheck storage error: Dracoon", "storage_dracoon", err.Error())
		return
	}

	c.SuccessResponse()
}

func (c *HealthCheckController) error(message, shortMessage, reason string) {
	LogMain.Log(logger.CreateError(message).SetErrorCode(15014))
	c.ShowError(shortMessage, 500, 15014, reason, true)
}
