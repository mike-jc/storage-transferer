package modelsRestDracoon

import "service-recordingStorage/models/data/accessSharer"

type UserKeyPairResponse struct {
	PrivateKeyContainer modelsDataAccessSharer.PrivateKeyContainer `json:"privateKeyContainer"`
	PublicKeyContainer  modelsDataAccessSharer.PublicKeyContainer  `json:"publicKeyContainer"`
}
