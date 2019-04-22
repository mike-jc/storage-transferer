package services

import (
	"errors"
	"github.com/astaxie/beego"
	"gitlab.com/24sessions/lib-go-logger/logger/services"
	"gitlab.com/24sessions/sdk-go-configurator"
	"gitlab.com/24sessions/sdk-go-configurator/data"
	"reflect"
	"service-recordingStorage/models/data/sqs"
	"strconv"
	"strings"
	"sync"
	"time"
)

const PseudoDomain = "24sessions.pseudo.com"
const ConfigExpiration = 5 * time.Minute

type configurationData struct {
	data     *sdksData.InstanceContainer
	expireAt time.Time
}

var testConfig *sdksData.InstanceContainer
var configurationsLock sync.Mutex
var configurations = make(map[string]configurationData)

func prodConfiguration(domain string) (config *sdksData.InstanceContainer, err error) {
	return new(sdks.Configurator).GetInstance(domain)
}

func testConfiguration() *sdksData.InstanceContainer {

	if testConfig == nil {
		testConfig = &sdksData.InstanceContainer{
			CompanyLocale:   "en",
			CompanyName:     "24sessions Test",
			CompanyTimezone: "UTC",
			CompanyStatus:   "trial",
		}

		if parentRoomId, err := beego.AppConfig.Int("dracoon.test.parentRoomId"); err != nil {
			LogMain.LogFatal(logger.CreateError("Wrong dracoon.test.parentRoomId in app configuration: " + err.Error()).SetErrorCode(15018))
			return nil
		} else {
			testConfig.DracoonParentRoomId = []byte(strconv.Itoa(parentRoomId))
		}

		baseUrl := beego.AppConfig.String("dracoon.test.baseUrl")
		if baseUrl == "" {
			LogMain.LogFatal(logger.CreateError("Empty dracoon.test.baseUrl in app configuration").SetErrorCode(15018))
			return nil
		} else {
			testConfig.DracoonBaseUrl = baseUrl
		}

		login := beego.AppConfig.String("dracoon.test.login")
		if login == "" {
			LogMain.LogFatal(logger.CreateError("Empty dracoon.test.login in app configuration").SetErrorCode(15018))
			return nil
		} else {
			testConfig.DracoonLogin = login
		}

		password := beego.AppConfig.String("dracoon.test.password")
		if password == "" {
			LogMain.LogFatal(logger.CreateError("Empty dracoon.test.password in app configuration").SetErrorCode(15018))
			return nil
		} else {
			testConfig.DracoonPassword = password
		}

		encPassword := beego.AppConfig.String("dracoon.test.encryptionPassword")
		if encPassword == "" {
			LogMain.LogFatal(logger.CreateError("Empty dracoon.test.encryptionPassword in app configuration").SetErrorCode(15018))
			return nil
		} else {
			testConfig.DracoonEncryptionPassword = encPassword
		}
	}

	return testConfig
}

func Configuration(domain string) (config *sdksData.InstanceContainer, err error) {
	configurationsLock.Lock()
	defer configurationsLock.Unlock()

	if beego.BConfig.RunMode == "test" || domain == PseudoDomain {
		config = testConfiguration()
	} else {
		if configData, ok := configurations[domain]; !ok || configData.expireAt.Before(time.Now()) {
			if config, err = prodConfiguration(domain); err == nil {
				configurations[domain] = configurationData{
					data:     config,
					expireAt: time.Now().Add(ConfigExpiration),
				}
			}
		} else {
			config = configData.data
		}
	}
	return
}

func MessageDomain(message interface{}) string {
	defaultDomain := PseudoDomain

	switch v := message.(type) {
	case nil:
		return defaultDomain
	case *modelsDataSqs.FileTransfererMessage:
		if v == nil {
			return defaultDomain
		} else {
			return message.(*modelsDataSqs.FileTransfererMessage).Instance.Domain
		}
	case modelsDataSqs.FileTransfererMessage:
		return message.(modelsDataSqs.FileTransfererMessage).Instance.Domain
	default:
		LogMain.LogFatal(logger.CreateError("Unknown type of message for getting instance domain: " + reflect.TypeOf(message).Name()).SetErrorCode(15018))
		return defaultDomain
	}
}

func DracoonBaseUrl(domain string) (url string, err error) {
	var config *sdksData.InstanceContainer
	if config, err = Configuration(domain); err == nil {
		if config.DracoonBaseUrl == "" {
			err = errors.New("Wrong Dracoon base URL in configuration of " + domain)
			return
		} else {
			url = strings.Trim(config.DracoonBaseUrl, "/")
		}
	}
	return
}

func DracoonParentRoomId(domain string) (id int, err error) {
	var config *sdksData.InstanceContainer
	if config, err = Configuration(domain); err == nil {
		id = config.GetDracoonParentRoomId()
	}
	return
}

func DracoonLogin(domain string) (login string, err error) {
	var config *sdksData.InstanceContainer
	if config, err = Configuration(domain); err == nil {
		if config.DracoonLogin == "" {
			err = errors.New("Wrong Dracoon login in configuration of " + domain)
			return
		} else {
			login = config.DracoonLogin
		}
	}
	return
}

func DracoonPassword(domain string) (password string, err error) {
	var config *sdksData.InstanceContainer
	if config, err = Configuration(domain); err == nil {
		if config.DracoonPassword == "" {
			err = errors.New("Wrong Dracoon password in configuration of " + domain)
			return
		} else {
			password = config.DracoonPassword
		}
	}
	return
}

func DracoonEncryptionPassword(domain string) (password string, err error) {
	var config *sdksData.InstanceContainer
	if config, err = Configuration(domain); err == nil {
		if config.DracoonEncryptionPassword == "" {
			err = errors.New("Wrong Dracoon encryption password in configuration of " + domain)
			return
		} else {
			password = config.DracoonEncryptionPassword
		}
	}
	return
}
