package servicesRest

import (
	"github.com/astaxie/beego/httplib"
	"gitlab.com/24sessions/lib-go-logger/logger/services"
	"net/http/httputil"
	"runtime/debug"
	"service-recordingStorage/errors/rest"
)

const RequestDumpMaxLength = 1000

type AbstractRestClient struct {
	logger *logger.Logger
}

func (c *AbstractRestClient) tuneRequest(request *httplib.BeegoHTTPRequest, headers interface{}) {
	request.Debug(false)

	if headers != nil {
		if headers, ok := headers.(map[string]string); ok {
			for name, value := range headers {
				request.Header(name, value)
			}
		}
	}
}

func (c *AbstractRestClient) addQueryParams(request *httplib.BeegoHTTPRequest, origParams interface{}) {
	if origParams == nil {
		return
	}

	if params, ok := origParams.(map[string]string); ok {
		for name, value := range params {
			request.Param(name, value)
		}
	} else if params, ok := origParams.(map[string][]string); ok {
		for name, values := range params {
			for _, subValue := range values {
				request.Param(name, subValue)
			}
		}
	}
}

func (c *AbstractRestClient) requestDump(request *httplib.BeegoHTTPRequest) string {
	cutDump := func(d []byte) string {
		if len(d) > RequestDumpMaxLength {
			return string(d[:RequestDumpMaxLength])
		} else {
			return string(d)
		}
	}

	if dump := request.DumpRequest(); len(dump) > 0 {
		return cutDump(dump)
	} else if dump, err := httputil.DumpRequest(request.GetRequest(), true); err == nil {
		return cutDump(dump)
	}
	return ""
}

func (c *AbstractRestClient) sendRequest(request *httplib.BeegoHTTPRequest) (output string, err errorsRest.ErrorContract) {
	requestDump := c.requestDump(request)

	response, reqErr := request.Response()
	if response != nil && response.Body != nil {
		defer response.Body.Close()
	}
	if reqErr != nil {
		err = errorsRest.NewError(reqErr.Error(), 0, "")
		c.logger.Log(logger.CreateError(err.Error()).SetErrorCode(15010).
			AddData("requestDump", requestDump).
			AddData("url", request.GetRequest().URL).
			AddData("response", err.Response()).
			AddData("stackTrace", string(debug.Stack())))
		return
	}

	if output, reqErr = request.String(); reqErr != nil {
		err = errorsRest.NewError(reqErr.Error(), 0, "")
	}

	if response.StatusCode < 200 || response.StatusCode >= 300 {
		err = errorsRest.NewError("Not OK status", response.StatusCode, output)
		c.logger.Log(logger.CreateError(err.Error()).SetErrorCode(15010).
			AddData("requestDump", requestDump).
			AddData("url", request.GetRequest().URL).
			AddData("response", err.Response()).
			AddData("stackTrace", string(debug.Stack())))
		return
	}
	return
}

func (c *AbstractRestClient) SetLogger(logger *logger.Logger) {
	c.logger = logger
}

func (c *AbstractRestClient) Logger() *logger.Logger {
	return c.logger
}
