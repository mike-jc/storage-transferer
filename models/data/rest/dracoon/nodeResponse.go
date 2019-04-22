package modelsRestDracoon

import (
	"strings"
	"time"
)

type NodeResponse struct {
	Id                 int             `json:"id"`
	Type               string          `json:"type"`
	Name               string          `json:"name"`
	ParentId           int             `json:"parentId"`
	ParentPath         string          `json:"parentPath"`
	ExpireAt           time.Time       `json:"expireAt"`
	Permissions        NodePermissions `json:"permissions"`
	IsEncrypted        bool            `json:"isEncrypted"`
	Classification     int             `json:"classification"`
	Size               int             `json:"size"`
	Quota              int             `json:"quota"`
	InheritPermissions bool            `json:"inheritPermissions"`
}

type NodePermissions struct {
	Manage bool `json:"manage"`
	Read   bool `json:"read"`
	Create bool `json:"create"`
	Change bool `json:"change"`
}

func (n *NodeResponse) FullPath() string {
	return n.ParentPath + n.Name
}

func (n *NodeResponse) RelativeFullPath(parentRoomPath string) string {
	relativePath := strings.Trim(strings.Replace(n.ParentPath, parentRoomPath, "", 1), "/")
	if relativePath == "" {
		return "/" + n.Name
	} else {
		return "/" + relativePath + "/" + n.Name
	}
}
