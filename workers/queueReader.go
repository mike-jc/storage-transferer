package workers

import (
	"errors"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/sqs"
	"gitlab.com/24sessions/lib-go-logger/logger/services"
	"math/rand"
	"service-recordingStorage/models/data/worker/options"
	"service-recordingStorage/services/amazon"
	"strings"
	"time"
)

type QueueReader struct {
	queueName string
	queueUrl  string
	options   modelsDataWorkerOptions.QueueReaderOptions

	amazon.Sqs
	AbstractWorker
}

func (o *QueueReader) WorkerName() string {
	return "Abstract QueueReader"
}

func (o *QueueReader) Init(options interface{}) error {
	if err := o.AbstractWorker.Init(options); err != nil {
		return err
	}

	o.self = o

	o.queueUrl = o.ConfigTopics().String(o.queueName)
	if o.queueUrl == "" {
		return errors.New(fmt.Sprintf("Topic %s is absent in PubSub SQS config", o.queueName))
	}

	if qrOptions, ok := options.(modelsDataWorkerOptions.QueueReaderOptions); !ok {
		return errors.New(fmt.Sprintf("Wrong options format for QueueReader. Should be modelsDataWorkerOptions.QueueReaderOptions, but got <%T> %+v", options, options))
	} else {
		o.options = qrOptions
	}

	return nil
}

func (o *QueueReader) SetQueueName(queueName string) {
	o.queueName = queueName
}

func (o *QueueReader) Run() {
	o.AbstractWorker.Run()

	for {
		output, err := o.Engine().ReceiveMessage(&sqs.ReceiveMessageInput{
			QueueUrl:            aws.String(o.queueUrl),
			VisibilityTimeout:   aws.Int64(o.options.VisibilityTimeout),
			WaitTimeSeconds:     aws.Int64(o.options.WaitTimeSeconds),
			MaxNumberOfMessages: aws.Int64(o.options.MaxNumberOfMessages),
		})
		if err != nil {
			o.Logger().Log(logger.CreateError(err.Error()).SetErrorCode(15005))
			time.Sleep(time.Duration(o.options.WaitTimeSeconds) * time.Second)
		} else {
			for _, message := range output.Messages {
				if o.self.ProcessSignal(message) {
					o.ackMessage(message)
				}
			}
		}
		time.Sleep(time.Duration(o.options.RequestDelay)*time.Second + time.Duration(rand.Intn(1000))*time.Millisecond)
	}
}

func (o *QueueReader) ackMessage(message *sqs.Message) {
	_, err := o.Engine().DeleteMessage(&sqs.DeleteMessageInput{
		QueueUrl:      aws.String(o.queueUrl),
		ReceiptHandle: message.ReceiptHandle,
	})
	if err != nil {
		o.LogError("Can not delete sqs message for "+o.queueName, 15005)
	}
}

func (o *QueueReader) HealthCheck() error {
	if err := o.AbstractWorker.HealthCheck(); err != nil {
		return err
	}

	urlParts := strings.Split(o.queueUrl, "/")

	// check if queue exists and we have access to it
	_, err := o.Engine().GetQueueUrl(&sqs.GetQueueUrlInput{
		QueueName: aws.String(urlParts[len(urlParts)-1]),
	})
	return err
}
