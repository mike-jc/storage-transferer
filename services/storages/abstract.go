package storages

import (
	"gitlab.com/24sessions/lib-go-logger/logger/services"
	"service-recordingStorage/errors/storage"
	"service-recordingStorage/models/data/instance"
	"service-recordingStorage/models/data/sqs"
	"service-recordingStorage/models/data/storage"
	"service-recordingStorage/services/accessSharers"
	"service-recordingStorage/services/encryptors"
	"time"
)

type AbstractStorage struct {
	meeting modelsDataInstance.Meeting

	chunkSize               int64
	removeAfterTransferring bool
	state                   modelsDataStorage.State
	fileInfo                modelsDataStorage.FileInfo

	encryptor encryptors.AbstractEncryptorContract
	sharer    accessSharers.AbstractSharerContract

	logger *logger.Logger
	self   AbstractStorageContract

	AbstractStorageContract
}

type AbstractStorageContract interface {
	StorageName() string
	SetLogger(lg *logger.Logger)
	Logger() *logger.Logger
	Init(message modelsDataSqs.FileTransfererMessage, logger *logger.Logger) errorsStorage.ErrorContract
	SetEncryptor(encryptor encryptors.AbstractEncryptorContract)
	SetSharer(sharer accessSharers.AbstractSharerContract)
	Reset()
	Meeting() modelsDataInstance.Meeting
	SetMeeting(meeting modelsDataInstance.Meeting) errorsStorage.ErrorContract
	Auth() errorsStorage.ErrorContract
	SourceFiles(message modelsDataSqs.FileTransfererMessage) (files []interface{}, err errorsStorage.ErrorContract)
	SetSourceFile(file interface{}) errorsStorage.ErrorContract
	CheckSourceFile() errorsStorage.ErrorContract
	DestinationFolders(message modelsDataSqs.FileTransfererMessage) (folders []interface{}, err errorsStorage.ErrorContract)
	SetDestinationFolder(folder interface{}) errorsStorage.ErrorContract
	CheckDestinationFolder() errorsStorage.ErrorContract
	FileInfo() modelsDataStorage.FileInfo
	SetFileInfo(info modelsDataStorage.FileInfo)
	SetRemoveAfterTransferring(remove bool)
	SetChunkSize(size int64)
	ChunkSize() int64
	State() modelsDataStorage.State
	Start() errorsStorage.ErrorContract
	Download() (data []byte, err errorsStorage.ErrorContract)
	Upload(data []byte, downloadState modelsDataStorage.State) errorsStorage.ErrorContract
	Close() errorsStorage.ErrorContract
	Revert()
	HealthCheck() errorsStorage.ErrorContract
}

func (s *AbstractStorage) StorageName() string {
	return "Abstract storage"
}

func (s *AbstractStorage) SetLogger(lg *logger.Logger) {
	s.logger = lg
}

func (s *AbstractStorage) Logger() *logger.Logger {
	return s.logger
}

func (s *AbstractStorage) Init(message modelsDataSqs.FileTransfererMessage, logger *logger.Logger) errorsStorage.ErrorContract {
	s.self = s
	s.logger = logger
	s.meeting = modelsDataInstance.Meeting{}
	return nil
}

func (s *AbstractStorage) SetEncryptor(encryptor encryptors.AbstractEncryptorContract) {
	s.encryptor = encryptor
}

func (s *AbstractStorage) SetSharer(sharer accessSharers.AbstractSharerContract) {
	s.sharer = sharer
}

// @Description need to be done before every downloading/uploading
func (s *AbstractStorage) Reset() {
	if s.encryptor != nil {
		s.encryptor.Reset()
	}
	if s.sharer != nil {
		s.sharer.Reset()
	}

	s.fileInfo = modelsDataStorage.FileInfo{}
	s.state = modelsDataStorage.State{}
}

func (s *AbstractStorage) Meeting() modelsDataInstance.Meeting {
	return s.meeting
}

func (s *AbstractStorage) SetMeeting(meeting modelsDataInstance.Meeting) errorsStorage.ErrorContract {
	s.meeting = meeting
	if meeting.Id == 0 {
		return errorsStorage.NewError("wrong meeting ID passed to the storage", errorsStorage.Message, errorsStorage.GeneralError)
	}
	if meeting.DateUTC.Unix() <= 0 {
		return errorsStorage.NewError("wrong meeting date passed to the storage", errorsStorage.Message, errorsStorage.GeneralError)
	}
	if _, err := time.LoadLocation(meeting.Timezone); err != nil {
		return errorsStorage.NewError("wrong meeting timezone passed to the storage", errorsStorage.Message, errorsStorage.GeneralError)
	}
	if meeting.User.Email == "" {
		return errorsStorage.NewError("wrong meeting user's email passed to the storage", errorsStorage.Message, errorsStorage.GeneralError)
	}
	return nil
}

// @Description authentication on the storage side
func (s *AbstractStorage) Auth() errorsStorage.ErrorContract {
	return nil
}

func (s *AbstractStorage) SourceFiles(message modelsDataSqs.FileTransfererMessage) (files []interface{}, err errorsStorage.ErrorContract) {
	return
}

func (s *AbstractStorage) SetSourceFile(file interface{}) errorsStorage.ErrorContract {
	return nil
}

// @Description check if source file exists
// (if we're gonna download from the storage)
// and set file info
func (s *AbstractStorage) CheckSourceFile() errorsStorage.ErrorContract {
	return nil
}

func (s *AbstractStorage) DestinationFolders(message modelsDataSqs.FileTransfererMessage) (folders []interface{}, err errorsStorage.ErrorContract) {
	return
}

func (s *AbstractStorage) SetDestinationFolder(folder interface{}) errorsStorage.ErrorContract {
	return nil
}

// @Description: check if there's free room, if quota isn't used up etc.
// (if we're gonna upload to the storage)
func (s *AbstractStorage) CheckDestinationFolder() errorsStorage.ErrorContract {
	return nil
}

func (s *AbstractStorage) FileInfo() modelsDataStorage.FileInfo {
	return s.fileInfo
}

func (s *AbstractStorage) SetFileInfo(info modelsDataStorage.FileInfo) {
	s.fileInfo = info
}

// @Description do some actions before starting transferring (e.g., open upload channel)
func (s *AbstractStorage) Start() errorsStorage.ErrorContract {
	return nil
}

// @Description set if we need to remove file after transferring
func (s *AbstractStorage) SetRemoveAfterTransferring(remove bool) {
	s.removeAfterTransferring = remove
}

// @Description set chunk size
func (s *AbstractStorage) SetChunkSize(size int64) {
	s.chunkSize = size
}

func (s *AbstractStorage) ChunkSize() int64 {
	return s.chunkSize
}

func (s *AbstractStorage) State() modelsDataStorage.State {
	return s.state
}

// @Description download data from the storage (probably chunk by chunk)
func (s *AbstractStorage) Download() (data []byte, err errorsStorage.ErrorContract) {
	return
}

// @Description upload data to the storage (probably chunk by chunk)
func (s *AbstractStorage) Upload(data []byte, downloadState modelsDataStorage.State) errorsStorage.ErrorContract {
	return nil
}

// @Description do some actions after completing downloading/uploading
// (e.g., close uploading channel, remove file from Amazon S3, sharing access for Dracoon users, etc)
func (s *AbstractStorage) Close() errorsStorage.ErrorContract {
	return nil
}

// @Description do some actions after failed downloading/uploading
func (s *AbstractStorage) Revert() {

}

// @Description check if service works fine
func (s *AbstractStorage) HealthCheck() errorsStorage.ErrorContract {
	return nil
}

func (s *AbstractStorage) logStatus(msg string) {
	s.self.Logger().Log(logger.CreateInfo(s.self.StorageName() + ": " + msg).SetDescriptionCode(15012))
}
