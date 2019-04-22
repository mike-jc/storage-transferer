package modelsDataSqs

import (
	"errors"
	"fmt"
	"service-recordingStorage/models/data/instance"
	"service-recordingStorage/models/data/storage/dracoon"
)

type FileTransfererMessage struct {
	Storage  Storage                    `json:"storage"`
	Meeting  modelsDataInstance.Meeting `json:"meeting"`
	Instance Instance                   `json:"instance"`
	Extra    interface{}                `json:"extra"`
}

type Storage struct {
	Type  string      `json:"type"`
	Extra interface{} `json:"extra"`
}

type Instance struct {
	Alias  string `json:"alias"`
	Domain string `json:"domain"`
}

func (m *FileTransfererMessage) StorageExtra() (storageExtra map[string]interface{}, err error) {
	var ok bool
	if storageExtra, ok = m.Storage.Extra.(map[string]interface{}); !ok {
		err = errors.New(fmt.Sprintf("FileTransfererMessage does not have storage:extra field: %+v", *m))
		return
	}
	return
}

func (m *FileTransfererMessage) S3Extra() (s3Extra map[string]interface{}, err error) {
	var storageExtra map[string]interface{}

	if storageExtra, err = m.StorageExtra(); err != nil {
		return
	} else {
		var ok bool
		if s3Extra, ok = storageExtra["aws.s3"].(map[string]interface{}); !ok {
			err = errors.New(fmt.Sprintf("FileTransfererMessage does not have storage:extra:aws.s3 field: %+v", *m))
			return
		}
	}
	return
}

func (m *FileTransfererMessage) S3BucketName() (bucketName string, err error) {
	var s3Extra map[string]interface{}

	if s3Extra, err = m.S3Extra(); err != nil {
		return
	} else {
		if bucketName = s3Extra["bucket"].(string); bucketName == "" {
			err = errors.New(fmt.Sprintf("Wrong bucket name is empty in FileTransfererMessage in storage:extra:aws.s3 field: %+v", *m))
			return
		}
	}
	return
}

func (m *FileTransfererMessage) S3BucketKeys() (keys []string, err error) {
	var s3Extra map[string]interface{}

	if s3Extra, err = m.S3Extra(); err != nil {
		return
	} else {
		// keys should be array
		if values, ok := s3Extra["keys"].([]interface{}); !ok {
			err = errors.New(fmt.Sprintf("Wrong bucket key list in FileTransfererMessage in storage:extra:aws.s3 field: %+v", *m))
			return
		} else {
			keys = make([]string, 0)
			for _, subValue := range values {
				// each key may be array of keys...
				if subList, ok := subValue.([]interface{}); ok {
					for _, v := range subList {
						if str := v.(string); len(str) > 0 {
							keys = append(keys, str)
						}
					}
					// ...or key itself
				} else if str := subValue.(string); len(str) > 0 {
					keys = append(keys, str)
				}
			}
		}
		if len(keys) == 0 {
			err = errors.New(fmt.Sprintf("Empty bucket key list in FileTransfererMessage in storage:extra:aws.s3 field: %+v", *m))
			return
		}
	}
	return
}

func (m *FileTransfererMessage) DracoonExtra() (dracoonExtra map[string]interface{}, err error) {
	var storageExtra map[string]interface{}

	if storageExtra, err = m.StorageExtra(); err != nil {
		return
	} else {
		var ok bool
		if dracoonExtra, ok = storageExtra["dracoon"].(map[string]interface{}); !ok {
			err = errors.New(fmt.Sprintf("FileTransfererMessage does not have storage:extra:dracoon field: %+v", *m))
			return
		}
	}
	return
}

func (m *FileTransfererMessage) DracoonParentRoomId() (id int, err error) {
	var dracoonExtra map[string]interface{}

	if dracoonExtra, err = m.DracoonExtra(); err != nil {
		return
	} else {
		if value, ok := dracoonExtra["parentRoomId"].(float64); !ok {
			err = errors.New(fmt.Sprintf("Wrong parentRoomId in FileTransfererMessage in storage:extra:dracoon field: %+v", *m))
			return
		} else {
			id = int(value)
		}
	}
	return
}

func (m *FileTransfererMessage) DracoonTargets() (targets []*modelsDataStorageDracoon.Target, err error) {
	var dracoonExtra map[string]interface{}

	if dracoonExtra, err = m.DracoonExtra(); err != nil {
		return
	} else {
		if values, ok := dracoonExtra["targets"].([]interface{}); !ok {
			err = errors.New(fmt.Sprintf("Wrong targets list in FileTransfererMessage in storage:extra:dracoon field: %+v", *m))
			return
		} else {
			targets = make([]*modelsDataStorageDracoon.Target, 0)

			for _, subValue := range values {
				if v, ok := subValue.(map[string]interface{}); !ok {
					err = errors.New(fmt.Sprintf("Wrong targets list in FileTransfererMessage in storage:extra:dracoon field: %+v", *m))
					return
				} else {
					if path := v["path"].(string); len(path) > 0 {
						targets = append(targets, &modelsDataStorageDracoon.Target{
							Path:       path,
							Expiration: v["expiration"].(string),
						})
					}
				}
			}
		}
	}
	return
}

// customer ID is optional (for some instances only)
func (m *FileTransfererMessage) CustomerId() string {
	if extra, ok := m.Extra.(map[string]interface{}); ok {
		if customerId, ok := extra["customerId"]; ok {
			return customerId.(string)
		}
	}
	return ""
}
