package tests

// make sure that initialization (`tests/init` import) is done before other imports
// noinspection GoUnsortedImport
import (
	"gitlab.com/24sessions/lib-go-logger/logger/services"
	"service-recordingStorage/models/data/instance"
	"service-recordingStorage/resources"
	_ "service-recordingStorage/tests/init" // keep it in the first row
	"time"
)

var LogMain *logger.Logger

func Init() {
	logger.OverrideConfigAsDummy(logger.LOG_LEVEL_DEBUG)
	LogMain = resources.InitLogger(resources.AppTypeTest)
}

func TestMeeting() modelsDataInstance.Meeting {
	return modelsDataInstance.Meeting{
		Id:       12345,
		DateUTC:  time.Now(),
		Timezone: "CET",
		Type: modelsDataInstance.MeetingType{
			Name:     "Test meeting",
			Duration: 30,
		},
		User: modelsDataInstance.User{
			Name:  "Test user",
			Email: "test.user@gmail.com",
		},
		Guest: modelsDataInstance.Guest{
			Name:  "Test guest",
			Email: "test.guest@gmail.com",
		},
	}
}

func TestDracoonTargets() []interface{} {
	return []interface{}{
		map[string]interface{}{
			"path":       "Data from unit tests/Room One",
			"expiration": "1 day",
		},
		map[string]interface{}{
			"path":       "Data from unit tests/Room Two",
			"expiration": "1 week",
		},
	}
}
