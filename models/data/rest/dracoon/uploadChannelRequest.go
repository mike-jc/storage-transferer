package modelsRestDracoon

import "time"

type UploadChannelRequest struct {
	ParentId   int                     `json:"parentId"`
	Name       string                  `json:"name"`
	Size       int                     `json:"size"`
	Expiration UploadChannelExpiration `json:"expiration"`
	Notes      string                  `json:"notes"`
}

type UploadChannelExpiration struct {
	EnableExpiration bool      `json:"enableExpiration"`
	ExpireAt         time.Time `json:"expireAt"`
}
