package modelsRestDracoon

type ErrorResponse struct {
	Message string `json:"message"`
	Info    string `json:"debugInfo"`
}
