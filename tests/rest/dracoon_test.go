package rest_test

import (
	"service-recordingStorage/models/data/storage/dracoon"
	"service-recordingStorage/services"
	"service-recordingStorage/services/rest/api"
	"service-recordingStorage/services/storages"
	"service-recordingStorage/tests"
	"testing"
)

var dracoon *servicesRestApi.Dracoon

func init() {
	tests.Init()

	dracoon = new(servicesRestApi.Dracoon)
	dracoon.SetLogger(tests.LogMain)
	if err := storages.InitDracoonRestClient(dracoon, nil); err != nil {
		panic(err.Error())
	}
}

func TestTokenValueOK(t *testing.T) {
	if _, err := dracoon.TokenValue(); err != nil {
		t.Fatalf("Getting auth token failed, should be successful: %s", err.Error())
	}
}

func TestChildRoomOK(t *testing.T) {
	if roomId, err := services.DracoonParentRoomId(services.MessageDomain(nil)); err != nil {
		t.Fatalf("Wrong dracoon parnt room id in app.conf")
	} else if _, err := dracoon.ChildRooms(roomId); err != nil {
		t.Fatalf("Getting child rooms failed, should be successful: %s", err.Error())
	}
}

func TestNodeOK(t *testing.T) {
	if roomId, err := services.DracoonParentRoomId(services.MessageDomain(nil)); err != nil {
		t.Fatalf("Wrong dracoon parnt room id in app.conf")
	} else if _, err := dracoon.Node(roomId); err != nil {
		t.Fatalf("Getting node info failed, should be successful: %s", err.Error())
	}
}

func TestFindOrCreateRoomOK(t *testing.T) {
	searchingPath := "Instance_1/Some new room"
	if roomId, err := services.DracoonParentRoomId(services.MessageDomain(nil)); err != nil {
		t.Fatalf("Wrong dracoon parent room id in app.conf")
	} else if _, err := dracoon.FindOrCreateRoom(roomId, searchingPath); err != nil {
		t.Fatalf("Finding room failed, should be successful: %s", err.Error())
	}
}

func TestCreateUploadChannelOK(t *testing.T) {
	if roomId, err := services.DracoonParentRoomId(services.MessageDomain(nil)); err != nil {
		t.Fatalf("Wrong dracoon parent room id in app.conf")
	} else {
		target := modelsDataStorageDracoon.Target{
			RoomId: roomId,
		}
		if _, err := dracoon.CreateUploadChannel(target, "test-file.txt", 12345, ""); err != nil {
			t.Fatalf("Creating of upload channel failed for target %+v, should be successful: %s", target, err.Error())
		}
	}
}

func TestUserKeyPairOK(t *testing.T) {
	if _, err := dracoon.UserKeyPair(); err != nil {
		t.Fatalf("Getting of user key pair failed, should be successful: %s", err.Error())
	}
}

func TestMissingFileKeysOK(t *testing.T) {
	if roomId, err := services.DracoonParentRoomId(services.MessageDomain(nil)); err != nil {
		t.Fatalf("Wrong dracoon parent room id in app.conf")
	} else if _, err := dracoon.MissingFileKeys(roomId, 0); err != nil {
		t.Fatalf("Getting of missing file keys failed, should be successful: %s", err.Error())
	}
}

func TestPingOK(t *testing.T) {
	if err := dracoon.Ping(); err != nil {
		t.Fatalf("Ping failed, should be successful: %s", err.Error())
	}
}
