package transferers

import (
	"github.com/astaxie/beego"
	"gitlab.com/24sessions/lib-go-logger/logger/services"
	"service-recordingStorage/models/data/sqs"
	"service-recordingStorage/services/storages"
)

const DefaultChunkSize = 1024 * 1024 // 1 Mb

type AmazonToDracoon struct {
	AbstractTransferer
}

// @Title set parameters
func (t *AmazonToDracoon) Init(message modelsDataSqs.FileTransfererMessage, lg *logger.Logger) error {
	if err := t.AbstractTransferer.Init(message, lg); err != nil {
		return err
	}

	t.self = t

	if err := t.initSource(message); err != nil {
		t.logger.Log(logger.CreateError("Could not initialize source storage: " + err.Error()).SetErrorCode(15009))
		return err
	}
	if err := t.initDestination(message); err != nil {
		t.logger.Log(logger.CreateError("Could not initialize destination storage: " + err.Error()).SetErrorCode(15009))
		return err
	}
	return nil
}

func (t *AmazonToDracoon) initSource(message modelsDataSqs.FileTransfererMessage) error {

	if t.source == nil {
		t.source = new(storages.S3)
		t.source.SetRemoveAfterTransferring(beego.AppConfig.DefaultBool("s3.removeAfterTransferring", false))
		t.source.SetChunkSize(beego.AppConfig.DefaultInt64("s3.chunkSize", DefaultChunkSize))
	}
	if err := t.source.Init(message, t.Logger()); err != nil {
		return err
	}
	if files, err := t.source.SourceFiles(message); err == nil {
		t.files = files
	} else {
		t.logger.Log(logger.CreateError("Could not get source files: " + err.Error()).SetErrorCode(15009))
		return err
	}
	if err := t.source.SetMeeting(message.Meeting); err != nil {
		t.logger.Log(logger.CreateError("Could not get meeting info: " + err.Error()).SetErrorCode(15009))
		return err
	}
	return nil
}

func (t *AmazonToDracoon) initDestination(message modelsDataSqs.FileTransfererMessage) error {

	if t.destination == nil {
		t.destination = new(storages.Dracoon)
	}
	if err := t.destination.Init(message, t.Logger()); err != nil {
		t.logger.Log(logger.CreateError("Could not initialize destination storage: " + err.Error()).SetErrorCode(15009))
		return err
	}
	if err := t.destination.Auth(); err != nil {
		t.logger.Log(logger.CreateError("Could not authenticate on destination storage: " + err.Error()).SetErrorCode(15009))
		return err
	}
	if folders, err := t.destination.DestinationFolders(message); err == nil {
		t.folders = folders
	} else {
		t.logger.Log(logger.CreateError("Could not get destination folders: " + err.Error()).SetErrorCode(15009))
		return err
	}
	if err := t.destination.SetMeeting(message.Meeting); err != nil {
		t.logger.Log(logger.CreateError("Could not get meeting info: " + err.Error()).SetErrorCode(15009))
		return err
	}
	return nil
}

func (t *AmazonToDracoon) TransfererName() string {
	return "AmazonToDracoon transferer"
}
