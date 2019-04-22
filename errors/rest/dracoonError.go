package errorsRest

import (
	"encoding/json"
	"errors"
	"fmt"
	"service-recordingStorage/models/data/rest"
	"service-recordingStorage/models/data/rest/dracoon"
)

type DracoonError struct {
	text     string
	response modelsRest.DracoonRestResponse

	DracoonErrorContract
}

type DracoonErrorContract interface {
	SetError(text string)
	Response() modelsRest.DracoonRestResponse

	error
}

func (e *DracoonError) SetError(text string) {
	e.text = text
}

func (e *DracoonError) Error() string {
	return e.text
}

func (e *DracoonError) Response() modelsRest.DracoonRestResponse {
	return e.response
}

func NewDracoonError(text string, response modelsRest.DracoonRestResponse) *DracoonError {
	return &DracoonError{
		text:     text,
		response: response,
	}
}

func NewDracoonErrorFromAbstract(err ErrorContract) *DracoonError {
	var jBody modelsRestDracoon.ErrorResponse
	var newResponse modelsRest.DracoonRestResponse

	response := err.Response()

	if jErr := json.Unmarshal([]byte(response.Body), &jBody); jErr == nil {
		newResponse = modelsRest.DracoonRestResponse{
			Code: response.Code,
			Body: jBody,
		}
	} else {
		newResponse = modelsRest.DracoonRestResponse{
			Code:         response.Code,
			ParsingError: errors.New(fmt.Sprintf("Can not parse the response: %s. Original body: %s", jErr.Error(), response.Body)),
		}
	}

	return &DracoonError{
		text:     err.Error(),
		response: newResponse,
	}
}
