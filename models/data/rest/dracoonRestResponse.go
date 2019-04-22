package modelsRest

import (
	"service-recordingStorage/models/data/rest/dracoon"
)

type DracoonRestResponse struct {
	Code         int
	Body         modelsRestDracoon.ErrorResponse
	ParsingError error
}
