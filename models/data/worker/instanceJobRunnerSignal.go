package modelsDataWorker

import "gitlab.com/24sessions/sdk-go-configurator/data"

type InstanceJobRunnerSignal struct {
	Domain string
	Config *sdksData.InstanceContainer
}
