package transferers_test

import (
	"service-recordingStorage/models/data/sqs"
	"service-recordingStorage/services/storages"
	"service-recordingStorage/services/transferers"
	"service-recordingStorage/tests"
	"testing"
)

var dracoon *transferers.AmazonToDracoon
var message *modelsDataSqs.FileTransfererMessage

func init() {
	tests.Init()

	dracoon = new(transferers.AmazonToDracoon)
	dracoon.SetSource(new(storages.Fake))
	dracoon.SetDestination(new(storages.Fake))

	message = new(modelsDataSqs.FileTransfererMessage)
	message.Meeting = tests.TestMeeting()
	message.Storage.Type = "dracoon"
}

func TestInitOK(t *testing.T) {
	message.Storage.Extra = map[string]interface{}{"parentRoomId": 1234567890}

	if err := dracoon.Init(*message, tests.LogMain); err != nil {
		t.Fatalf("%s initialization failed, should be successful: %s", dracoon.TransfererName(), err.Error())
	}
}
