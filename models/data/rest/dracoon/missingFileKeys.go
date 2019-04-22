package modelsRestDracoon

import "service-recordingStorage/models/data/accessSharer"

type MissingFileKeys struct {
	Items map[int]map[int]bool
	Users map[int]*modelsDataAccessSharer.PublicKeyContainer
	Files map[int]*FileKey
}
