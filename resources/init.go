package resources

import (
	_ "gitlab.com/24sessions/lib-beego-logger/logger"
	logModels "gitlab.com/24sessions/lib-go-logger/logger/models"
	"gitlab.com/24sessions/lib-go-logger/logger/services"
	"service-recordingStorage/services"
)

const AppTypeRegular = 1
const AppTypeCli = 2
const AppTypeTest = 3

var LogMain *logger.Logger

func InitLogger(appType int) *logger.Logger {
	var lg = new(logger.Logger).SetSubject("application", "")
	var logRow *logModels.ContLog

	switch appType {
	case AppTypeRegular:
		logRow = logger.CreateInfo("Starting application").SetDescriptionCode(15000)
	case AppTypeCli:
		logRow = logger.CreateInfo("Starting CLI application").SetDescriptionCode(15001)
	case AppTypeTest:
		logRow = logger.CreateInfo("Starting test application")
	}

	if !lg.LogForce(logRow) {
		panic("Can not init Logger")
	}

	LogMain = lg
	services.LogMain = lg
	return lg
}
