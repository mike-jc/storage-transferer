package workers

import (
	"errors"
	"fmt"
	"gitlab.com/24sessions/sdk-go-configurator"
	"service-recordingStorage/models/data/worker"
	"service-recordingStorage/models/data/worker/options"
	"time"
)

type InstanceJobRunner struct {
	options      modelsDataWorkerOptions.InstanceJobRunnerOptions
	configurator *sdks.Configurator

	AbstractWorker
}

func (o *InstanceJobRunner) WorkerName() string {
	return "InstanceJobRunner"
}

func (o *InstanceJobRunner) Init(options interface{}) error {
	if err := o.AbstractWorker.Init(options); err != nil {
		return err
	}

	if iOptions, ok := options.(modelsDataWorkerOptions.InstanceJobRunnerOptions); !ok {
		return errors.New(fmt.Sprintf("Wrong options format for InstanceJobRunner. Should be modelsDataWorkerOptions.InstanceJobRunnerOptions, but got <%T> %+v", options, options))
	} else {
		if iOptions.JobDelay < 1 {
			return errors.New("modelsDataWorkerOptions.InstanceJobRunnerOptions::JobDelay can not be less than 1")
		}
		o.options = iOptions
	}

	o.self = o

	o.configurator = new(sdks.Configurator)
	return nil
}

func (o *InstanceJobRunner) Run() {
	o.AbstractWorker.Run()

	for {
		if domains, err := o.configurator.GetListOfInstances(); err != nil {
			o.LogError(fmt.Sprintf("Can not get list of instances: %s", err.Error()), 15013)
		} else {
			for _, domain := range domains {
				if len(domain) == 0 {
					continue
				}
				if config, err := o.configurator.GetInstance(domain); err != nil || config == nil {
					o.LogError(fmt.Sprintf("Can not get instances configuration for %s: %s", domain, err.Error()), 15013)
					continue
				} else {
					o.self.ProcessSignal(modelsDataWorker.InstanceJobRunnerSignal{
						Domain: domain,
						Config: config,
					})
				}
			}
		}
		time.Sleep(time.Duration(o.options.JobDelay) * time.Minute)
	}
}
