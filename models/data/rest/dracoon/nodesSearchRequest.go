package modelsRestDracoon

import (
	"strconv"
)

type NodesSearchRequest struct {
	ParentId     int    `json:"parent_id"`
	DepthLevel   int    `json:"depth_level"`
	SearchString string `json:"search_string"`
	Filter       string `json:"filter"`
}

func (r *NodesSearchRequest) Map() map[string]string {
	m := make(map[string]string)
	m["depth_level"] = strconv.Itoa(r.DepthLevel)

	if r.ParentId >= 0 {
		m["parent_id"] = strconv.Itoa(r.ParentId)
	}
	if len(r.SearchString) > 0 {
		m["search_string"] = r.SearchString
	}
	if len(r.Filter) > 0 {
		m["filter"] = r.Filter
	}
	return m
}
