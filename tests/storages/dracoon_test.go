package storages_test

import (
	"fmt"
	"io"
	"os"
	"service-recordingStorage/models/data/sqs"
	"service-recordingStorage/models/data/storage"
	"service-recordingStorage/models/data/storage/dracoon"
	"service-recordingStorage/services/storages"
	"service-recordingStorage/system"
	"service-recordingStorage/tests"
	"testing"
)

var dracoon *storages.Dracoon
var message modelsDataSqs.FileTransfererMessage

func init() {
	tests.Init()

	dracoon = new(storages.Dracoon)
	if err := dracoon.Init(message, tests.LogMain); err != nil {
		panic("Could not initialize Dracoon storage: " + err.Error())
	}
}

func TestAuthOK(t *testing.T) {
	if err := dracoon.Auth(); err != nil {
		t.Fatalf("Authentication failed, should be successful: %s", err.Error())
	}
}

func TestDestinationFoldersOK(t *testing.T) {
	message = modelsDataSqs.FileTransfererMessage{
		Storage: modelsDataSqs.Storage{
			Type: "dracoon",
			Extra: map[string]interface{}{
				"dracoon": map[string]interface{}{
					"targets": tests.TestDracoonTargets(),
				},
			},
		},
	}
	if _, err := dracoon.DestinationFolders(message); err != nil {
		t.Fatalf("Getting destination folders failed, should be successful: %s", err.Error())
	}
}

func TestCheckDestinationFailed(t *testing.T) {
	_ = dracoon.SetDestinationFolder(&modelsDataStorageDracoon.Target{
		RoomId: 5981635,
	})
	if err := dracoon.CheckDestinationFolder(); err == nil {
		t.Fatalf("Check destionation is successful, should be failed")
	}
}

func TestCheckDestinationOK(t *testing.T) {
	_ = dracoon.SetDestinationFolder(&modelsDataStorageDracoon.Target{
		RoomId: 6922227, // room for mike@24sessions.com account
	})
	if err := dracoon.CheckDestinationFolder(); err != nil {
		t.Fatalf("Check destionation failed, should be successful: %s", err.Error())
	}
}

func TestStartOK(t *testing.T) {
	dracoon.SetFileInfo(modelsDataStorage.FileInfo{
		Name:      "test-video",
		Extension: "mp4",
		Size:      123456,
	})
	_ = dracoon.SetMeeting(tests.TestMeeting())
	if err := dracoon.Start(); err != nil {
		t.Fatalf("Starting of transferring failed, should be successful: %s", err.Error())
	}
}

func TestUploadOK(t *testing.T) {
	fileName := "test-video.mp4"

	// open test file
	file, err := os.Open(fmt.Sprintf("%s/fixtures/%s", system.AppDir(), fileName))
	if err != nil {
		t.Fatalf("Can not open file %s from fixtures: %s", fileName, err.Error())
	}
	defer file.Close()

	// get file status
	fileStat, err := file.Stat()
	if err != nil {
		t.Fatalf("Can not get status for file %s from fixtures: %s", fileName, err.Error())
	}

	// prepare some data
	_ = dracoon.SetMeeting(tests.TestMeeting())
	dracoon.SetFileInfo(modelsDataStorage.FileInfo{
		Name:      "test-video",
		Extension: "mp4",
		Size:      int(fileStat.Size()),
	})
	downloadState := modelsDataStorage.State{
		Range: modelsDataStorage.LoadRange{
			Limit: fileStat.Size() - 1,
		},
	}

	// read chunk from file till EOF
	buffer := make([]byte, fileStat.Size()/3)
	for {
		bytes, err := file.Read(buffer)
		if err != nil {
			if err == io.EOF {
				break
			} else {
				t.Fatalf("Can not read file %s from fixtures: %s", fileName, err.Error())
			}
		}

		downloadState.Range.End = downloadState.Range.Start + int64(bytes) - 1

		// try to upload chunk
		if err := dracoon.Upload(buffer, downloadState); err != nil {
			t.Fatalf("Uploading failed, should be successful: %s", err.Error())
		}

		downloadState.Range.Start = downloadState.Range.End + 1
	}
}

func TestCloseDracoonOK(t *testing.T) {
	if err := dracoon.Close(); err != nil {
		t.Fatalf("Closing of storage failed, should be successful: %s", err.Error())
	}
}

func TestDracoonHealthCheckOK(t *testing.T) {
	if err := dracoon.HealthCheck(); err != nil {
		t.Fatalf("Health check failed, should be successful: %s", err.Error())
	}
}
