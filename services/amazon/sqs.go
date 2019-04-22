package amazon

import (
	"github.com/astaxie/beego/config"
	_ "github.com/astaxie/beego/config/yaml"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/sqs"
	"gitlab.com/24sessions/lib-go-logger/logger/services"
	"service-recordingStorage/services"
	"service-recordingStorage/system"
)

type Sqs struct {
	engine       *sqs.SQS
	configTopics config.Configer

	AbstractService
}

func (o *Sqs) SetConfigTopics(config config.Configer) {
	o.configTopics = config
}

func (o *Sqs) ConfigTopics() config.Configer {
	if o.configTopics == nil {
		cnf, err := config.NewConfig("yaml", system.AppDir()+"/conf/pubsub/sqs/topics.yml")
		if err != nil {
			services.LogMain.LogFatal(logger.CreateError("Can not read SQS configuration from file: " + err.Error()).SetErrorCode(15005))
			return nil
		}
		o.configTopics = cnf
	}
	return o.configTopics
}

func (o *Sqs) Engine() *sqs.SQS {
	if o.engine == nil {
		o.engine = sqs.New(o.Session(), aws.NewConfig() /*.WithLogLevel(aws.LogDebugWithHTTPBody)*/)
	}
	return o.engine
}
