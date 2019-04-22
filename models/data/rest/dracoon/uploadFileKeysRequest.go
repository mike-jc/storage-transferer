package modelsRestDracoon

type UploadFileKeysRequest struct {
	Items []*UploadFileKeysItem `json:"items"`
}

type UploadFileKeysItem struct {
	FileId  int     `json:"fileId"`
	UserId  int     `json:"userId"`
	FileKey FileKey `json:"fileKey"`
}
