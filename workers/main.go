package workers

import (
	"github.com/astaxie/beego"
	"gitlab.com/24sessions/lib-go-logger/logger/services"
	"service-recordingStorage/models/data/worker/options"
)

var RunningStatus map[string]int

func StartWorkers(log *logger.Logger) {

	// Get workers' configuration
	countFileTransferer, err := beego.AppConfig.Int("workers.FileTransferer.count")
	if err != nil {
		log.LogFatal(logger.CreateAlert("workers.FileTransferer.count is not set in config: " + err.Error()).SetErrorCode(15002))
	}
	if beego.BConfig.RunMode == "prod" && countFileTransferer < 1 {
		log.LogFatal(logger.CreateAlert("workers.FileTransferer.count can not be less than 1").SetErrorCode(15002))
	}

	countDracoonGlobalSharer, err := beego.AppConfig.Int("workers.DracoonGlobalSharer.count")
	if err != nil {
		log.LogFatal(logger.CreateAlert("workers.DracoonGlobalSharer.count is not set in config: " + err.Error()).SetErrorCode(15002))
	}
	if countDracoonGlobalSharer > 1 {
		log.LogFatal(logger.CreateAlert("workers.DracoonGlobalSharer.count can not be more than 1").SetErrorCode(15002))
	}
	delayDracoonGlobalSharer, err := beego.AppConfig.Int("workers.DracoonGlobalSharer.jobDelayMinutes")
	if err != nil {
		log.LogFatal(logger.CreateAlert("workers.DracoonGlobalSharer.jobDelayMinutes is not set in config: " + err.Error()).SetErrorCode(15002))
	}
	if delayDracoonGlobalSharer < 1 {
		log.LogFatal(logger.CreateAlert("workers.DracoonGlobalSharer.jobDelayMinutes can not be less than 1").SetErrorCode(15002))
	}

	// Make workers' options
	ftOptions := modelsDataWorkerOptions.FileTransfererOptions{
		QueueReaderOptions: modelsDataWorkerOptions.QueueReaderOptions{
			VisibilityTimeout:   12 * 60 * 60, // in seconds; 12 hours
			WaitTimeSeconds:     3,            // in seconds
			MaxNumberOfMessages: 1,
			RequestDelay:        10, // in seconds
		},
	}
	dsOptions := modelsDataWorkerOptions.DracoonGlobalSharerOptions{
		InstanceJobRunnerOptions: modelsDataWorkerOptions.InstanceJobRunnerOptions{
			JobDelay: delayDracoonGlobalSharer,
		},
	}

	// Start workers
	RunningStatus = make(map[string]int)

	RunningStatus["FileTransferer"] = 0
	for i := 0; i < countFileTransferer; i++ {
		transferer := new(FileTransferer)
		if err := transferer.Init(ftOptions); err != nil {
			log.LogFatal(logger.CreateAlert("Can not init worker FileTransferer: " + err.Error()).SetErrorCode(15002))
		}

		RunningStatus["FileTransferer"]++
		go func() {
			defer func() { RunningStatus["FileTransferer"]-- }()
			transferer.Run()
		}()
	}

	RunningStatus["DracoonGlobalSharer"] = 0
	for i := 0; i < countDracoonGlobalSharer; i++ {
		sharer := new(DracoonGlobalSharer)
		if err := sharer.Init(dsOptions); err != nil {
			log.LogFatal(logger.CreateAlert("Can not init worker DracoonGlobalSharer: " + err.Error()).SetErrorCode(15002))
		}

		RunningStatus["DracoonGlobalSharer"]++
		go func() {
			defer func() { RunningStatus["DracoonGlobalSharer"]++ }()
			sharer.Run()
		}()
	}
}
