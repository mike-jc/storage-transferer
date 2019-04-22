package workers

import (
	"errors"
	"fmt"
	"gitlab.com/24sessions/lib-go-logger/logger/services"
	"service-recordingStorage/models/data/accessSharer"
	"service-recordingStorage/models/data/encryptor"
	"service-recordingStorage/models/data/rest/dracoon"
	"service-recordingStorage/models/data/worker"
	"service-recordingStorage/models/data/worker/options"
	"service-recordingStorage/services/accessSharers"
	"service-recordingStorage/services/encryptors"
	"service-recordingStorage/services/rest/api"
)

type DracoonGlobalSharer struct {
	options   modelsDataWorkerOptions.DracoonGlobalSharerOptions
	apiClient *servicesRestApi.Dracoon
	sharer    *accessSharers.Dracoon

	InstanceJobRunner
}

func (o *DracoonGlobalSharer) WorkerName() string {
	return "DracoonGlobalSharer"
}

func (o *DracoonGlobalSharer) Init(options interface{}) error {

	if dOptions, ok := options.(modelsDataWorkerOptions.DracoonGlobalSharerOptions); !ok {
		return errors.New(fmt.Sprintf("Wrong options format for DracoonGlobalSharer. Should be modelsDataWorkerOptions.DracoonGlobalSharerOptions, but got <%T> %+v", options, options))
	} else {
		o.options = dOptions
	}

	if err := o.InstanceJobRunner.Init(o.options.InstanceJobRunnerOptions); err != nil {
		return err
	}

	o.self = o

	// API client
	o.apiClient = new(servicesRestApi.Dracoon)
	o.apiClient.SetLogger(o.Logger())

	// Sharer
	encryptor := new(encryptors.HybridCrypto)
	if err := encryptor.Init(modelsDataEncryptor.NewHybridCryptoOptionsForDracoon()); err != nil {
		return errors.New(fmt.Sprintf("Can not initialize Dracoon encryptor: %s", err.Error()))
	}

	o.sharer = new(accessSharers.Dracoon)
	o.sharer.SetClient(o.apiClient)
	o.sharer.SetEncryptor(encryptor)
	if err := o.sharer.Init(); err != nil {
		return errors.New(fmt.Sprintf("Can not initialize Dracoon sharer: %s", err.Error()))
	}
	return nil
}

func (o *DracoonGlobalSharer) ProcessSignal(signal interface{}) bool {
	if iSignal, ok := signal.(modelsDataWorker.InstanceJobRunnerSignal); !ok {
		o.LogError(fmt.Sprintf("Wrong signal format. Should be modelsDataWorker.InstanceJobRunnerSignal, got <%T> %+v", signal, signal), 15016)
		return false
	} else {
		if iSignal.Config.RecordingsStorage == "dracoon" {
			o.Logger().Log(logger.CreateInfo(fmt.Sprintf("%s [%s]: signal processing has started for %s", o.WorkerName(), o.Logger().GetSubjectId(), iSignal.Domain)).
				SetDescriptionCode(15017).
				AddData("signal", signal))

			if iSignal.Config.GetDracoonParentRoomId() == 0 {
				o.LogError(fmt.Sprintf("Instance %s should have recordings on Dracoon but Dracoon parent room ID is empty", iSignal.Domain), 15016)
				return false
			} else if len(iSignal.Config.DracoonBaseUrl) == 0 {
				o.LogError(fmt.Sprintf("Instance %s should have recordings on Dracoon but Dracoon base URL is empty", iSignal.Domain), 15016)
				return false
			} else if len(iSignal.Config.DracoonLogin) == 0 {
				o.LogError(fmt.Sprintf("Instance %s should have recordings on Dracoon but Dracoon account's login is empty", iSignal.Domain), 15016)
				return false
			} else if len(iSignal.Config.DracoonPassword) == 0 {
				o.LogError(fmt.Sprintf("Instance %s should have recordings on Dracoon but Dracoon account's password is empty", iSignal.Domain), 15016)
				return false
			} else if len(iSignal.Config.DracoonEncryptionPassword) == 0 {
				o.LogError(fmt.Sprintf("Instance %s should have recordings on Dracoon but Dracoon account's encryption password is empty", iSignal.Domain), 15016)
				return false
			} else {
				o.share(iSignal.Config.GetDracoonParentRoomId(), iSignal.Config.DracoonBaseUrl, iSignal.Config.DracoonLogin, iSignal.Config.DracoonPassword, iSignal.Config.DracoonEncryptionPassword)
				return true
			}
		}
		return true
	}
}

func (o *DracoonGlobalSharer) share(roomId int, baseUrl, login, password, encryptionPassword string) {
	o.apiClient.SetBaseUrl(baseUrl)
	o.apiClient.SetCredentials(login, password)
	o.apiClient.Reset()

	o.sharer.SetPassword(encryptionPassword)

	// Get auth user keys
	if keys, err := o.apiClient.UserKeyPair(); err != nil {
		o.LogError(fmt.Sprintf("Can not get Dracoon user's keys for this instance: %s", err.Error()), 15016)
		return
	} else {
		o.sharer.SetKeyPair(modelsDataAccessSharer.KeyPair{
			PrivateKeyContainer: keys.PrivateKeyContainer,
			PublicKeyContainer:  keys.PublicKeyContainer,
		})
	}

	// Collect all rooms
	roomMap := make(map[int]*modelsRestDracoon.NodeResponse)
	if parentRoom, err := o.apiClient.Node(roomId); err != nil {
		o.LogError(fmt.Sprintf("Can not get info of parent room %d on Dracoon: %s", roomId, err.Error()), 15016)
		return
	} else {
		roomMap[roomId] = parentRoom
	}
	if rooms, err := o.apiClient.ChildRooms(roomId); err != nil {
		o.LogError(fmt.Sprintf("Can not get child rooms on Dracoon for room %d: %s", roomId, err.Error()), 15016)
		return
	} else {
		for _, room := range rooms {
			roomMap[room.Id] = room
		}
	}

	// Check every room and share keys for it
	for roomId, room := range roomMap {
		err := o.sharer.Share(modelsDataAccessSharer.DracoonOptions{
			RoomId: roomId,
		})
		if err != nil {
			o.LogError(fmt.Sprintf("Error for Dracoon room %s [%d]: %s", room.FullPath(), roomId, err.Error()), 15016)
			// try to share access for other rooms even if we've got error here
		}
	}
}
