package storages

import (
	"fmt"
	"gitlab.com/24sessions/lib-go-logger/logger/services"
	"service-recordingStorage/errors/storage"
	"service-recordingStorage/models/data/accessSharer"
	"service-recordingStorage/models/data/encryptor"
	"service-recordingStorage/models/data/rest/dracoon"
	"service-recordingStorage/models/data/sqs"
	"service-recordingStorage/models/data/storage"
	"service-recordingStorage/models/data/storage/dracoon"
	"service-recordingStorage/services"
	"service-recordingStorage/services/accessSharers"
	"service-recordingStorage/services/encryptors"
	"service-recordingStorage/services/rest/api"
	"service-recordingStorage/system"
	"strings"
)

type Dracoon struct {
	client *servicesRestApi.Dracoon

	target    modelsDataStorageDracoon.Target
	extraInfo string

	uploadChannel struct {
		id    string
		token string
	}

	AbstractStorage
}

func InitDracoonRestClient(client *servicesRestApi.Dracoon, message *modelsDataSqs.FileTransfererMessage) errorsStorage.ErrorContract {
	var baseUrl, login, password string
	var err error
	if baseUrl, err = services.DracoonBaseUrl(services.MessageDomain(message)); err != nil {
		return errorsStorage.NewError(err.Error(), errorsStorage.Message, errorsStorage.InstanceSettingsError)
	}
	if login, err = services.DracoonLogin(services.MessageDomain(message)); err != nil {
		return errorsStorage.NewError(err.Error(), errorsStorage.Message, errorsStorage.CredentialsError)
	}
	if password, err = services.DracoonPassword(services.MessageDomain(message)); err != nil {
		return errorsStorage.NewError(err.Error(), errorsStorage.Message, errorsStorage.CredentialsError)
	}
	client.SetBaseUrl(baseUrl)
	client.SetCredentials(login, password)
	client.Reset()
	return nil
}

func (s *Dracoon) StorageName() string {
	return "Dracoon storage"
}

func (s *Dracoon) Init(message modelsDataSqs.FileTransfererMessage, logger *logger.Logger) errorsStorage.ErrorContract {
	err := s.AbstractStorage.Init(message, logger)
	if err != nil {
		return err
	}

	s.self = s

	// client
	if s.client == nil {
		s.client = new(servicesRestApi.Dracoon)
	}
	s.client.SetLogger(logger)
	if err := InitDracoonRestClient(s.client, &message); err != nil {
		return err
	}

	// encryptor
	if s.encryptor == nil {
		encryptor := new(encryptors.HybridCrypto)
		if iErr := encryptor.Init(modelsDataEncryptor.NewHybridCryptoOptionsForDracoon()); iErr != nil {
			return errorsStorage.NewError(iErr.Error(), errorsStorage.Settings, errorsStorage.OptionsError)
		}
		s.encryptor = encryptor
	}

	// sharer
	if s.sharer == nil {
		sharer := new(accessSharers.Dracoon)
		sharer.SetClient(s.client)
		sharer.SetEncryptor(s.encryptor)
		if password, mErr := services.DracoonEncryptionPassword(services.MessageDomain(message)); mErr != nil {
			return errorsStorage.NewError(mErr.Error(), errorsStorage.Settings, errorsStorage.InstanceSettingsError)
		} else {
			sharer.SetPassword(password)
		}
		if iErr := sharer.Init(); iErr != nil {
			return errorsStorage.NewError(iErr.Error(), errorsStorage.Settings, errorsStorage.GeneralError)
		}
		s.sharer = sharer
	}

	// other
	s.SetExtraInfo(message.CustomerId())

	return nil
}

func (s *Dracoon) SetClient(client *servicesRestApi.Dracoon) {
	s.client = client
}

func (s *Dracoon) SetExtraInfo(info string) {
	s.extraInfo = info
}

// @Description need to be done before every uploading
func (s *Dracoon) Reset() {
	s.AbstractStorage.Reset()
	s.client.Reset()

	s.uploadChannel.id = ""
	s.uploadChannel.token = ""
}

// @Description authentication on the storage side
func (s *Dracoon) Auth() errorsStorage.ErrorContract {
	err := s.AbstractStorage.Auth()
	if err != nil {
		return err
	}

	if _, cErr := s.client.TokenValue(); cErr != nil {
		return errorsStorage.NewError(cErr.Error(), errorsStorage.Auth, errorsStorage.AuthError)
	}
	return nil
}

func (s *Dracoon) DestinationFolders(message modelsDataSqs.FileTransfererMessage) (folders []interface{}, err errorsStorage.ErrorContract) {
	var targets []*modelsDataStorageDracoon.Target
	var mErr error

	if targets, mErr = message.DracoonTargets(); mErr == nil {
		var parentRoomId int

		if parentRoomId, mErr = services.DracoonParentRoomId(services.MessageDomain(message)); mErr == nil {
			folders = make([]interface{}, len(targets))

			if len(targets) > 0 {
				for i, target := range targets {
					var room *modelsRestDracoon.NodeResponse
					if room, mErr = s.client.FindOrCreateRoom(parentRoomId, target.Path); mErr == nil {
						target.RoomId = room.Id
						folders[i] = target
					} else {
						return
					}
				}
			} else {
				folders[0] = &modelsDataStorageDracoon.Target{
					RoomId: parentRoomId,
				}
			}
		}
		if mErr != nil {
			err = errorsStorage.NewError(mErr.Error(), errorsStorage.Destination, errorsStorage.ApiError)
		}
	} else {
		err = errorsStorage.NewError(mErr.Error(), errorsStorage.Message, errorsStorage.FolderError)
	}
	return
}

func (s *Dracoon) SetDestinationFolder(folder interface{}) errorsStorage.ErrorContract {
	if err := s.AbstractStorage.SetDestinationFolder(folder); err != nil {
		return err
	}

	if target, ok := folder.(*modelsDataStorageDracoon.Target); ok {
		s.target = *target
	} else {
		return errorsStorage.NewError(fmt.Sprintf("Wrong Dracoon target: %+v", target), errorsStorage.Message, errorsStorage.FolderError)
	}
	return nil
}

// @Description: check if there's free room, if quota isn't used up etc.
// (if we're gonna upload to the storage)
func (s *Dracoon) CheckDestinationFolder() errorsStorage.ErrorContract {
	if err := s.AbstractStorage.CheckDestinationFolder(); err != nil {
		return err
	}

	var makeError = func(info *modelsRestDracoon.NodeResponse, msg string) errorsStorage.ErrorContract {
		return errorsStorage.NewError(fmt.Sprintf("Dracoon destination %s: %s", info.Name, msg), errorsStorage.Destination, errorsStorage.FolderError)
	}

	if info, cErr := s.client.Node(s.target.RoomId); cErr != nil {
		return errorsStorage.NewError(cErr.Error(), errorsStorage.Destination, errorsStorage.ApiError)
	} else {
		if !info.IsEncrypted {
			return makeError(info, "room is not encrypted but should be")
		}
		if info.Type != "room" {
			return makeError(info, "is "+info.Type+", should be a room")
		}
		if !info.Permissions.Manage && !info.Permissions.Create {
			return makeError(info, "is not accessible for our account")
		}
		if info.Quota > 0 && info.Size >= info.Quota {
			return makeError(info, "room is full")
		}
	}
	return nil
}

func (s *Dracoon) fileName() string {
	maxExtraInfoLen := 100 - len(s.meeting.User.Email) - len(s.fileInfo.Extension)

	extraInfoInFileName := system.RemoveSpecialChars(s.extraInfo, `/\<>:"|?*`)
	extraInfoInFileName = strings.Replace(extraInfoInFileName, " ", "+", -1)
	extraInfoInFileName = system.TruncateByMaxLen(extraInfoInFileName, maxExtraInfoLen)

	return fmt.Sprintf("%s--%s--%s.%s",
		s.meeting.User.Email,
		extraInfoInFileName,
		s.meeting.DateInTz().Format("2006-01-02_15-04_MST"),
		s.fileInfo.Extension)
}

func (s *Dracoon) fileNotes() string {
	notes := fmt.Sprintf("Meeting ID: %d. \n", s.meeting.Id) +
		fmt.Sprintf("Type: %s. \n", s.meeting.Type.NameAndDuration()) +
		fmt.Sprintf("Date: %s. \n", s.meeting.DateInTz().Format("Mon, 2 Jan 2006 15:04:05 MST")) +
		fmt.Sprintf("User: %s. \n", s.meeting.User.NameAndEmail()) +
		fmt.Sprintf("Guest: %s. \n", s.meeting.Guest.NameAndEmail())

	notes += fmt.Sprintf("Extra info: %s", s.extraInfo)
	return system.TruncateByMaxLen(notes, 255)
}

// @Description do some actions before starting transferring (e.g., open upload channel)
func (s *Dracoon) Start() errorsStorage.ErrorContract {
	if err := s.AbstractStorage.Start(); err != nil {
		return err
	}
	if channel, cErr := s.client.CreateUploadChannel(s.target, s.fileName(), s.fileInfo.Size, s.fileNotes()); cErr != nil {
		return errorsStorage.NewError(cErr.Error(), errorsStorage.Destination, errorsStorage.UploadError)
	} else {
		s.uploadChannel.id = channel.Id
		s.uploadChannel.token = channel.Token
	}
	return nil
}

// @Description upload data to the storage (probably chunk by chunk)
func (s *Dracoon) Upload(data []byte, downloadState modelsDataStorage.State) errorsStorage.ErrorContract {
	err := s.AbstractStorage.Upload(data, downloadState)
	if err != nil {
		return err
	}

	if s.encryptor != nil {
		s.logStatus(s.encryptor.Status())

		if encryptedData, cErr := s.encryptor.Encrypt(data); cErr != nil {
			return errorsStorage.NewError(cErr.Error(), errorsStorage.Destination, errorsStorage.EncryptionError)
		} else {
			if cErr := s.client.Upload(s.uploadChannel.token, encryptedData, downloadState.Range.ContentRangeHeader()); cErr != nil {
				return errorsStorage.NewError(cErr.Error(), errorsStorage.Destination, errorsStorage.UploadError)
			}
		}
	} else {
		if cErr := s.client.Upload(s.uploadChannel.token, data, downloadState.Range.ContentRangeHeader()); cErr != nil {
			return errorsStorage.NewError(cErr.Error(), errorsStorage.Destination, errorsStorage.UploadError)
		}
	}
	return nil
}

// @Description do some actions after completing uploading
// (e.g., close uploading channel, sharing access for Dracoon users, etc)
func (s *Dracoon) Close() errorsStorage.ErrorContract {
	var cErr error
	var fileKey modelsRestDracoon.FileKey
	var userKeys modelsRestDracoon.UserKeyPairResponse

	// encrypt file key
	if s.encryptor != nil {
		if cErr = s.encryptor.EncryptFinal(); cErr != nil {
			return errorsStorage.NewError(cErr.Error(), errorsStorage.Destination, errorsStorage.EncryptionError)
		}

		encryptionInfo := s.encryptor.EncryptionInfo()
		if encryptionInfo.ValidAesGcmInfo() {
			fileKey = modelsRestDracoon.FileKeyFromEncryptionInfo(encryptionInfo)

			// get user keys
			if userKeys, cErr = s.client.UserKeyPair(); cErr != nil {
				return errorsStorage.NewError(cErr.Error(), errorsStorage.Destination, errorsStorage.ApiError)
			}

			// encrypt file key with user public key
			var encryptedKey []byte
			if encryptedKey, cErr = s.encryptor.EncryptWithKey([]byte(fileKey.Key), []byte(userKeys.PublicKeyContainer.PublicKey)); cErr != nil {
				return errorsStorage.NewError(cErr.Error(), errorsStorage.Destination, errorsStorage.EncryptionError)
			}
			fileKey.Key = string(encryptedKey)

			// base 64 encoding
			fileKey = fileKey.BinaryToString()
		}
	}

	// close upload channel
	var file modelsRestDracoon.NodeResponse
	if file, cErr = s.client.FinishUploading(s.uploadChannel.token, s.fileName(), fileKey); cErr != nil {
		return errorsStorage.NewError(cErr.Error(), errorsStorage.Destination, errorsStorage.UploadError)
	}

	// share access
	if s.sharer != nil {
		s.logStatus(s.sharer.Status())

		s.sharer.SetKeyPair(modelsDataAccessSharer.KeyPair{
			PrivateKeyContainer: userKeys.PrivateKeyContainer,
			PublicKeyContainer:  userKeys.PublicKeyContainer,
		})

		options := modelsDataAccessSharer.DracoonOptions{
			FileId: file.Id,
		}
		if cErr := s.sharer.Share(options); cErr != nil {
			return errorsStorage.NewError(cErr.Error(), errorsStorage.Destination, errorsStorage.AccessSharingError)
		}
	}

	return nil
}

// @Description do some actions after failed downloading/uploading
func (s *Dracoon) Revert() {
	s.AbstractStorage.Revert()
	_ = s.client.CancelUploading(s.uploadChannel.token) // try and don't care of error
}

// @Description check if service works fine
func (s *Dracoon) HealthCheck() errorsStorage.ErrorContract {
	if err := s.AbstractStorage.HealthCheck(); err != nil {
		return err
	}

	// test connection to the server
	if cErr := s.client.Ping(); cErr != nil {
		return errorsStorage.NewError(cErr.Error(), errorsStorage.Destination, errorsStorage.ApiError)
	}
	return nil
}
