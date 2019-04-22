package servicesRest

import (
	"github.com/astaxie/beego/httplib"
	"service-recordingStorage/errors/rest"
)

type RestJsonBinaryClient struct {
	RestJsonClient
}

func (c *RestJsonBinaryClient) tuneRequest(request *httplib.BeegoHTTPRequest, headers interface{}) {
	c.RestJsonClient.tuneRequest(request, headers)

	request.Header("Content-Type", "application/octet-stream")
}

func (c *RestJsonBinaryClient) Post(url string, data []byte, headers interface{}, result interface{}) errorsRest.ErrorContract {
	request := httplib.Post(url)
	request.Body(data)
	c.tuneRequest(request, headers)

	if err := c.sendRequest(request, result); err != nil {
		return err
	}
	return nil
}
