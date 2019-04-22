package amazon

import (
	"github.com/astaxie/beego"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"gitlab.com/24sessions/lib-go-logger/logger/services"
	"service-recordingStorage/services"
)

type AbstractService struct {
	session *session.Session
}

func (o *AbstractService) Session() *session.Session {
	if o.session == nil {
		region := beego.AppConfig.String("aws.region")
		if region == "" {
			services.LogMain.LogFatal(logger.CreateError("Incorrect AWS region in app configuration").SetErrorCode(15003))
			return nil
		}

		var err error
		o.session, err = session.NewSession(&aws.Config{
			Region: aws.String(region),
		})
		if err != nil {
			o.session = nil
			services.LogMain.LogFatal(logger.CreateError("Can not create AWS session: " + err.Error()).SetErrorCode(15003))
		}
	}
	return o.session
}
