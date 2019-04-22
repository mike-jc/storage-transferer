package servicesRestApi

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"gitlab.com/24sessions/lib-go-logger/logger/services"
	"math/rand"
	"service-recordingStorage/errors/rest"
	"service-recordingStorage/models/data/accessSharer"
	"service-recordingStorage/models/data/rest"
	"service-recordingStorage/models/data/rest/dracoon"
	"service-recordingStorage/models/data/storage/dracoon"
	"service-recordingStorage/services/rest"
	"service-recordingStorage/system"
	"strings"
	"time"
)

const TokenExpirationPeriod = 2 * time.Hour

const UploadMaxRetries = 3
const UploadRetryDelay = 10 * time.Second

type Dracoon struct {
	baseUrl  string
	login    string
	password string

	token struct {
		value    string
		expireAt time.Time
	}

	childRooms map[int][]*modelsRestDracoon.NodeResponse
	nodes      map[int]*modelsRestDracoon.NodeResponse

	client       *servicesRest.RestJsonClient
	clientBinary *servicesRest.RestJsonBinaryClient
}

func (r *Dracoon) SetLogger(logger *logger.Logger) {
	r.Client().SetLogger(logger)
	r.ClientBinary().SetLogger(logger)
}

func (r *Dracoon) Client() *servicesRest.RestJsonClient {
	if r.client == nil {
		r.client = new(servicesRest.RestJsonClient)
	}
	return r.client
}

func (r *Dracoon) ClientBinary() *servicesRest.RestJsonBinaryClient {
	if r.clientBinary == nil {
		r.clientBinary = new(servicesRest.RestJsonBinaryClient)
	}
	return r.clientBinary
}

func (r *Dracoon) SetBaseUrl(url string) {
	r.token.value = "" // force token re-request since login and password are changed
	r.baseUrl = url
}

func (r *Dracoon) SetCredentials(login, password string) {
	r.token.value = "" // force token re-request since login and password are changed
	r.login = login
	r.password = password
}

func (r *Dracoon) Reset() {
	r.childRooms = make(map[int][]*modelsRestDracoon.NodeResponse)
	r.nodes = make(map[int]*modelsRestDracoon.NodeResponse)
}

func (r *Dracoon) absoluteUrl(url string) string {
	return fmt.Sprintf("%s/%s", r.baseUrl, url)
}

func (r *Dracoon) IsExpired(t time.Time) bool {
	if t.Unix() <= 0 {
		return false
	}
	return time.Now().Add(5 * time.Minute).After(t)
}

func (r *Dracoon) TokenValue() (value string, err errorsRest.DracoonErrorContract) {

	// get new token from Dracoon if there's no any token or it's expired
	if r.token.value == "" || r.IsExpired(r.token.expireAt) {
		var params = modelsRestDracoon.AuthRequest{
			Login:    r.login,
			Password: r.password,
		}
		var response modelsRestDracoon.AuthResponse

		if cErr := r.Client().Post(r.absoluteUrl("auth/login"), params, nil, &response); cErr == nil {
			r.token.value = response.Token
			r.token.expireAt = time.Now().Add(TokenExpirationPeriod)
		} else {
			cErr.SetError("Can not get auth token from Dracoon: " + cErr.Error())
			err = errorsRest.NewDracoonErrorFromAbstract(cErr)
			return
		}
	} else {
		// every usage resets expiration period
		r.token.expireAt = time.Now().Add(TokenExpirationPeriod)
	}

	value = r.token.value
	return
}

func (r *Dracoon) actualizeCache(node *modelsRestDracoon.NodeResponse) {
	r.nodes[node.Id] = node

	if r.childRooms[node.ParentId] == nil {
		r.childRooms[node.ParentId] = make([]*modelsRestDracoon.NodeResponse, 0)
	}
	r.childRooms[node.ParentId] = append(r.childRooms[node.ParentId], node)
}

func (r *Dracoon) CreateNode(parentNodeId int, name string) (node *modelsRestDracoon.NodeResponse, err errorsRest.DracoonErrorContract) {

	var token string
	if token, err = r.TokenValue(); err != nil {
		err.SetError("Can not create node on Dracoon: " + err.Error())
		return
	}

	var url = r.absoluteUrl("nodes/rooms")
	var params = modelsRestDracoon.NodeRequest{
		Name:               name,
		ParentId:           parentNodeId,
		InheritPermissions: true,
	}

	var headers = map[string]string{
		"X-Sds-Auth-Token":  token,
		"X-Sds-Date-Format": "UTC",
	}
	node = new(modelsRestDracoon.NodeResponse)

	if cErr := r.Client().Post(url, params, headers, node); cErr != nil {
		cErr.SetError("Can not create node on Dracoon: " + cErr.Error())
		err = errorsRest.NewDracoonErrorFromAbstract(cErr)
		return
	}

	r.actualizeCache(node)
	return
}

func (r *Dracoon) ChildRooms(parentRoomId int) (rooms []*modelsRestDracoon.NodeResponse, err errorsRest.DracoonErrorContract) {

	var ok bool
	if rooms, ok = r.childRooms[parentRoomId]; ok {
		return
	} else {
		var token string
		if token, err = r.TokenValue(); err != nil {
			err.SetError("Can not get full list of rooms from Dracoon: " + err.Error())
			return
		}

		var url = r.absoluteUrl("nodes/search")
		var params = modelsRestDracoon.NodesSearchRequest{
			SearchString: "*",            // any
			DepthLevel:   -1,             // full tree
			Filter:       "type:eq:room", // only rooms
			ParentId:     parentRoomId,
		}
		var headers = map[string]string{
			"X-Sds-Auth-Token":  token,
			"X-Sds-Date-Format": "UTC",
		}
		var list modelsRestDracoon.NodesResponse

		if cErr := r.Client().Get(url, params.Map(), headers, &list); cErr != nil {
			cErr.SetError("Can not get full list of rooms from Dracoon: " + cErr.Error())
			err = errorsRest.NewDracoonErrorFromAbstract(cErr)
			return
		}

		rooms = list.Items
		r.childRooms[parentRoomId] = list.Items
		return
	}
}

func (r *Dracoon) Node(nodeId int) (node *modelsRestDracoon.NodeResponse, err errorsRest.DracoonErrorContract) {

	var ok bool
	if node, ok = r.nodes[nodeId]; ok {
		return
	} else {
		var token string
		if token, err = r.TokenValue(); err != nil {
			err.SetError("Can not get node info from Dracoon: " + err.Error())
			return
		}

		var url = r.absoluteUrl(fmt.Sprintf("nodes/%d", nodeId))
		var headers = map[string]string{
			"X-Sds-Auth-Token":  token,
			"X-Sds-Date-Format": "UTC",
		}
		node = new(modelsRestDracoon.NodeResponse)

		if cErr := r.Client().Get(url, nil, headers, node); cErr != nil {
			cErr.SetError("Can not get node info from Dracoon: " + cErr.Error())
			err = errorsRest.NewDracoonErrorFromAbstract(cErr)
			return
		}

		r.nodes[nodeId] = node
		return
	}
}

// @Description find node or create if it does not exist
func (r *Dracoon) FindOrCreateRoom(parentNodeId int, path string) (room *modelsRestDracoon.NodeResponse, err errorsRest.DracoonErrorContract) {
	// clear caches - make sure that previously created room are visible here
	r.Reset()

	// get parent room info
	var parentRoom *modelsRestDracoon.NodeResponse
	if parentRoom, err = r.Node(parentNodeId); err != nil {
		return
	}

	// get child rooms of parent room (recursively)
	var rooms []*modelsRestDracoon.NodeResponse
	if rooms, err = r.ChildRooms(parentNodeId); err != nil {
		return
	}

	// collect relative paths to make search faster
	paths := make(map[string]*modelsRestDracoon.NodeResponse)
	parentRoomPath := parentRoom.FullPath()
	paths[parentRoomPath] = parentRoom
	for _, room := range rooms {
		paths[room.FullPath()] = room
	}

	// try to find by path
	path = parentRoomPath + "/" + strings.Trim(path, "/")
	var ok bool

	if room, ok = paths[path]; ok {
		return
	} else {
		// firstly, find the deepest parent that exists
		// and if not, create all other parents and the node itself
		var parentPath string
		var parentNodeId int
		var prevParentNode, parentNode *modelsRestDracoon.NodeResponse

		for _, part := range strings.Split(path, "/") {
			if part == "" {
				continue
			}

			parentPath += "/" + part

			if parentNode, ok = paths[parentPath]; !ok {
				if prevParentNode == nil {
					parentNodeId = 0 // root room
				} else {
					parentNodeId = prevParentNode.Id
				}

				// path doesn't exist, create its node
				if room, err = r.CreateNode(parentNodeId, part); err != nil {
					return
				}
				prevParentNode = room
			} else {
				prevParentNode = parentNode
			}
		}
	}
	return
}

func (r *Dracoon) CreateUploadChannel(target modelsDataStorageDracoon.Target, fileName string, fileSize int, notes string) (channel modelsRestDracoon.UploadChannelResponse, err errorsRest.DracoonErrorContract) {
	var token string
	if token, err = r.TokenValue(); err != nil {
		err.SetError("Can not open upload channel on Dracoon: " + err.Error())
		return
	}

	var url = r.absoluteUrl("nodes/files/uploads")
	var params = modelsRestDracoon.UploadChannelRequest{
		ParentId: target.RoomId,
		Name:     fileName,
		Size:     fileSize,
		Notes:    notes,
	}
	if len(target.Expiration) > 0 {
		params.Expiration = modelsRestDracoon.UploadChannelExpiration{
			EnableExpiration: true,
			ExpireAt:         system.AddDurationFromString(time.Now().UTC(), target.Expiration),
		}
	}
	var headers = map[string]string{
		"X-Sds-Auth-Token": token,
	}

	if cErr := r.Client().Post(url, params, headers, &channel); cErr != nil {
		cErr.SetError("Can not open upload channel on Dracoon: " + cErr.Error())
		err = errorsRest.NewDracoonErrorFromAbstract(cErr)
		return
	}
	return
}

func (r *Dracoon) Upload(uploadToken string, data []byte, contentRange string) errorsRest.DracoonErrorContract {

	var url = r.absoluteUrl(fmt.Sprintf("uploads/%s", uploadToken))
	var headers = make(map[string]string)
	if contentRange != "" {
		headers["Content-Range"] = contentRange
	}
	var response modelsRestDracoon.UploadResponse

	dataMd5 := md5.Sum(data)
	hash := hex.EncodeToString(dataMd5[:])
	retries := 0

	for {
		retries++
		if cErr := r.ClientBinary().Post(url, data, headers, &response); cErr != nil {
			cErr.SetError("Can not upload data to Dracoon room: " + cErr.Error())
			return errorsRest.NewDracoonErrorFromAbstract(cErr)
		}
		if response.Hash == hash {
			return nil
		} else if retries >= UploadMaxRetries {
			return errorsRest.NewDracoonError(fmt.Sprintf("Can not upload data to Dracoon room: transmission error (hash is [%s], should be [%s])", response.Hash, hash), modelsRest.DracoonRestResponse{})
		} else {
			time.Sleep(UploadRetryDelay + time.Duration(rand.Intn(1000))*time.Millisecond)
		}
	}
}

func (r *Dracoon) UserKeyPair() (keyPair modelsRestDracoon.UserKeyPairResponse, err errorsRest.DracoonErrorContract) {
	var token string
	if token, err = r.TokenValue(); err != nil {
		err.SetError("Can not get current user's key pair: " + err.Error())
		return
	}

	var url = r.absoluteUrl("user/account/keypair")
	var headers = map[string]string{
		"X-Sds-Auth-Token": token,
	}

	if cErr := r.Client().Get(url, nil, headers, &keyPair); cErr != nil {
		cErr.SetError("Can not get current user's key pair: " + cErr.Error())
		err = errorsRest.NewDracoonErrorFromAbstract(cErr)
		return
	}
	return
}

func (r *Dracoon) FinishUploading(uploadToken string, fileName string, fileKey modelsRestDracoon.FileKey) (file modelsRestDracoon.NodeResponse, err errorsRest.DracoonErrorContract) {

	var url = r.absoluteUrl(fmt.Sprintf("uploads/%s", uploadToken))
	var params = modelsRestDracoon.UploadFinishRequest{
		ResolutionStrategy: "autorename",
		FileName:           fileName,
	}
	var headers = map[string]string{
		"X-Sds-Date-Format": "UTC",
	}

	if fileKey.Valid() {
		params.FileKey = fileKey
	}

	if cErr := r.Client().Put(url, params, headers, &file); cErr != nil {
		cErr.SetError("Can not finish uploading on Dracoon: " + cErr.Error())
		err = errorsRest.NewDracoonErrorFromAbstract(cErr)
		return
	}
	return
}

func (r *Dracoon) CancelUploading(uploadToken string) errorsRest.DracoonErrorContract {
	var url = r.absoluteUrl(fmt.Sprintf("uploads/%s", uploadToken))

	if cErr := r.Client().Delete(url, nil); cErr != nil {
		cErr.SetError("Can not cancel uploading on Dracoon: " + cErr.Error())
		return errorsRest.NewDracoonErrorFromAbstract(cErr)
	}
	return nil
}

func (r *Dracoon) MissingFileKeys(roomId int, fileId int) (keys modelsRestDracoon.MissingFileKeys, err errorsRest.DracoonErrorContract) {
	var token string
	if token, err = r.TokenValue(); err != nil {
		err.SetError("Can not get missing file keys on Dracoon: " + err.Error())
		return
	}

	var url = r.absoluteUrl("nodes/missingFileKeys")
	var params = modelsRestDracoon.MissingFileKeysRequest{
		Limit: modelsRestDracoon.DefaultLimit,
	}
	if roomId > 0 {
		params.RoomId = roomId
	}
	if fileId > 0 {
		params.FileId = fileId
	}
	var headers = map[string]string{
		"X-Sds-Auth-Token": token,
	}
	var result modelsRestDracoon.MissingFileKeysResponse

	keys.Items = make(map[int]map[int]bool)
	keys.Users = make(map[int]*modelsDataAccessSharer.PublicKeyContainer)
	keys.Files = make(map[int]*modelsRestDracoon.FileKey)

	for {
		cErr := r.Client().Get(url, params.Map(), headers, &result)

		// check range
		if cErr != nil {
			cErr.SetError("Can not get missing file keys on Dracoon: " + cErr.Error())
			err = errorsRest.NewDracoonErrorFromAbstract(cErr)
			return
		} else if len(result.Items) == 0 {
			break
		}

		// collect result
		for _, item := range result.Items {
			if keys.Items[item.UserId] == nil {
				keys.Items[item.UserId] = make(map[int]bool)
			}
			keys.Items[item.UserId][item.FileId] = true
		}
		for _, user := range result.Users {
			keys.Users[user.Id] = &user.PublicKeyContainer
		}
		for _, file := range result.Files {
			keys.Files[file.Id] = &file.FileKeyContainer
		}

		// move to the next range
		params.Offset = result.Range.Total
	}

	return
}

func (r *Dracoon) UploadFileKeys(params modelsRestDracoon.UploadFileKeysRequest) errorsRest.DracoonErrorContract {
	var token string
	var err errorsRest.DracoonErrorContract

	if token, err = r.TokenValue(); err != nil {
		err.SetError("Can not upload missing file keys to Dracoon: " + err.Error())
		return err
	}

	var url = r.absoluteUrl("nodes/files/keys")
	var headers = map[string]string{
		"X-Sds-Auth-Token": token,
	}

	if cErr := r.Client().Post(url, params, headers, nil); cErr != nil {
		return errorsRest.NewDracoonErrorFromAbstract(cErr)
	}
	return nil
}

func (r *Dracoon) Ping() errorsRest.DracoonErrorContract {

	var url = r.absoluteUrl("auth/ping")
	if cErr := r.Client().Get(url, nil, nil, nil); cErr != nil {
		cErr.SetError("Can not ping to Dracoon: " + cErr.Error())
		return errorsRest.NewDracoonErrorFromAbstract(cErr)
	}
	return nil
}
