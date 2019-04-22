package transferers

import (
	"errors"
	"gitlab.com/24sessions/lib-go-logger/logger/services"
	"service-recordingStorage/models/data/sqs"
	"service-recordingStorage/services/storages"
)

type Fake struct {
	sourceAccessible bool

	AbstractTransferer
}

// @Title set parameters
func (t *Fake) Init(message modelsDataSqs.FileTransfererMessage, logger *logger.Logger) error {
	t.self = t
	t.logger = logger

	t.files = []interface{}{"fake file"}
	t.folders = []interface{}{"fake folder"}

	t.source = new(storages.Fake)
	t.destination = new(storages.Fake)

	// we can make source not accessible for tests
	t.sourceAccessible = true
	storageExtra, ok := message.Storage.Extra.(map[string]interface{})
	if ok {
		sourceAccessible, ok := storageExtra["sourceAccessible"]
		if ok && sourceAccessible.(bool) == false {
			t.sourceAccessible = false
		}
	}

	return t.AbstractTransferer.Init(message, logger)
}

func (t *Fake) TransfererName() string {
	return "Fake transferer"
}

// @Title check source
// @Description authenticate to source storage, define file path and check its existence
func (t *Fake) CheckSource() error {
	err := t.AbstractTransferer.CheckSource()
	if err != nil {
		return err
	}

	if !t.sourceAccessible {
		return errors.New("source is not accessible")
	}

	return err
}
