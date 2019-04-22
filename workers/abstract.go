package workers

import (
	"fmt"
	"gitlab.com/24sessions/lib-go-logger/logger"
	"gitlab.com/24sessions/lib-go-logger/logger/services"
	"service-recordingStorage/models/data/universe"
	"service-recordingStorage/resources"
	"time"
)

type AbstractWorker struct {
	logger      *logger.Logger
	subjectId   string
	reportLevel int

	self AbstractWorkerContract

	AbstractWorkerContract
}

type AbstractWorkerContract interface {
	WorkerName() string
	SubjectId() string
	Init(options interface{}) error
	Run()
	LogError(err string, code int)
	LogProfile(t time.Time)
	Status() interface{}
	ProcessSignal(signal interface{}) bool
	HealthCheck() error
}

func (o *AbstractWorker) SetLogger(logger *logger.Logger) {
	o.logger = logger
}

func (o *AbstractWorker) Logger() *logger.Logger {
	if o.logger == nil {
		o.logger = new(logger.Logger).SetSubject(o.self.WorkerName(), o.self.SubjectId()).SetParent(resources.LogMain)
	}
	return o.logger
}

func (o *AbstractWorker) SetInstance(instanceDomain string) {
	o.Logger().SetInstance(instanceDomain)
}

func (o *AbstractWorker) SetSubject(subjectDomain, subjectId string) {
	o.subjectId = subjectId
	o.Logger().SetSubject(subjectDomain, subjectId)
}

func (o *AbstractWorker) Init(options interface{}) error {
	o.self = o
	return nil
}

func (o *AbstractWorker) WorkerName() string {
	return "Abstract worker"
}

func (o *AbstractWorker) SubjectId() string {
	if o.subjectId == "" {
		o.subjectId = uuid.NewV4String()
	}
	return o.subjectId
}

func (o *AbstractWorker) Run() {
	o.Logger().Log(logger.CreateInfo(fmt.Sprintf("%s [%s] started", o.self.WorkerName(), o.Logger().GetSubjectId())))
}

func (o *AbstractWorker) LogError(err string, code int) {
	o.Logger().Log(logger.CreateError(fmt.Sprintf("%s [%s]: %s", o.self.WorkerName(), o.Logger().GetSubjectId(), err)).SetErrorCode(code))
}

func (o *AbstractWorker) LogReport(t time.Time) {
	if o.reportLevel >= modelsDataUniverse.ReportLevelProfile {
		message := fmt.Sprintf("%s [%s] passed: %s sec", o.self.WorkerName(), o.Logger().GetSubjectId(), time.Now().Sub(t).String())
		o.Logger().Log(logger.CreateInfo(message).SetDescriptionCode(15008))
	}
}

func (o *AbstractWorker) Status() interface{} {
	return nil
}

func (o *AbstractWorker) ProcessSignal(signal interface{}) bool {
	return true
}

func (o *AbstractWorker) HealthCheck() error {
	return nil
}
