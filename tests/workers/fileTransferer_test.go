package workers_test

import (
	"encoding/json"
	"github.com/aws/aws-sdk-go/service/sqs"
	"service-recordingStorage/models/data/sqs"
	"service-recordingStorage/models/data/worker/options"
	"service-recordingStorage/tests"
	"service-recordingStorage/workers"
	"testing"
)

var ft *workers.FileTransferer
var sqsMessage *sqs.Message

func init() {
	tests.Init()

	ft = new(workers.FileTransferer)
	sqsMessage = new(sqs.Message)

	ft.SetLogger(tests.LogMain)
}

func TestInit(t *testing.T) {
	if err := ft.Init(modelsDataWorkerOptions.DefaultFileTransfererOptions()); err != nil {
		t.Fatalf("Can not initialize FileTransferer: %s", err.Error())
	}
}

func TestProcessSignalNotSqsMessage(t *testing.T) {
	// not SQS message
	message := struct {
		Body string
	}{"some body"}

	if ft.ProcessSignal(message) == true {
		t.Fatalf("Type casting OK for message that is not SQS message, should be error")
	}
}

func TestProcessSignalNotParsable(t *testing.T) {
	// incorrect JSON in SQS message body
	sqsMessage.SetBody("{\"field\":\"value\"")

	if ft.ProcessSignal(sqsMessage) == true {
		t.Fatalf("Parsing OK for non-parsable JSON, should be error")
	}
}

func TestProcessSignalTransfererNotFound(t *testing.T) {
	// unknown type of storage
	message := modelsDataSqs.FileTransfererMessage{
		Storage: modelsDataSqs.Storage{
			Type: "unknown",
		},
	}
	jsonMessage, _ := json.Marshal(message)
	sqsMessage.SetBody(string(jsonMessage))

	if ft.ProcessSignal(sqsMessage) == true {
		t.Fatalf("Transferer found, should be error for unknown type")
	}
}

func TestProcessSignalFakeTransfererOK(t *testing.T) {
	// unknown type of storage
	message := modelsDataSqs.FileTransfererMessage{
		Storage: modelsDataSqs.Storage{
			Type: "fake",
		},
	}
	jsonMessage, _ := json.Marshal(message)
	sqsMessage.SetBody(string(jsonMessage))

	if ft.ProcessSignal(sqsMessage) == false {
		t.Fatalf("Fake transferer not found, should exist")
	}
}

func TestProcessSignalFakeTransfererFailed(t *testing.T) {
	// unknown type of storage
	message := modelsDataSqs.FileTransfererMessage{
		Storage: modelsDataSqs.Storage{
			Type:  "fake",
			Extra: map[string]interface{}{"sourceAccessible": false},
		},
	}
	jsonMessage, _ := json.Marshal(message)
	sqsMessage.SetBody(string(jsonMessage))

	if ft.ProcessSignal(sqsMessage) == true {
		t.Fatalf("Fake transferer is initialized successfully, should be error")
	}
}
