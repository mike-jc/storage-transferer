package modelsRestDracoon

import (
	"encoding/base64"
	"service-recordingStorage/models/data/encryptor"
)

type FileKey struct {
	Key     string `json:"key"`
	IV      string `json:"iv"`
	Version string `json:"version"`
	Tag     string `json:"tag"`
}

type UserFileKey struct {
	UserId  int     `json:"userId"`
	FileKey FileKey `json:"fileKey"`
}

type UserFileKeyList struct {
	Items []*UserFileKey `json:"items"`
}

func (k *FileKey) Valid() bool {
	return len(k.Key) > 0 && len(k.IV) > 0 && len(k.Tag) > 0
}

func (k *FileKey) BinaryToString() FileKey {
	return FileKey{
		Key:     base64.StdEncoding.EncodeToString([]byte(k.Key)),
		IV:      base64.StdEncoding.EncodeToString([]byte(k.IV)),
		Version: k.Version,
		Tag:     base64.StdEncoding.EncodeToString([]byte(k.Tag)),
	}
}

func FileKeyFromEncryptionInfo(info modelsDataEncryptor.EncryptionInfo) FileKey {
	return FileKey{
		Key:     string(info.RandomKey),
		IV:      string(info.InitialVector),
		Version: "A",
		Tag:     string(info.AuthenticationTag),
	}
}
