package modelsRestDracoon

type UploadResponse struct {
	TotalSize int    `json:"size"`
	Hash      string `json:"hash"`
}
