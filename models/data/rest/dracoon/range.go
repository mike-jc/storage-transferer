package modelsRestDracoon

const DefaultLimit = 100

type ResponseRange struct {
	Offset int `json:"offset"`
	Limit  int `json:"limit"`
	Total  int `json:"total"`
}
