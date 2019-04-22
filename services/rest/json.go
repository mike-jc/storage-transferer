package servicesRest

import (
	"encoding/json"
	"github.com/astaxie/beego/httplib"
	"service-recordingStorage/errors/rest"
)

type RestJsonClient struct {
	AbstractRestClient
}

func (c *RestJsonClient) tuneRequest(request *httplib.BeegoHTTPRequest, headers interface{}) {
	c.AbstractRestClient.tuneRequest(request, headers)

	request.Header("Accept", "application/json;charset=UTF-8")
	request.Header("Content-Type", "application/json;charset=UTF-8")
}

func (c *RestJsonClient) sendRequest(request *httplib.BeegoHTTPRequest, result interface{}) errorsRest.ErrorContract {
	output, err := c.AbstractRestClient.sendRequest(request)
	if err != nil {
		return err
	}

	if result != nil {
		if jErr := json.Unmarshal([]byte(output), result); jErr != nil {
			return errorsRest.NewError("Can not parse the response: "+jErr.Error(), 0, "")
		}
	}
	return nil
}

func (c *RestJsonClient) Post(url string, params interface{}, headers interface{}, result interface{}) errorsRest.ErrorContract {
	request := httplib.Post(url)
	c.tuneRequest(request, headers)

	if params != nil {
		body, jErr := json.Marshal(params)
		if jErr != nil {
			return errorsRest.NewError(jErr.Error(), 0, "")
		}
		request.Body(body)
	}

	if err := c.sendRequest(request, result); err != nil {
		return err
	}
	return nil
}

func (c *RestJsonClient) Put(url string, params interface{}, headers interface{}, result interface{}) errorsRest.ErrorContract {
	request := httplib.Put(url)
	c.tuneRequest(request, headers)

	if params != nil {
		body, jErr := json.Marshal(params)
		if jErr != nil {
			return errorsRest.NewError(jErr.Error(), 0, "")
		}
		request.Body(body)
	}

	if err := c.sendRequest(request, result); err != nil {
		return err
	}
	return nil
}

func (c *RestJsonClient) Get(url string, params interface{}, headers interface{}, result interface{}) errorsRest.ErrorContract {
	request := httplib.Get(url)
	c.tuneRequest(request, headers)
	c.addQueryParams(request, params)

	if err := c.sendRequest(request, result); err != nil {
		return err
	}
	return nil
}

func (c *RestJsonClient) Delete(url string, headers interface{}) errorsRest.ErrorContract {
	request := httplib.Delete(url)
	c.tuneRequest(request, headers)

	if err := c.sendRequest(request, nil); err != nil {
		return err
	}
	return nil
}
