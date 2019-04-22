package modelsDataStorageDracoon

type Target struct {
	Path       string `json:"path"`
	Expiration string `json:"expiration"`
	RoomId     int    `json:"-"`
}
