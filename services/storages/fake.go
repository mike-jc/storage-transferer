package storages

import (
	"gitlab.com/24sessions/lib-go-logger/logger/services"
	"service-recordingStorage/errors/storage"
	"service-recordingStorage/models/data/sqs"
	"service-recordingStorage/services/accessSharers"
	"service-recordingStorage/services/encryptors"
)

type Fake struct {
	AbstractStorage
}

func (s *Fake) StorageName() string {
	return "Fake storage"
}

func (s *Fake) Init(message modelsDataSqs.FileTransfererMessage, logger *logger.Logger) errorsStorage.ErrorContract {
	err := s.AbstractStorage.Init(message, logger)
	if err != nil {
		return err
	}

	s.self = s
	s.encryptor = new(encryptors.Fake)
	s.sharer = new(accessSharers.Fake)

	return err
}
