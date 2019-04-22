package modelsRestDracoon

import "service-recordingStorage/models/data/accessSharer"

type MissingFileKeysResponse struct {
	Range ResponseRange          `json:"range"`
	Items []*MissingFileKeysItem `json:"items"`
	Users []*MissingFileKeysUser `json:"users"`
	Files []*MissingFileKeysFile `json:"files"`
}

type MissingFileKeysItem struct {
	UserId int `json:"userId"`
	FileId int `json:"fileId"`
}

type MissingFileKeysUser struct {
	Id                 int                                       `json:"id"`
	PublicKeyContainer modelsDataAccessSharer.PublicKeyContainer `json:"publicKeyContainer"`
}

type MissingFileKeysFile struct {
	Id               int     `json:"id"`
	FileKeyContainer FileKey `json:"fileKeyContainer"`
}
