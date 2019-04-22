package errorsRest

import (
	"service-recordingStorage/models/data/rest"
)

type Error struct {
	text     string
	response modelsRest.RestResponse

	ErrorContract
}

type ErrorContract interface {
	SetError(text string)
	Response() modelsRest.RestResponse

	error
}

func (e *Error) SetError(text string) {
	e.text = text
}

func (e *Error) Error() string {
	return e.text
}

func (e *Error) Response() modelsRest.RestResponse {
	return e.response
}

func NewError(text string, code int, body string) *Error {
	return &Error{
		text: text,
		response: modelsRest.RestResponse{
			Code: code,
			Body: body,
		},
	}
}
