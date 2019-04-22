package workers

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/astaxie/beego"
	"github.com/aws/aws-sdk-go/service/sqs"
	"gitlab.com/24sessions/lib-go-logger/logger/services"
	"service-recordingStorage/errors/storage"
	"service-recordingStorage/models/data/sqs"
	"service-recordingStorage/models/data/universe"
	"service-recordingStorage/models/data/worker/options"
	"service-recordingStorage/services/transferers"
	"strconv"
	"time"
)

type FileTransferer struct {
	options     modelsDataWorkerOptions.FileTransfererOptions
	factory     *transferers.FileTransfererFactory
	transferers map[string]transferers.AbstractTransfererContract

	QueueReader
}

func (o *FileTransferer) Init(options interface{}) error {

	qName := beego.AppConfig.String("workers.FileTransferer.queue")
	if qName == "" {
		return errors.New("Incorrect queue name in app for workers.FileTransferer.queue")
	}

	if ftOptions, ok := options.(modelsDataWorkerOptions.FileTransfererOptions); !ok {
		return errors.New(fmt.Sprintf("Wrong options format for FileTransferer. Should be modelsDataWorkerOptions.FileTransfererOptions but got <%T> %+v", options, options))
	} else {
		o.options = ftOptions
	}

	o.QueueReader.SetQueueName(qName)
	if err := o.QueueReader.Init(o.options.QueueReaderOptions); err != nil {
		return err
	}

	o.self = o

	o.factory = new(transferers.FileTransfererFactory)

	o.reportLevel, _ = beego.AppConfig.Int("workers.FileTransferer.reportLevel")
	if o.reportLevel <= 0 {
		o.reportLevel = modelsDataUniverse.ReportLevelAll
	}

	return nil
}

func (o *FileTransferer) WorkerName() string {
	return "FileTransferer"
}

// @Description Process signal (queue message)
// @Result false - ignore signal; true - mark signal as processed
func (o *FileTransferer) ProcessSignal(signal interface{}) bool {
	if !o.QueueReader.ProcessSignal(signal) {
		return false
	}

	timeStart := time.Now()

	sqsMessage, ok := signal.(*sqs.Message)
	if !ok {
		o.LogError("Got signal that is not queue message", 15007)
		return false
	}

	var message modelsDataSqs.FileTransfererMessage
	err := json.Unmarshal([]byte(*sqsMessage.Body), &message)

	if err != nil {
		o.LogError("Can not parse queue message from "+o.queueName+": "+err.Error(), 15007)
		return false
	} else {
		o.logger.SetObject("meeting", strconv.FormatInt(message.Meeting.Id, 10))

		if transferer, err := o.transfererByType(message.Storage.Type); err != nil {
			o.LogError(err.Error(), 15007)
			return false
		} else {
			o.Logger().Log(logger.CreateInfo(fmt.Sprintf("%s [%s]: signal processing has started", o.WorkerName(), o.Logger().GetSubjectId())).
				SetDescriptionCode(15006).
				AddData("signal", signal))

			if err := o.transfer(transferer, message); err != nil {
				o.LogError(err.Error(), 15007)
				return false
			}
		}
	}

	o.LogReport(timeStart)
	return true
}

func (o *FileTransferer) transfererByType(tType string) (transferer transferers.AbstractTransfererContract, err error) {
	var ok bool
	transferer, ok = o.transferers[tType]
	if !ok {
		transferer, err = o.factory.FileTransferByType(tType)
	}
	return
}

func (o *FileTransferer) transfer(transferer transferers.AbstractTransfererContract, message modelsDataSqs.FileTransfererMessage) error {
	var err error
	var processedFiles = 0

	// determine all source files and destination folders
	// and other parameters from message and configuration
	if err = transferer.Init(message, o.Logger()); err != nil {
		return err
	}

	// transfer all source file...
	for _, file := range transferer.Files() {
		// ...to all destination folders
		for _, folder := range transferer.Folders() {

			// reset source storage before new file downloading for every folder
			if err = transferer.SetSourceFile(file); err != nil {
				return err
			}
			// reset destination storage before uploading to the folder
			if err = transferer.SetDestinationFolder(folder); err != nil {
				return err
			}

			// checkings
			if err = transferer.CheckSource(); err != nil {
				// no errors, just skip if there's no source file
				// (it allows to process all source files)
				if sErr, ok := err.(errorsStorage.ErrorContract); ok && sErr.Code() == errorsStorage.NotFoundError {
					continue
				}

				return err
			}
			if err = transferer.CheckDestination(); err != nil {
				return err
			}

			// transferring
			if err = transferer.Start(); err != nil {
				return err
			}

			var data []byte
			for {
				if data, err = transferer.DownloadData(); err != nil {
					transferer.Revert()
					return err
				}
				if data == nil {
					break
				}
				if err = transferer.UploadData(data); err != nil {
					transferer.Revert()
					return err
				}
			}

			if err = transferer.CloseDestination(); err != nil {
				if sErr, ok := err.(errorsStorage.ErrorContract); ok && sErr.Code() == errorsStorage.AccessSharingError {
					// we can ignore access sharing error
				} else {
					return err
				}
			}
		}

		// close source storage only when source file is transferred to all folders
		if err = transferer.CloseSource(); err != nil {
			return err
		}

		processedFiles++
	}

	if processedFiles == 0 {
		err = errors.New("No files were processed")
		o.Logger().Log(logger.CreateError(fmt.Sprintf("%s [%s]: %s", o.WorkerName(), o.Logger().GetSubjectId(), err.Error())).
			SetErrorCode(15007).
			AddData("files", transferer.Files()))
		return err
	}

	transferer.Reset()

	return err
}

func (o *FileTransferer) HealthCheck() error {
	if err := o.QueueReader.HealthCheck(); err != nil {
		return err
	}
	return nil
}
