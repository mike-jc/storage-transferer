package transferers

import (
	"errors"
	"gitlab.com/24sessions/lib-go-logger/logger/services"
	"service-recordingStorage/models/data/sqs"
	"service-recordingStorage/services"
	"service-recordingStorage/services/storages"
)

type AbstractTransferer struct {
	source      storages.AbstractStorageContract
	destination storages.AbstractStorageContract

	files   []interface{}
	folders []interface{}

	logger *logger.Logger
	self   AbstractTransfererContract

	AbstractTransfererContract
}

type AbstractTransfererContract interface {
	Init(message modelsDataSqs.FileTransfererMessage, logger *logger.Logger) error
	Reset()
	TransfererName() string
	Logger() *logger.Logger
	SetSource(source storages.AbstractStorageContract)
	SetDestination(destination storages.AbstractStorageContract)
	Files() []interface{}
	Folders() []interface{}
	Start() error
	SetSourceFile(file interface{}) error
	CheckSource() error
	SetDestinationFolder(folder interface{}) error
	CheckDestination() error
	DownloadData() (data []byte, err error)
	UploadData(data []byte) error
	CloseSource() error
	CloseDestination() error
	Revert()
}

// @Title set parameters
func (t *AbstractTransferer) Init(message modelsDataSqs.FileTransfererMessage, logger *logger.Logger) error {
	t.self = t
	t.logger = logger
	t.logger.SetInstance(services.MessageDomain(message))
	t.logStatus("initialization")
	return nil
}

func (t *AbstractTransferer) Reset() {
	if t.logger != nil {
		t.logger.SetInstance("").SetObject("", "")
	}
}

func (t *AbstractTransferer) TransfererName() string {
	return "Abstract transferer"
}

func (t *AbstractTransferer) Logger() *logger.Logger {
	return t.logger
}

func (t *AmazonToDracoon) SetSource(source storages.AbstractStorageContract) {
	t.source = source
}

func (t *AmazonToDracoon) SetDestination(destination storages.AbstractStorageContract) {
	t.destination = destination
}

func (t *AbstractTransferer) Files() []interface{} {
	return t.files
}

func (t *AbstractTransferer) Folders() []interface{} {
	return t.folders
}

// @Description do some actions before starting transferring (e.g., open upload channel)
func (t *AbstractTransferer) Start() error {
	t.logStatus("start transferring")

	if err := t.source.Start(); err != nil {
		return err
	}
	if err := t.destination.Start(); err != nil {
		return err
	}
	return nil
}

func (t *AbstractTransferer) SetSourceFile(file interface{}) error {
	if t.source == nil {
		return errors.New("Source is not set for " + t.self.TransfererName())
	}

	t.source.Reset()
	return t.source.SetSourceFile(file)
}

// @Title check source
// @Description authenticate to source storage, define file path and check its existence
func (t *AbstractTransferer) CheckSource() error {
	t.logStatus("source checking")

	if err := t.source.CheckSourceFile(); err != nil {
		return err
	}

	// after checking source file we have file info and pass it to the destination storage
	t.destination.SetFileInfo(t.source.FileInfo())
	return nil
}

func (t *AbstractTransferer) SetDestinationFolder(folder interface{}) error {
	if t.destination == nil {
		return errors.New("Destination is not set for " + t.self.TransfererName())
	}

	t.destination.Reset()
	return t.destination.SetDestinationFolder(folder)
}

// @Title check destination
// @Description authenticate to destination storage and check quotas, room etc.
func (t *AbstractTransferer) CheckDestination() error {
	t.logStatus("destination checking")

	return t.destination.CheckDestinationFolder()
}

// @Title download data
// @Description download file chunk from Amazon S3
// If data and err are both nil it means that there no more data and transferring is completed
func (t *AbstractTransferer) DownloadData() (data []byte, err error) {
	t.logStatus("data downloading")

	return t.source.Download()
}

// @Title upload data
// @Description upload data to Dracoon via upload channel;
// encrypt the data before uploading if necessary
// @Param data []byte data
func (t *AbstractTransferer) UploadData(data []byte) error {
	t.logStatus("data uploading")

	return t.destination.Upload(data, t.source.State())
}

// @Description do some actions after completing transferring of the source file
// to all upload targets (e.g., remove file from Amazon S3)
func (t *AbstractTransferer) CloseSource() error {
	t.logStatus("close source")

	if err := t.source.Close(); err != nil {
		return err
	}
	return nil
}

// @Description do some actions after completing uploading
// (e.g., close uploading channel, share access etc.)
func (t *AbstractTransferer) CloseDestination() error {
	t.logStatus("close destination")

	if err := t.destination.Close(); err != nil {
		return err
	}
	return nil
}

// @Description do some actions after completing transferring
// (e.g., close uploading channel, remove file from Amazon S3)
func (t *AbstractTransferer) Revert() {
	t.logStatus("revert transferring")

	t.source.Revert()
	t.destination.Revert()
}

func (t *AbstractTransferer) logStatus(msg string) {
	t.Logger().Log(logger.CreateInfo(t.self.TransfererName() + ": " + msg).SetDescriptionCode(15006))
}
