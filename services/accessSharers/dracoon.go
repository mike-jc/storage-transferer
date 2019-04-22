package accessSharers

import (
	"encoding/base64"
	"errors"
	"fmt"
	"math/rand"
	"service-recordingStorage/errors/rest"
	"service-recordingStorage/models/data/accessSharer"
	"service-recordingStorage/models/data/rest/dracoon"
	"service-recordingStorage/services/encryptors"
	"service-recordingStorage/services/rest/api"
	"time"
)

const SharingMaxRetries = 3
const SharingRetryDelay = 10 * time.Second

type Dracoon struct {
	client    *servicesRestApi.Dracoon
	encryptor *encryptors.AbstractEncryptorContract

	AbstractSharer
}

func (s *Dracoon) SharerName() string {
	return "Dracoon sharer"
}

func (s *Dracoon) Init() error {
	err := s.AbstractSharer.Init()
	if err != nil {
		return err
	}

	s.self = s

	return err
}

func (s *Dracoon) SetClient(client *servicesRestApi.Dracoon) {
	s.client = client
}

func (s *Dracoon) SetEncryptor(encryptor encryptors.AbstractEncryptorContract) {
	s.encryptor = &encryptor
}
func (s *Dracoon) Share(options interface{}) error {
	if err := s.AbstractSharer.Share(options); err != nil {
		return err
	}

	if s.client == nil {
		return errors.New("API client is not set for " + s.SharerName())
	}
	if s.encryptor == nil {
		return errors.New("Encryptor is not set for " + s.SharerName())
	}

	if len(s.keyPair.PrivateKeyContainer.PrivateKey) == 0 {
		return errors.New("Can not share access: empty private key for the current Dracoon account")
	}
	if len(s.password) == 0 {
		return errors.New("Can not share access: empty password for the current Dracoon account")
	}

	// parse options
	dOptions, ok := options.(modelsDataAccessSharer.DracoonOptions)
	if !ok {
		return errors.New(fmt.Sprintf("Wrong options for Dracoon sharer: %+v", options))
	}

	retries := 0
	for {
		// get missed keys
		if missingKeys, err := s.client.MissingFileKeys(dOptions.RoomId, dOptions.FileId); err != nil {
			return err
		} else {
			newFileKeys := modelsRestDracoon.UploadFileKeysRequest{
				Items: make([]*modelsRestDracoon.UploadFileKeysItem, 0),
			}

			// re-encrypt keys with users' public keys
			for userId, userFiles := range missingKeys.Items {
				if len(userFiles) == 0 {
					continue
				}

				publicKeyContainer, ok := missingKeys.Users[userId]
				if !ok {
					continue
				}

				for fileId, _ := range userFiles {
					fileKeyContainer, ok := missingKeys.Files[fileId]
					if !ok {
						continue
					}

					if decryptedFileKey, err := (*s.encryptor).DecryptWithKey([]byte(fileKeyContainer.Key), []byte(s.keyPair.PrivateKeyContainer.PrivateKey), []byte(s.password)); err != nil {
						return errors.New(fmt.Sprintf("Can not decrypt file key for file %d: %s", fileId, err.Error()))
					} else {
						if reEncryptedFileKey, err := (*s.encryptor).EncryptWithKey(decryptedFileKey, []byte(publicKeyContainer.PublicKey)); err != nil {
							return errors.New(fmt.Sprintf("Can not re-encrypt file key for user %d and file %d: %s", userId, fileId, err.Error()))
						} else {
							newFileKeys.Items = append(newFileKeys.Items, &modelsRestDracoon.UploadFileKeysItem{
								FileId: fileId,
								UserId: userId,
								FileKey: modelsRestDracoon.FileKey{
									Key:     base64.StdEncoding.EncodeToString(reEncryptedFileKey),
									IV:      fileKeyContainer.IV,
									Version: fileKeyContainer.Version,
									Tag:     fileKeyContainer.Tag,
								},
							})
						}
					}
				}
			}

			// share keys with other users who has access to this room and file
			if len(newFileKeys.Items) > 0 {
				if err := s.client.UploadFileKeys(newFileKeys); err != nil {
					// try to re-fetch missing keys if there's a conflict
					if dErr, ok := err.(errorsRest.DracoonErrorContract); ok && dErr.Response().Code == 409 {
						if retries < SharingMaxRetries {
							retries++
							time.Sleep(SharingRetryDelay + time.Duration(rand.Intn(1000))*time.Millisecond)
							continue
						} else {
							return errorsRest.NewDracoonError("Max retries number exceeded: "+dErr.Error(), dErr.Response())
						}
					}
					return err
				}
			}
			return nil
		}
	}
}
