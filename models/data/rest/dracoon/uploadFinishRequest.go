package modelsRestDracoon

type UploadFinishRequest struct {
	ResolutionStrategy string          `json:"resolutionStrategy"`
	KeepShareLinks     bool            `json:"keepShareLinks"`
	FileName           string          `json:"fileName"`
	FileKey            FileKey         `json:"fileKey"`
	UserFileKeyList    UserFileKeyList `json:"userFileKeyList"`
}
