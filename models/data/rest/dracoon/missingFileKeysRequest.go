package modelsRestDracoon

import "strconv"

type MissingFileKeysRequest struct {
	Offset int `json:"offset"`
	Limit  int `json:"limit"`
	RoomId int `json:"room_id"`
	FileId int `json:"file_id"`
}

func (r *MissingFileKeysRequest) Map() map[string]string {
	m := make(map[string]string)
	m["offset"] = strconv.Itoa(r.Offset)

	if r.Limit > 0 {
		m["limit"] = strconv.Itoa(r.Limit)
	}
	if r.RoomId > 0 {
		m["room_id"] = strconv.Itoa(r.RoomId)
	}
	if r.FileId > 0 {
		m["file_id"] = strconv.Itoa(r.FileId)
	}
	return m
}
