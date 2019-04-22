package modelsRestDracoon

type NodeRequest struct {
	Name               string `json:"name"`
	ParentId           int    `json:"parentId"`
	Quota              int    `json:"quota"`
	Notes              string `json:"notes"`
	InheritPermissions bool   `json:"inheritPermissions"`
}
