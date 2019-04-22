package workers_test

import (
	"os"
	"os/exec"
	"service-recordingStorage/models/data/worker/options"
	"service-recordingStorage/tests"
	"service-recordingStorage/workers"
	"testing"
)

var qr *workers.QueueReader

func init() {
	tests.Init()

	qr = new(workers.QueueReader)
	qr.SetLogger(tests.LogMain)
}

func TestInitQueueNotFound(t *testing.T) {

	// need this approach since log.Fatal exits with code 1 when queue doesn't exist
	if os.Getenv("NON_EXISTING_QUEUE") == "1" {
		qr.SetQueueName("non-existing-queue")
		if err := qr.Init(modelsDataWorkerOptions.DefaultQueueReaderOptions()); err != nil {
			t.Fatalf("Can not initialize QueueReader: %s", err.Error())
		}
	} else {
		cmd := exec.Command(os.Args[0], "-test.run=TestInitQueueNotFound")
		cmd.Env = append(os.Environ(), "NON_EXISTING_QUEUE=1")

		err := cmd.Run()
		if e, ok := err.(*exec.ExitError); ok && !e.Success() {
			return
		}
		t.Fatalf("For not existing queue QueueReader is initialized successfully, should be error")
	}
}

func TestInitQueueOK(t *testing.T) {
	qr.SetQueueName("recordings-transferring")
	if err := qr.Init(modelsDataWorkerOptions.DefaultQueueReaderOptions()); err != nil {
		t.Fatalf("Can not initialize QueueReader: %s", err.Error())
	}
}

func TestHealthCheckOK(t *testing.T) {
	qr.SetQueueName("recordings-transferring")
	if err := qr.HealthCheck(); err != nil {
		t.Fatalf("Health check failed, should be successful: %s", err.Error())
	}
}
